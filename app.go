package main

import (
	"context"
	"strings"
)

// App 是暴露给前端的 Wails 绑定结构体。
// 前端通过 window.go.main.App.<Method> 调用这些方法。
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

// startup 在应用启动时由 Wails 调用，保存上下文。
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// DetectAll 返回所有 IDE（内置 + 自定义）及各自检测到的 GitLens。
func (a *App) DetectAll() detectResult {
	return detectAll()
}

// AddCustomDir 添加并保存自定义目录，返回刷新后的检测结果。
func (a *App) AddCustomDir(dir string) (detectResult, error) {
	if err := addCustomDir(dir); err != nil {
		return detectResult{}, err
	}
	return detectAll(), nil
}

// RemoveCustomDir 删除自定义目录，返回刷新后的检测结果。
func (a *App) RemoveCustomDir(dir string) (detectResult, error) {
	if err := removeCustomDir(dir); err != nil {
		return detectResult{}, err
	}
	return detectAll(), nil
}

// Activate 批量激活，返回每个目录的操作结果。
func (a *App) Activate(dirs []string) []actionResult {
	return batchAction(dirs, activateExtensionDir, "激活成功")
}

// Restore 批量恢复，返回每个目录的操作结果。
func (a *App) Restore(dirs []string) []actionResult {
	return batchAction(dirs, restoreExtensionDir, "恢复成功")
}

// actionResult 单个扩展目录的操作结果。
type actionResult struct {
	DirPath string `json:"dirPath"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// batchAction 对 dirs 逐个执行 fn，汇总每个目录的成功/失败信息。
func batchAction(dirs []string, fn func(string) error, okMsg string) []actionResult {
	results := make([]actionResult, 0, len(dirs))
	for _, d := range dirs {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}
		if err := fn(d); err != nil {
			results = append(results, actionResult{DirPath: d, OK: false, Message: err.Error()})
		} else {
			results = append(results, actionResult{DirPath: d, OK: true, Message: okMsg})
		}
	}
	return results
}
