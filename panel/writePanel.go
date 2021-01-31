package panel

import (
	"fmt"

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

var img *canvas.Image = canvas.NewImageFromResource(theme.FyneLogo())

//GetWritePanel 获取写入面板
func GetWritePanel(topWindow fyne.Window, globalData *Data) *fyne.Container {
	labelIn := widget.NewLabel("")
	labelOut := widget.NewLabel("")

	box5 := fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(500, 300)),
		img,
	)
	box1 := container.NewHBox(
		widget.NewButton("intput File (.jpg or .png)", func() {
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil && reader == nil {
					return
				}
				if err != nil {
					dialog.ShowError(err, topWindow)
					return
				}
				img, path, err := readImg(reader, globalData)
				labelIn.SetText(path)
				globalData.imgData = img.Resource.Content()
				box5.Objects = []fyne.CanvasObject{img}
				box5.Refresh()

			}, topWindow)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", "jpeg"}))
			fd.Show()
		}),
		widget.NewSeparator(),
		labelIn,
	)

	//----------------
	box2 := container.NewHBox(
		widget.NewButton("output Directory", func() {
			dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, topWindow)
					return
				}
				if list == nil {
					return
				}
				out := fmt.Sprintf("Folder is:%s", list.String())
				labelOut.SetText(out)
				globalData.outputPath = list.String()
			}, topWindow)
		}),
		widget.NewSeparator(),
		labelOut,
	)

	//----------------
	entry := widget.NewEntry()
	entry.SetPlaceHolder("copyrights info")
	entry.OnChanged = func(s string) {
		globalData.watermark = s
	}
	box3 := container.NewVBox(
		entry,
	)
	//----------------
	box4 := container.NewVBox(
		widget.NewButton("start", func() {
			prog := dialog.NewProgress("processing", "please wait...", topWindow)
			prog.Show()
			fmt.Printf("input: %s\n", globalData.inputName)
			fmt.Printf("input: %s\n", globalData.inputExt)
			fmt.Printf("output: %s\n", globalData.outputPath)
			fmt.Printf("watermark: %s\n", globalData.watermark)
			fmt.Printf("img_data: %d\n", len(globalData.imgData))
			ch := make(chan float64)
			go func(ch chan float64) {
				for num := range ch {
					prog.SetValue(num)
				}
				prog.SetValue(1)
				prog.Hide()
			}(ch)
			data, err := shadow.DealImg(globalData.imgData, globalData.watermark, globalData.inputExt, ch)
			if err != nil {
				fmt.Printf("%s", err.Error())
				dialog.ShowConfirm("ops!", "there is something wrong happend!:"+err.Error(), func(b bool) {
					prog.Hide()
				}, topWindow)
				return
			}
			err = writeImg(data, globalData.outputPath, topWindow, globalData, ch)
			if err != nil {
				fmt.Printf("%s", err.Error())
				dialog.ShowConfirm("ops!", "there is something wrong happend!:"+err.Error(), func(b bool) {
					prog.Hide()
				}, topWindow)
				return
			}

			prog.Show()
		}),
	)
	//----------------

	grid := container.NewVBox(
		widget.NewLabelWithStyle("write copyrights info", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewSeparator(),
		box1,
		box2,
		box3,
		box4,
		box5)

	return grid
}
