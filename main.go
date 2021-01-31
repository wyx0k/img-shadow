package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"wyx0k.com/img-shadow/panel"
)

var topWindow fyne.Window

var globalData panel.Data = panel.Data{}

func main() {
	application := app.NewWithID("img-shadow")
	theme := &imgShadowTheme{}
	application.Settings().SetTheme(theme)
	mainWindow := application.NewWindow("img-shadow")
	topWindow = mainWindow
	header := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabelWithStyle("wyx0k for ellie", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewSeparator(),
	)
	// 右侧内容
	content := container.NewMax()
	// 左侧菜单
	menu := widget.NewHBox(widget.NewVBox(
		widget.NewButton("write", func() {
			write := fyne.CanvasObject(panel.GetWritePanel(topWindow, &globalData))
			onSelected(content, &write)
		}),
		widget.NewButton("read", func() {
			read := fyne.CanvasObject(panel.GetReadPanel(topWindow, &globalData))
			onSelected(content, &read)
		}),
	), widget.NewSeparator())
	// 整体结构
	border := container.NewBorder(header, nil, menu, nil, content)

	topWindow.SetContent(border)
	topWindow.SetFixedSize(true)
	topWindow.Resize(fyne.NewSize(800, 500))
	topWindow.ShowAndRun()
}

func onSelected(content *fyne.Container, item *fyne.CanvasObject) {
	content.Objects = []fyne.CanvasObject{*item}
	content.Refresh()
}
