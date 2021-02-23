package panel

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"wyx0k.com/img-shadow/shadow"
)

var imgRead *canvas.Image = canvas.NewImageFromResource(theme.FyneLogo())

//GetReadPanel 获取读取面板
func GetReadPanel(topWindow fyne.Window, globalData *Data) *fyne.Container {
	labelIn := widget.NewLabel("")
	box2 := fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(500, 300)),
		imgRead,
	)
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
				imgRead, path, err := readImg(reader, globalData)
				data, err := shadow.ReadImg(imgRead.Resource.Content(), globalData.inputExt)
				res := fyne.NewStaticResource("readed", data)

				box2.Objects = []fyne.CanvasObject{canvas.NewImageFromResource(res)}
				labelIn.SetText(path)
				box2.Refresh()
			}, topWindow)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".txt"}))
			fd.Show()
		}),
		widget.NewSeparator(),
		labelIn,
	)

	box := container.NewVBox(
		widget.NewLabelWithStyle("read copyrights info", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewSeparator(),
		box1,
		box2,
	)
	return box
}
