package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "GitLens Activator",
		Width:  950,
		Height: 768,
		// 工具栏单行需要 920px（内容 876 + 左右 22px 边距）。
		// 950 为实测可行的最小宽度，再窄工具栏就会换行。
		MinWidth: 950,
		// 固定部分（标题栏 71 + 工具栏 77 + 状态行 41 + 底栏 37）共 226px，
		// 卡片单行 108px：500 可完整显示 2 行并露出第 3 行，也容得下弹窗。
		MinHeight: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
