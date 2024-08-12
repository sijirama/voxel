package main

import (
	"context"
	"embed"
	"voxel/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	util := utils.NewUtils()

	err := wails.Run(&options.App{
		Title:  "voxel",
		Width:  600,
		Height: 708,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
			util,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
