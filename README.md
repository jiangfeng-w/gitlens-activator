# GitLens Activator

为 GitLens Pro 提供一键激活 / 恢复的桌面 GUI 工具，逻辑复刻自 gitlens-pro（原项目仅作参考，未改动）。

## 特性

- 桌面端应用，双击即开原生窗口，不再启动本地服务或依赖浏览器
- 以列表形式展示所有 IDE 的安装状态与检测到的 GitLens 版本
- 支持勾选多个 IDE 批量激活 / 批量恢复，也支持单版本操作
- 支持自定义扩展目录，添加后持久保存（保存在系统用户配置目录）
- 支持 GitLens 15.x ~ 19.x：
  - 15.x：订阅枚举替换
  - 16.x ~ 17.x：license 对象注入（与原项目 16.x 逻辑一致）
  - 18.x ~ 19.x：license 对象注入 + **未登录激活**（强制初始订阅为 Pro，无需登录即可解锁）
  - 17.x：额外解锁提交图（graph.js）
- 激活前自动备份原文件为 `.backup`，恢复时一键还原
- 已激活的版本重复激活时幂等跳过，不会重复注入

## 内置支持的 IDE

| IDE | 扩展目录 |
|-----|----------|
| VS Code | `~/.vscode/extensions` |
| VS Code Insiders | `~/.vscode-insiders/extensions` |
| Cursor | `~/.cursor/extensions` |
| Windsurf | `~/.windsurf/extensions` |
| Kiro | `~/.kiro/extensions` |
| Antigravity | `~/.antigravity/extensions` |
| Trae | `~/.trae/extensions` |
| Trae CN | `~/.trae-cn/extensions` |
| CodeBuddy | `~/.codebuddy/extensions` |
| CodeBuddy CN | `~/.codebuddycn/extensions` |
| Qoder | `~/.qoder/extensions` |
| Qoder CN | `~/.qoder-cn/extensions` |
| CatPaw（美团） | `~/.catpaw/extensions` |

> Qoder / CatPaw 的扩展目录为按 VS Code 系惯例推断，若与实际不符，请用「添加目录」手动添加真实路径。

## 开发与构建

前置要求：

- Go 1.22+
- Wails CLI：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- 各平台系统 WebView 组件：
  - Windows：WebView2 Runtime（Win10/11 自带）
  - macOS：Xcode Command Line Tools（WKWebView）
  - Linux：`libwebkit2gtk-4.0-dev`、`libgtk-3-dev`

```bash
go test ./...    # 测试
wails dev        # 开发模式（热重载）
wails build      # 生产构建，输出到 build/bin/
```

## 使用

双击构建产物启动应用窗口：

- **扫描**：重新检测所有 IDE 目录
- **激活选中 / 恢复选中**：勾选左侧复选框后，批量处理该 IDE 下的所有 GitLens 版本
- **每行按钮**：对单个 GitLens 版本执行激活 / 恢复
- **自定义目录**：输入任意扩展目录路径后点「添加目录」，保存到系统用户配置目录；点行内「删除」可移除

自定义目录配置文件位置：

- Windows：`%AppData%\gitlens-activator\custom_dirs.json`
- macOS：`~/Library/Application Support/gitlens-activator/custom_dirs.json`
- Linux：`~/.config/gitlens-activator/custom_dirs.json`

## 说明

- 激活前会自动备份原文件为 `.backup`（已备份则跳过，重复激活幂等返回成功）
- 恢复时会将 `.backup` 还原为原文件；若没有备份（例如被其它工具修改过），请卸载重装扩展还原
- 支持 `eamodio.gitlens-x.y.z` 和 `eamodio.gitlens-x.y.z-universal` 两种目录格式
- 18.x 起 GitLens 移除了 17.x 的提交图限制逻辑，因此 18.x+ 无需额外 patch graph.js
- 18.x+ 的未登录激活补丁会把 GitLens 的本地订阅状态也存为 Pro（在 IDE 的 globalStorage 中）；仅还原 `.backup` 文件不会清除该状态，如需完全还原请在 GitLens 中退出登录或清除 IDE 的 GitLens 存储数据
