package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type editorCandidate struct {
	Key           string             `json:"key"`
	Name          string             `json:"name"`
	ExtensionsDir string             `json:"extensionsDir"`
	Installed     bool               `json:"installed"`
	Custom        bool               `json:"custom"`
	Extensions    []gitlensExtension `json:"extensions,omitempty"`
}

type gitlensExtension struct {
	DirName   string `json:"dirName"`
	DirPath   string `json:"dirPath"`
	Version   string `json:"version"`
	Universal bool   `json:"universal"`
	HasBackup bool   `json:"hasBackup"`
	Activated bool   `json:"activated"`
}

type detectResult struct {
	Presets []editorCandidate `json:"presets"`
	Customs []editorCandidate `json:"customs"`
}

var (
	patternVSCode = `^eamodio\.gitlens-(\d+\.\d+\.\d+)$`
	patternOther  = `^eamodio\.gitlens-(\d+\.\d+\.\d+)-universal$`
	// 16.x~19.x 注入点：let r,s,n,o={id:e.user.id,name:e.user.name
	// 支持多变量声明，与原项目 processVersion16File 保持一致
	licenseInjectionPattern = regexp.MustCompile(`let ([a-zA-Z,]+)=\{id:e\.user\.id,name:e\.user\.name`)
	// 17.x graph.js 提交图解锁注入点（18.x 起已不存在）
	graphInjectionPattern = regexp.MustCompile(`(\(|=)(this\.(graphS|s)tate\.allowed)`)
	// 16.x+ 注入的 mock 对象，覆盖函数参数 e（与原项目 insertCode 一致）
	insertCode = `e={user:{id:"88888888-8888-8888-8888-888888888888",name:"Neo",email:"x@x.com",status:"activated",createdDate:"2000-01-01T00:00:00.000Z"},licenses:{paidLicenses:{},effectiveLicenses:{"gitlens-pro":{organizationId:"Linux",latestStatus:"active",latestStartDate:"2024-01-01",latestEndDate:"2999-01-01",reactivationCount:99,nextOptInDate:"2999-01-01"}}},nextOptInDate:"2999-01-01"};`
	// mock 用户 ID，用于检测文件是否已被激活
	mockUserID = "88888888-8888-8888-8888-888888888888"
	// 18.x+ 未登录激活注入点：changeSubscription 的默认订阅表达式
	// 形如 e??={plan:{actual:(0,ty.le)("community",!1,0,void 0),effective:(0,ty.le)("community",!1,0,void 0)},account:void 0,state:iJ.z.Community}
	// 16.x/17.x 没有这段（走登录后 i9 校验路径），匹配不到时静默跳过。
	noLoginPatchPattern = regexp.MustCompile(`e\?\?=\{plan:\{actual:\(0,([A-Za-z_$][A-Za-z0-9_$]*)\.le\)\("community",!1,0,void 0\),effective:\(0,([A-Za-z_$][A-Za-z0-9_$]*)\.le\)\("community",!1,0,void 0\)\},account:void 0,state:([A-Za-z_$][A-Za-z0-9_$]*)\.z\.Community\}`)
	// license 注入（补丁 1）的检测标记：该字符串只出现在 insertCode 里
	licenseInjectionMarker = `licenses:{paidLicenses:{},effectiveLicenses:{"gitlens-pro"`
	// 未登录激活（补丁 2）的检测标记：changeSubscription 被改为无条件 e= 覆盖
	noLoginPatchMarker = `changeSubscription(e,t,i){e={plan:{actual:(0,`
)

// presetEditorCandidates 返回内置 IDE 的扩展目录。
// 全部为 VS Code 系 IDE，扩展统一存放在 ~/.<ide>/extensions。
// Qoder / CatPaw 目录为社区惯例推断，若与实际不符可用自定义目录添加。
func presetEditorCandidates() []editorCandidate {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	presets := []struct{ key, name, dir string }{
		{"vscode", "VS Code", filepath.Join(home, ".vscode", "extensions")},
		{"vscode-insiders", "VS Code Insiders", filepath.Join(home, ".vscode-insiders", "extensions")},
		{"cursor", "Cursor", filepath.Join(home, ".cursor", "extensions")},
		{"windsurf", "Windsurf", filepath.Join(home, ".windsurf", "extensions")},
		{"kiro", "Kiro", filepath.Join(home, ".kiro", "extensions")},
		{"antigravity", "Antigravity", filepath.Join(home, ".antigravity", "extensions")},
		{"trae", "Trae", filepath.Join(home, ".trae", "extensions")},
		{"traecn", "Trae CN", filepath.Join(home, ".trae-cn", "extensions")},
		{"codebuddy", "CodeBuddy", filepath.Join(home, ".codebuddy", "extensions")},
		{"codebuddycn", "CodeBuddy CN", filepath.Join(home, ".codebuddycn", "extensions")},
		{"qoder", "Qoder", filepath.Join(home, ".qoder", "extensions")},
		{"qodercn", "Qoder CN", filepath.Join(home, ".qoder-cn", "extensions")},
		{"catpaw", "CatPaw", filepath.Join(home, ".catpawai", "extensions")},
	}
	candidates := make([]editorCandidate, 0, len(presets))
	for _, p := range presets {
		candidates = append(candidates, editorCandidate{
			Key:           p.key,
			Name:          p.name,
			ExtensionsDir: p.dir,
		})
	}
	return candidates
}

// ---- 自定义目录持久化（config 文件放在用户配置目录） ----

func configFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "gitlens-activator", "custom_dirs.json"), nil
}

func loadCustomDirs() []string {
	path, err := configFilePath()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var dirs []string
	if err := json.Unmarshal(data, &dirs); err != nil {
		return nil
	}
	return dirs
}

func saveCustomDirs(dirs []string) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(dirs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// addCustomDir 新增自定义目录（去重后保存）
func addCustomDir(dir string) error {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return fmt.Errorf("目录不能为空")
	}
	st, err := os.Stat(dir)
	if err != nil || !st.IsDir() {
		return fmt.Errorf("目录不存在: %s", dir)
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	dirs := loadCustomDirs()
	for _, d := range dirs {
		dAbs, err := filepath.Abs(d)
		if err == nil && strings.EqualFold(dAbs, abs) {
			return fmt.Errorf("目录已存在: %s", d)
		}
	}
	dirs = append(dirs, abs)
	return saveCustomDirs(dirs)
}

// removeCustomDir 删除自定义目录
func removeCustomDir(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	dirs := loadCustomDirs()
	kept := dirs[:0]
	removed := false
	for _, d := range dirs {
		dAbs, err := filepath.Abs(d)
		if err == nil && strings.EqualFold(dAbs, abs) {
			removed = true
			continue
		}
		kept = append(kept, d)
	}
	if !removed {
		return fmt.Errorf("目录不在自定义列表中")
	}
	return saveCustomDirs(kept)
}

// ---- 检测 ----

func detectAll() detectResult {
	res := detectResult{}
	for _, c := range presetEditorCandidates() {
		res.Presets = append(res.Presets, detectOne(c))
	}
	for _, d := range loadCustomDirs() {
		res.Customs = append(res.Customs, detectOne(editorCandidate{
			Key:           "custom-" + d,
			Name:          filepath.Base(d),
			ExtensionsDir: d,
			Custom:        true,
		}))
	}
	return res
}

func detectOne(c editorCandidate) editorCandidate {
	if st, err := os.Stat(c.ExtensionsDir); err == nil && st.IsDir() {
		exts, _ := findGitLensExtensions(c.ExtensionsDir)
		c.Installed = true
		c.Extensions = exts
	} else {
		c.Installed = false
	}
	return c
}

func findGitLensExtensions(extensionsDir string) ([]gitlensExtension, error) {
	entries, err := os.ReadDir(extensionsDir)
	if err != nil {
		return nil, err
	}
	var result []gitlensExtension
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		var pattern *regexp.Regexp
		if strings.HasSuffix(name, "-universal") {
			pattern = regexp.MustCompile(patternOther)
		} else {
			pattern = regexp.MustCompile(patternVSCode)
		}
		m := pattern.FindStringSubmatch(name)
		if len(m) < 2 {
			continue
		}
		ver := m[1]
		backup := filepath.Join(extensionsDir, name, "dist", "gitlens.js.backup")
		_, err := os.Stat(backup)
		result = append(result, gitlensExtension{
			DirName:   name,
			DirPath:   filepath.Join(extensionsDir, name),
			Version:   ver,
			Universal: strings.HasSuffix(name, "-universal"),
			HasBackup: err == nil,
			Activated: isExtensionActivated(filepath.Join(extensionsDir, name), ver),
		})
	}
	return result, nil
}

func extractVersionFromDirName(name string) (string, error) {
	pattern := patternVSCode
	if strings.HasSuffix(name, "-universal") {
		pattern = patternOther
	}
	m := regexp.MustCompile(pattern).FindStringSubmatch(name)
	if len(m) < 2 {
		return "", fmt.Errorf("invalid extension dir name: %s", name)
	}
	return m[1], nil
}

// isExtensionActivated 读取 gitlens.js 检测激活补丁标记，判断该扩展是否已被本工具激活。
// 15.x 无唯一标记（枚举替换），以存在备份文件视为已激活。
// 16.x~19.x：license 注入标记（补丁 1）存在即视为已激活；
// 若该版本有未登录补丁点但补丁 2 缺失，则视为未完全激活。
func isExtensionActivated(dirPath, version string) bool {
	jsFile := filepath.Join(dirPath, "dist", "gitlens.js")
	data, err := os.ReadFile(jsFile)
	if err != nil {
		return false
	}
	content := string(data)
	if hasLicenseInjection(content) {
		if noLoginPatchPattern.MatchString(content) && !hasNoLoginPatch(content) {
			// 有 18.x+ 注入点但未登录补丁缺失：其他工具只打了 license 补丁，视为未激活
			return false
		}
		return true
	}
	// 无 license 标记：
	// 15.x 无唯一标记（枚举替换），以 mock 用户 ID 检测；
	// 只有 15.x 才用备份存在性兜底（本工具激活 15.x 必然生成 .backup，
	// 且恢复后不删除备份，无法用备份区分，故仅作 15.x 的弱兜底）。
	if strings.Contains(content, mockUserID) {
		return true
	}
	if strings.Split(version, ".")[0] == "15" {
		if _, err := os.Stat(jsFile + ".backup"); err == nil {
			return true
		}
	}
	return false
}

// ---- 备份 / 恢复 ----

func ensureBackup(jsFilePath string) error {
	backupPath := jsFilePath + ".backup"
	if _, err := os.Stat(backupPath); err == nil {
		return nil
	}
	content, err := os.ReadFile(jsFilePath)
	if err != nil {
		return err
	}
	return os.WriteFile(backupPath, content, 0644)
}

func restoreBackup(filePath string) error {
	backupPath := filePath + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file %s does not exist", backupPath)
	}
	// 用复制而非移动，恢复后保留 .backup 文件，避免多次恢复时报 "no backup files found"
	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("read backup failed: %w", err)
	}
	return os.WriteFile(filePath, content, 0644)
}

// ---- 激活 ----

// activateForVersion15 15.x：枚举替换（长名字在前，避免短前缀污染）
func activateForVersion15(content string) string {
	replacements := []struct{ old, new string }{
		{"qn.CommunityWithAccount", "qn.Enterprise"},
		{"qn.Community", "qn.Enterprise"},
		{"qn.Pro", "qn.Enterprise"},
	}
	for _, r := range replacements {
		content = strings.ReplaceAll(content, r.old, r.new)
	}
	return content
}

// hasLicenseInjection 报告是否已注入过 mock license（补丁 1）
func hasLicenseInjection(content string) bool {
	return strings.Contains(content, licenseInjectionMarker)
}

// hasNoLoginPatch 报告是否已打过未登录激活补丁（补丁 2）
func hasNoLoginPatch(content string) bool {
	return strings.Contains(content, noLoginPatchMarker)
}

// applyNoLoginPatch 18.x+ 未登录激活：把 changeSubscription 的默认订阅改成 Pro。
// 原代码：e??={plan:{actual:community,effective:community},account:void 0,state:Community}
// 改为：  e={plan:{actual:pro,effective:pro},account:mock}
// 使用 e= 无条件覆盖，因此无论登录 / 登出 / 本地存储状态如何，订阅始终为 Pro，
// 从而在未登录的情况下也能解锁 Pro 功能。
func applyNoLoginPatch(content string) (string, bool, error) {
	m := noLoginPatchPattern.FindStringSubmatch(content)
	if m == nil {
		return content, false, nil
	}
	if m[1] != m[2] {
		return content, false, fmt.Errorf("no-login patch: plan module mismatch %s vs %s", m[1], m[2])
	}
	replacement := fmt.Sprintf(
		`e={plan:{actual:(0,%s.le)("pro",!1,0,void 0),effective:(0,%s.le)("pro",!1,0,void 0)},account:{id:"88888888-8888-8888-8888-888888888888",name:"Neo",email:"x@x.com",verified:!0,createdOn:"2000-01-01T00:00:00.000Z"}}`,
		m[1], m[1])
	return strings.Replace(content, m[0], replacement, 1), true, nil
}

// activateForVersion16 16.x~19.x 及未知新版本通用激活，包含两个独立补丁：
//  1. license 注入：在 i9/iL 的账号对象声明前覆盖函数参数 e，
//     使登录后的订阅校验也拿到 pro license（与原项目 processVersion16File 一致）；
//  2. 未登录激活：强制 changeSubscription 的初始订阅为 Pro，
//     使从未登录过的 IDE 也能直接解锁 Pro 功能（18.x+ 才有该注入点）。
func activateForVersion16(content string) (string, bool, error) {
	changed := false
	patched := content

	if !hasLicenseInjection(patched) {
		if matches := licenseInjectionPattern.FindStringSubmatch(patched); len(matches) >= 2 {
			exactMatch := fmt.Sprintf("let %s={id:e.user.id,name:e.user.name", matches[1])
			patched = strings.Replace(patched, exactMatch, insertCode+exactMatch, 1)
			changed = true
		}
	}

	if !hasNoLoginPatch(patched) {
		p, ok, err := applyNoLoginPatch(patched)
		if err != nil {
			return patched, changed, err
		}
		if ok {
			patched = p
			changed = true
		}
	}

	return patched, changed, nil
}

// activateForVersion17Graph 17.x 提交图解锁：所有 this.graphState.allowed 前加感叹号（与原项目一致）
func activateForVersion17Graph(content string) (string, bool) {
	if !graphInjectionPattern.MatchString(content) {
		return content, false
	}
	return graphInjectionPattern.ReplaceAllString(content, `$1!$2`), true
}

// activateExtensionDir 激活单个 GitLens 扩展目录（dirPath 形如 .../eamodio.gitlens-17.0.1）
func activateExtensionDir(dirPath string) error {
	version, err := extractVersionFromDirName(filepath.Base(dirPath))
	if err != nil {
		return err
	}
	jsFile := filepath.Join(dirPath, "dist", "gitlens.js")
	if _, err := os.Stat(jsFile); os.IsNotExist(err) {
		return fmt.Errorf("gitlens.js not found")
	}
	content, err := os.ReadFile(jsFile)
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}
	orig := string(content)
	major := strings.Split(version, ".")[0]

	var newContent string
	switch major {
	case "15":
		newContent = activateForVersion15(orig)
		if newContent == orig {
			return fmt.Errorf("no patch applied (already activated?)")
		}
	default:
		// 16.x ~ 19.x 及未知新版本统一走 16+ 注入逻辑（license 注入 + 未登录激活）
		if hasLicenseInjection(orig) && hasNoLoginPatch(orig) {
			// 两个补丁都已存在，幂等成功
			return nil
		}
		if hasLicenseInjection(orig) && !noLoginPatchPattern.MatchString(orig) {
			// 16.x/17.x 只有 license 注入点：license 已注入即视为已激活，幂等成功
			return nil
		}
		patched, changed, err := activateForVersion16(orig)
		if err != nil {
			return err
		}
		if !changed {
			return fmt.Errorf("injection point not found in gitlens.js")
		}
		newContent = patched
	}

	if err := ensureBackup(jsFile); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}
	if err := os.WriteFile(jsFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	if major == "17" {
		// 17.x 需要额外 patch webviews/graph.js 解锁提交图
		graphFile := filepath.Join(dirPath, "dist", "webviews", "graph.js")
		if _, err := os.Stat(graphFile); err == nil {
			graphContent, err := os.ReadFile(graphFile)
			if err != nil {
				return fmt.Errorf("read graph failed: %w", err)
			}
			newGraph, ok := activateForVersion17Graph(string(graphContent))
			if !ok {
				return fmt.Errorf("graph injection point not found")
			}
			if err := ensureBackup(graphFile); err != nil {
				return fmt.Errorf("graph backup failed: %w", err)
			}
			if err := os.WriteFile(graphFile, []byte(newGraph), 0644); err != nil {
				return fmt.Errorf("write graph failed: %w", err)
			}
		}
	}
	return nil
}

// restoreExtensionDir 恢复单个 GitLens 扩展目录的备份
func restoreExtensionDir(dirPath string) error {
	var restored []string
	targets := []string{
		filepath.Join(dirPath, "dist", "gitlens.js"),
		filepath.Join(dirPath, "dist", "webviews", "graph.js"),
	}
	for _, t := range targets {
		if _, err := os.Stat(t + ".backup"); err == nil {
			if err := restoreBackup(t); err != nil {
				return err
			}
			restored = append(restored, filepath.Base(t))
		}
	}
	if len(restored) == 0 {
		return fmt.Errorf("no backup files found (if patched by another tool, reinstall the extension to restore)")
	}
	return nil
}
