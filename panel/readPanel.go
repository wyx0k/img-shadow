package panel

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
)

//GetReadPanel 获取读取面板
func GetReadPanel(topWindow fyne.Window, globalData *Data) *fyne.Container {
	box1 := container.NewHBox(
		widget.NewButton("select File (.jpg or .png)", func() {
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil && reader == nil {
					return
				}
				if err != nil {
					dialog.ShowError(err, topWindow)
					return
				}

				// fileOpened(reader)
			}, topWindow)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".txt"}))
			fd.Show()
		}),
		widget.NewSeparator(),
	)
	box2 := fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(500, 300)),
		canvas.NewRasterWithPixels(rgbGradient),
	)
	box := container.NewVBox(
		widget.NewLabelWithStyle("read copyrights info", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewSeparator(),
		box1,
		box2,
	)
	return box
}
