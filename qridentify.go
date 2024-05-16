package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/ncruces/zenity"
	"image"
	"os"
	"strings"
)

func init() {
	//设置中文字体:解决中文乱码问题
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("QR code recognition")
	w.Resize(fyne.NewSize(800, 600))

	/*显示结果的对话框*/
	resultEntry := widget.NewEntry()
	//resultEntry.SetText(result)
	imageContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(800, 400)))
	/*设置一个接收返回路径的参数*/
	filePathLabel := widget.NewLabel("")
	//imageContainer := container.NewVBox(filePathLabel)
	/*选择图片按钮*/
	picBtn := widget.NewButton("打开", func() {
		/*利用zenity库选择文件*/
		file, err := zenity.SelectFile()
		if err != nil {
			fmt.Println(err)
			return
		}
		filePathLabel.SetText(file)

		// 清空imageContainer的内容
		imageContainer.Objects = []fyne.CanvasObject{}
		//fmt.Println(file)
		/*创建一个容器，用于接收图片*/
		imageFile := canvas.NewImageFromFile(filePathLabel.Text)
		imageFile.FillMode = canvas.ImageFillOriginal
		imageFile.ScaleMode = canvas.ImageScaleSmooth
		//imageContainer = container.New(layout.NewGridWrapLayout(fyne.NewSize(800, 400)), imageFile)

		/*解码*/
		fl, err := os.Open(filePathLabel.Text)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fl.Close()

		img, _, _ := image.Decode(fl)

		bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

		/*decode image*/
		qrReader := qrcode.NewQRCodeReader()
		result, _ := qrReader.Decode(bmp, nil)
		/*如果result为空则输出识别失败*/
		if result == nil {
			resultEntry.SetText("识别失败")
			return
		}
		/*把图片添加到imageContainer容器中*/
		imageContainer.Add(imageFile)
		// 刷新窗口内容
		w.Content().Refresh()

		/*将识别到的文本写入到输入框*/
		resultEntry.SetText(result.String())

	})
	/*提示信息*/
	tipsLabel := widget.NewLabel("扫码结果：")

	//filePath := "D:\\program\\go\\goProgram\\2024\\05\\15\\qrcodeproject\\微信图片_20240515224901.jpg"

	content := container.NewVBox(picBtn, imageContainer, tipsLabel, resultEntry)

	//fmt.Println(qrmatrix.Content)

	w.SetContent(content)
	w.ShowAndRun()

}
