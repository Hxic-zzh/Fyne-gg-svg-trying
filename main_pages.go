// main_pages.go
package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"2025-12-18-ggAndPng/tools"
)

// BuildMainPage 构建主页面（包含原来的所有按钮内容）
func BuildMainPage() fyne.CanvasObject {
	// 创建不同颜色的按钮
	redStyle := tools.ParticleButtonStyle{
		BaseColor:     color.RGBA{R: 255, G: 100, B: 100, A: 255},
		CanvasBorder:  3,
		CanvasOffsetY: -3,
	}
	redBtn := tools.NewParticleButtonWithStyle(
		func() { println("红色按钮被点击了！") },
		"红色按钮",
		redStyle,
	)

	// 新增：无粒子特效的红色按钮
	noParticleBtn := tools.NewParticleButtonWithStyle(
		func() { println("无粒子按钮被点击了！") },
		"无粒子按钮",
		redStyle,
	)
	noParticleBtn.EnableParticle = false

	purpleStyle := tools.ParticleButtonStyle{
		BaseColor:  color.RGBA{R: 200, G: 100, B: 255, A: 255},
		UseGGFont:  true,
		GGFontType: "chinese",
	}
	purpleBtn := tools.NewParticleButtonWithStyle(
		func() { println("紫色按钮被点击了！") },
		"紫色按钮",
		purpleStyle,
	)

	dynamicStyle := tools.ParticleButtonStyle{
		BaseColor:    color.RGBA{R: 150, G: 150, B: 150, A: 255},
		AutoColorful: true,
	}
	dynamicBtn := tools.NewParticleButtonWithStyle(
		func() { println("动态按钮被点击了！") },
		"点击换色",
		dynamicStyle,
	)

	greenStyle := tools.ParticleButtonStyle{
		BaseColor:     color.RGBA{R: 143, G: 196, B: 0, A: 255},
		UseGGFont:     true,
		GGFontType:    "english",
		GGFontSize:    32,
		GGFontColor:   color.RGBA{0, 128, 0, 255},
		GGFontOffsetY: 10,
	}
	greenBtn := tools.NewParticleButtonWithStyle(
		func() { println("绿色按钮被点击了！") },
		"bro",
		greenStyle,
	)

	// 创建BorderButton演示
	borderBtn := tools.NewBorderButton(func(active bool) { println("边框按钮被点击了！") }, "边框按钮")

	// 反色样式边框按钮
	inverseBorderBtn := tools.NewInverseBorderButton(func(active bool) {
		println("反色边框按钮被点击了！激活状态:", active)
	}, "反色边框按钮")
	inverseBorderBtn.ContourWidthScale = 0.8
	inverseBorderBtn.ContourHeightScale = 0.7

	// 透明样式边框按钮
	transparentBorderBtn := tools.NewTransparentBorderButton(func(active bool) {
		println("透明边框按钮被点击了！激活状态:", active)
	}, "透明边框按钮")
	transparentBorderBtn.ContourWidthScale = 0.75
	transparentBorderBtn.ContourHeightScale = 0.65

	// 自定义颜色边框按钮
	customStyle := tools.BorderButtonStyle{
		DefaultColor:       color.RGBA{R: 255, G: 230, B: 200, A: 255},
		DefaultText:        color.RGBA{R: 150, G: 75, B: 0, A: 255},
		HoverColor:         color.RGBA{R: 255, G: 210, B: 170, A: 255},
		PressedColor:       color.RGBA{R: 255, G: 190, B: 140, A: 255},
		PressedText:        color.RGBA{R: 120, G: 60, B: 0, A: 255},
		ActiveColor:        color.RGBA{R: 200, G: 255, B: 220, A: 255},
		ActiveContour:      color.RGBA{R: 100, G: 220, B: 150, A: 255},
		ActiveText:         color.RGBA{R: 0, G: 150, B: 80, A: 255},
		BorderRadius:       12,
		ContourWidthScale:  0.9,
		ContourHeightScale: 0.8,
		ContourLineWidth:   2.5,
		UseGGFont:          true,
		GGFontType:         "chinese",
		GGFontSize:         16,
		GGFontColor:        color.RGBA{0, 0, 0, 255},
		GGFontOffsetY:      2,
	}
	customBorderBtn := tools.NewBorderButtonWithStyle(func(active bool) {
		println("自定义边框按钮被点击了！激活状态:", active)
	}, "自定义边框按钮", customStyle)

	// 设置每个按钮的固定尺寸
	redBtn.SetSize(220, 56)
	noParticleBtn.SetSize(220, 56)
	purpleBtn.SetSize(220, 56)
	dynamicBtn.SetSize(220, 56)
	greenBtn.SetSize(220, 56)
	borderBtn.SetSize(220, 56)
	inverseBorderBtn.SetSize(220, 56)
	transparentBorderBtn.SetSize(220, 56)
	customBorderBtn.SetSize(220, 56)

	// 打开高分辨率修复 感觉没有太大的用处
	//redBtn.SetHighDPI(true)
	//purpleBtn.SetHighDPI(true)

	// 创建按钮容器
	btnsVBox := container.NewVBox(
		widget.NewLabelWithStyle("粒子按钮演示", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(container.NewStack(redBtn)),
		container.NewCenter(container.NewStack(noParticleBtn)),
		container.NewCenter(container.NewStack(purpleBtn)),
		container.NewCenter(container.NewStack(dynamicBtn)),
		container.NewCenter(container.NewStack(greenBtn)),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("边框按钮演示", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(container.NewStack(borderBtn.WrapWithIsolationContainer())),
		container.NewCenter(container.NewStack(inverseBorderBtn.WrapWithIsolationContainer())),
		container.NewCenter(container.NewStack(transparentBorderBtn.WrapWithIsolationContainer())),
		container.NewCenter(container.NewStack(customBorderBtn.WrapWithIsolationContainer())),
	)

	// 启动粒子更新协程
	buttons := []*tools.ParticleButton{redBtn, noParticleBtn, purpleBtn, dynamicBtn, greenBtn}
	go func() {
		ticker := time.NewTicker(16 * time.Millisecond)
		defer ticker.Stop()

		for {
			<-ticker.C
			for _, btn := range buttons {
				if btn != nil {
					btn.UpdateParticles()
					if btn.IsAnimating() {
						fyne.Do(func() {
							btn.Refresh()
						})
					}
				}
			}
		}
	}()

	return container.NewPadded(btnsVBox)
}

// BuildTogglePage 构建开关页面
func BuildTogglePage() fyne.CanvasObject {
	return container.NewPadded(BuildToggleSwitchPanel())
}

// BuildCustomListPage 构建自定义列表页面
func BuildCustomListPage() fyne.CanvasObject {
	// 创建一个简单的自定义列表示例
	items := []tools.ListItem{
		{Title: "项目一", Icon: nil},
		{Title: "项目二", Icon: nil},
		{Title: "项目三", Icon: nil},
		{Title: "项目四", Icon: nil},
		{Title: "项目五", Icon: nil},
	}

	config := tools.CustomListConfig{
		Items:           items,
		OnSelected:      func(id int) { println("选中项目:", id) },
		InitialSelected: 0,
	}

	list := tools.NewCustomList(config)

	return container.NewPadded(
		container.NewVBox(
			widget.NewLabelWithStyle("自定义列表演示", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			list,
		),
	)
}

// BuildAboutPage 构建关于页面 - 简化版本
func BuildAboutPage() fyne.CanvasObject {
	return container.NewPadded(
		container.NewVBox(
			widget.NewLabelWithStyle("关于", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			container.NewCenter(
				widget.NewLabel("自定义UI演示程序"),
			),
			widget.NewSeparator(),
			container.NewCenter(
				widget.NewLabel("功能模块:"),
			),
			container.NewHBox(
				container.NewVBox(
					widget.NewLabel("• 粒子按钮系统"),
					widget.NewLabel("• 边框按钮系统"),
				),
				container.NewVBox(
					widget.NewLabel("• 动画开关控件"),
					widget.NewLabel("• 谷歌风格输入框"),
					widget.NewLabel("• 自定义列表组件"),
				),
			),
			widget.NewSeparator(),
			container.NewCenter(
				widget.NewLabel("版本: 1.0.0\n2025年12月"),
			),
		),
	)
}
