package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"2025-12-18-ggAndPng/tools"
)

// BuildToggleSwitchPanel 创建6种动画风格的ToggleSwitch按钮演示区
func BuildToggleSwitchPanel() fyne.CanvasObject {
	labels := []string{
		"Effect 0: 默认",
		"Effect 1: 双横",
		"Effect 2: 双纵",
		"Effect 3: Y旋",
		"Effect 4: 矩形",
		"Effect 5: 矩遮",
	}

	// 定义不同的自定义配置
	customConfigs := []tools.SwitchConfig{
		// 配置1: 默认配置
		tools.DefaultConfig,

		// 配置2: 红绿主题 - 中文
		{
			YesLabel:      "开启",
			NoLabel:       "关闭",
			FontPath:      "ttf/chinese.ttf",                          // 中文字体
			YesColor:      color.RGBA{R: 76, G: 175, B: 80, A: 255},   // 绿色
			NoColor:       color.RGBA{R: 244, G: 67, B: 54, A: 255},   // 红色
			YesBgColor:    color.RGBA{R: 232, G: 245, B: 233, A: 255}, // 浅绿背景
			NoBgColor:     color.RGBA{R: 252, G: 228, B: 236, A: 255}, // 浅红背景
			TextColor:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色文字
			TextDarkColor: color.RGBA{R: 97, G: 97, B: 97, A: 255},    // 深灰色文字
			YesValue:      true,
			NoValue:       false,
		},

		// 配置3: 蓝黄主题 - 英文
		{
			YesLabel:      "ON",
			NoLabel:       "OFF",
			FontPath:      "ttf/toggle_switch.ttf",                    // 英文字体
			YesColor:      color.RGBA{R: 33, G: 150, B: 243, A: 255},  // 蓝色
			NoColor:       color.RGBA{R: 255, G: 193, B: 7, A: 255},   // 黄色
			YesBgColor:    color.RGBA{R: 227, G: 242, B: 253, A: 255}, // 浅蓝背景
			NoBgColor:     color.RGBA{R: 255, G: 248, B: 225, A: 255}, // 浅黄背景
			TextColor:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色文字
			TextDarkColor: color.RGBA{R: 66, G: 66, B: 66, A: 255},    // 深灰色文字
			YesValue:      true,
			NoValue:       false,
		},

		// 配置4: 紫色主题 - 中文
		{
			YesLabel:      "是",
			NoLabel:       "否",
			FontPath:      "ttf/chinese.ttf",                          // 中文字体
			YesColor:      color.RGBA{R: 156, G: 39, B: 176, A: 255},  // 紫色
			NoColor:       color.RGBA{R: 255, G: 152, B: 0, A: 255},   // 橙色
			YesBgColor:    color.RGBA{R: 243, G: 229, B: 245, A: 255}, // 浅紫背景
			NoBgColor:     color.RGBA{R: 255, G: 243, B: 224, A: 255}, // 浅橙背景
			TextColor:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色文字
			TextDarkColor: color.RGBA{R: 117, G: 117, B: 117, A: 255}, // 中灰色文字
			YesValue:      true,
			NoValue:       false,
		},

		// 配置5: 暗黑主题 - 中文
		{
			YesLabel:      "启用",
			NoLabel:       "禁用",
			FontPath:      "ttf/chinese.ttf",                          // 中文字体
			YesColor:      color.RGBA{R: 0, G: 200, B: 83, A: 255},    // 绿色
			NoColor:       color.RGBA{R: 255, G: 61, B: 0, A: 255},    // 红色
			YesBgColor:    color.RGBA{R: 40, G: 40, B: 40, A: 255},    // 深灰背景
			NoBgColor:     color.RGBA{R: 60, G: 60, B: 60, A: 255},    // 深灰背景
			TextColor:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色文字
			TextDarkColor: color.RGBA{R: 200, G: 200, B: 200, A: 255}, // 浅灰色文字
			YesValue:      true,
			NoValue:       false,
		},

		// 配置6: 简约主题 - 英文
		{
			YesLabel:      "Y",
			NoLabel:       "N",
			FontPath:      "ttf/toggle_switch.ttf",                    // 英文字体
			YesColor:      color.RGBA{R: 0, G: 150, B: 136, A: 255},   // 青色
			NoColor:       color.RGBA{R: 213, G: 0, B: 0, A: 255},     // 深红
			YesBgColor:    color.RGBA{R: 224, G: 242, B: 241, A: 255}, // 浅青背景
			NoBgColor:     color.RGBA{R: 255, G: 235, B: 238, A: 255}, // 浅红背景
			TextColor:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色文字
			TextDarkColor: color.RGBA{R: 96, G: 125, B: 139, A: 255},  // 蓝灰色文字
			YesValue:      true,
			NoValue:       false,
		},
	}

	var toggles []fyne.CanvasObject

	// 演示1: 每个效果使用不同的自定义配置
	for i := 0; i < 6; i++ {
		config := customConfigs[i%len(customConfigs)]

		ts := tools.NewToggleSwitch(false).
			SetEffect(tools.SwitchEffect(i)).
			SetConfig(config).
			SetSize(140, 65)

		localTs := ts // 捕获当前 ts
		localTs.OnChanged = func(_ bool) {
			println("开关值:", localTs.Value())
		}

		row := container.NewHBox(
			widget.NewLabel(labels[i]+" ("+config.YesLabel+"/"+config.NoLabel+")"),
			container.NewCenter(ts.WrapWithIsolationContainer()),
		)
		toggles = append(toggles, row)
	}

	// 演示2: 使用链式调用设置单个属性
	demoLabel := widget.NewLabel("演示链式调用设置属性:")

	// 使用链式调用自定义单个开关 - 中文
	customToggle1 := tools.NewToggleSwitch(true).
		SetEffect(tools.EffectSlide).
		SetYesLabel("开").
		SetNoLabel("关").
		SetFontPath("ttf/chinese.ttf").                       // 设置中文字体
		SetYesColor(color.RGBA{R: 0, G: 200, B: 83, A: 255}). // 绿色
		SetNoColor(color.RGBA{R: 255, G: 61, B: 0, A: 255}).  // 红色
		SetYesBgColor(color.RGBA{R: 232, G: 245, B: 233, A: 255}).
		SetNoBgColor(color.RGBA{R: 252, G: 228, B: 236, A: 255}).
		SetSize(160, 70)

	// 英文开关
	customToggle2 := tools.NewToggleSwitch(false).
		SetEffect(tools.EffectTwoBallSwap).
		SetYesLabel("ONLINE").
		SetNoLabel("OFFLINE").
		SetFontPath("ttf/toggle_switch.ttf").                   // 设置英文字体
		SetYesColor(color.RGBA{R: 3, G: 169, B: 244, A: 255}).  // 蓝色
		SetNoColor(color.RGBA{R: 158, G: 158, B: 158, A: 255}). // 灰色
		SetTextColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}).
		SetSize(180, 70)

	// 另一个中文开关示例
	customToggle4 := tools.NewToggleSwitch(false).
		SetEffect(tools.EffectProjectionFlip).
		SetYesLabel("激活").
		SetNoLabel("停用").
		SetFontPath("ttf/chinese.ttf").                         // 中文字体
		SetYesColor(color.RGBA{R: 156, G: 39, B: 176, A: 255}). // 紫色
		SetNoColor(color.RGBA{R: 255, G: 152, B: 0, A: 255}).   // 橙色
		SetYesBgColor(color.RGBA{R: 243, G: 229, B: 245, A: 255}).
		SetNoBgColor(color.RGBA{R: 255, G: 243, B: 224, A: 255}).
		SetSize(170, 65)

	panel := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("ToggleSwitch 6种动画风格演示："),
		widget.NewSeparator(),
	)

	panel.Objects = append(panel.Objects, toggles...)

	panel.Add(widget.NewSeparator())
	panel.Add(demoLabel)

	panel.Add(container.NewHBox(
		widget.NewLabel("自定义开关1 (中文):"),
		container.NewCenter(customToggle1.WrapWithIsolationContainer()),
	))

	panel.Add(container.NewHBox(
		widget.NewLabel("自定义开关2 (英文):"),
		container.NewCenter(customToggle2.WrapWithIsolationContainer()),
	))

	panel.Add(container.NewHBox(
		widget.NewLabel("自定义开关4 (中文):"),
		container.NewCenter(customToggle4.WrapWithIsolationContainer()),
	))

	return panel
}
