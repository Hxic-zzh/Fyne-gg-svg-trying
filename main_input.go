// main_input.go
package main

import (
	"2025-12-18-ggAndPng/tools"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// BuildGoogleInputPage 返回一个带谷歌风格输入框的 fyne.CanvasObject
func BuildGoogleInputPage() fyne.CanvasObject {
	// 1. 蓝色主题大号输入框（英文）
	blueInput := tools.NewMaterialEntry("400 60 22", 400, 60)
	blueInput.SetFontPath("ttf/english.ttf")
	blueInput.SetStyle(tools.MaterialEntryStyle{
		Width:           400,
		Height:          60,
		FontSize:        22,
		LabelColor:      color.RGBA{82, 100, 174, 255},
		TextColor:       color.RGBA{82, 100, 174, 255},
		BorderColor:     color.RGBA{82, 100, 174, 255},
		UnderlineColor:  color.RGBA{82, 100, 174, 255},
		UnderlineHeight: 5,
	})
	blueInput.SetCustomBackground(color.White)
	blueInput.SetCornerRadius(16)

	// 2. 红色主题中文输入框
	redInput := tools.NewMaterialEntry("中文输入", 400, 60)
	redInput.SetFontPath("ttf/chinese.ttf")
	redInput.SetStyle(tools.MaterialEntryStyle{
		Width:           400,
		Height:          60,
		FontSize:        24,
		LabelColor:      color.RGBA{244, 67, 54, 255},
		TextColor:       color.RGBA{244, 67, 54, 255},
		BorderColor:     color.RGBA{244, 67, 54, 255},
		UnderlineColor:  color.RGBA{244, 67, 54, 255},
		UnderlineHeight: 5,
	})
	redInput.SetCustomBackground(color.White)
	redInput.SetCornerRadius(16)

	// 3. 绿色主题输入框
	greenInput := tools.NewMaterialEntry("请输入密码", 400, 60)
	greenInput.SetFontPath("ttf/chinese.ttf")
	greenInput.SetStyle(tools.MaterialEntryStyle{
		Width:           400,
		Height:          60,
		FontSize:        20,
		LabelColor:      color.RGBA{76, 175, 80, 255},
		TextColor:       color.RGBA{76, 175, 80, 255},
		BorderColor:     color.RGBA{76, 175, 80, 255},
		UnderlineColor:  color.RGBA{76, 175, 80, 255},
		UnderlineHeight: 4,
	})
	greenInput.SetCustomBackground(color.White)
	greenInput.SetCornerRadius(12)

	return container.NewPadded(
		container.NewVBox(
			widget.NewLabelWithStyle("谷歌风格输入框演示", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			widget.NewLabel("蓝色主题英文输入框:"),
			container.NewCenter(container.NewPadded(blueInput.WrapWithIsolationContainer())),
			widget.NewLabel("红色主题中文输入框:"),
			container.NewCenter(container.NewPadded(redInput.WrapWithIsolationContainer())),
			widget.NewLabel("绿色主题密码输入框:"),
			container.NewCenter(container.NewPadded(greenInput.WrapWithIsolationContainer())),
		),
	)
}
