// main_check.go
package main

import (
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"2025-12-18-ggAndPng/tools"
)

// BuildCheckboxPage 构建复选框演示页面
func BuildCheckboxPage() fyne.CanvasObject {
	// 创建标题
	titleLabel := widget.NewLabelWithStyle("Material Design 复选框演示 (支持自定义动画颜色)",
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// SVG图标路径列表
	svgPaths := []string{
		"svg/1.svg",
		"svg/2.svg",
		"svg/3.svg",
		"svg/4.svg",
		"svg/5.svg",
		"svg/6.svg",
	}

	// 创建6个不同样式的复选框，模仿CSS示例
	checkboxes := []struct {
		name    string
		checked bool
		style   tools.MaterialCheckboxStyle
	}{
		{
			name:    "Discord",
			checked: false,
			style: tools.MaterialCheckboxStyle{
				TileWidth:    112, // 7rem = 112px
				TileHeight:   112,
				IconColor:    color.RGBA{114, 137, 218, 255}, // Discord蓝
				LabelColor:   color.RGBA{112, 112, 112, 255}, // #707070
				BorderColor:  tools.CheckboxColorDefault,
				BgColor:      color.White,
				CornerRadius: 8, // 0.5rem = 8px
				IconPath:     svgPaths[0],
				// 自定义动画颜色
				HoverColor:     color.RGBA{114, 137, 218, 100}, // Discord蓝半透明
				SelectedColor:  color.RGBA{114, 137, 218, 255}, // Discord蓝
				CircleColor:    color.RGBA{114, 137, 218, 255}, // Discord蓝
				CheckmarkColor: color.White,
			},
		},
		{
			name:    "Framer",
			checked: true, // 默认选中
			style: tools.MaterialCheckboxStyle{
				TileWidth:    112,
				TileHeight:   112,
				IconColor:    color.RGBA{0, 0, 0, 255},    // 黑色
				LabelColor:   color.RGBA{85, 85, 85, 255}, // 深灰色
				BorderColor:  color.RGBA{0, 0, 0, 255},    // 黑色
				BgColor:      color.White,
				CornerRadius: 8,
				IconPath:     svgPaths[1],
				// 自定义动画颜色
				HoverColor:     color.RGBA{0, 0, 0, 100}, // 黑色半透明
				SelectedColor:  color.RGBA{0, 0, 0, 255}, // 黑色
				CircleColor:    color.RGBA{0, 0, 0, 255}, // 黑色
				CheckmarkColor: color.White,
			},
		},
		{
			name:    "Sketch",
			checked: false,
			style: tools.MaterialCheckboxStyle{
				TileWidth:    112,
				TileHeight:   112,
				IconColor:    color.RGBA{253, 176, 34, 255}, // Sketch橙色
				LabelColor:   tools.CheckboxColorText,
				BorderColor:  tools.CheckboxColorDefault,
				BgColor:      color.White,
				CornerRadius: 8,
				IconPath:     svgPaths[2],
				// 自定义动画颜色
				HoverColor:     color.RGBA{253, 176, 34, 100}, // Sketch橙半透明
				SelectedColor:  color.RGBA{253, 176, 34, 255}, // Sketch橙
				CircleColor:    color.RGBA{253, 176, 34, 255}, // Sketch橙
				CheckmarkColor: color.White,
			},
		},
		{
			name:    "Instagram",
			checked: false,
			style: tools.MaterialCheckboxStyle{
				TileWidth:    112,
				TileHeight:   112,
				IconColor:    color.RGBA{225, 48, 108, 255}, // Instagram紫红
				LabelColor:   tools.CheckboxColorText,
				BorderColor:  tools.CheckboxColorDefault,
				BgColor:      color.White,
				CornerRadius: 8,
				IconPath:     svgPaths[3],
				// 自定义动画颜色
				HoverColor:     color.RGBA{225, 48, 108, 100}, // Instagram紫红半透明
				SelectedColor:  color.RGBA{225, 48, 108, 255}, // Instagram紫红
				CircleColor:    color.RGBA{225, 48, 108, 255}, // Instagram紫红
				CheckmarkColor: color.White,
			},
		},
		{
			name:    "Dribbble",
			checked: false,
			style: tools.MaterialCheckboxStyle{
				TileWidth:    112,
				TileHeight:   112,
				IconColor:    color.RGBA{234, 76, 137, 255}, // Dribbble粉红
				LabelColor:   tools.CheckboxColorText,
				BorderColor:  tools.CheckboxColorDefault,
				BgColor:      color.White,
				CornerRadius: 8,
				IconPath:     svgPaths[4],
				// 自定义动画颜色
				HoverColor:     color.RGBA{234, 76, 137, 100}, // Dribbble粉红半透明
				SelectedColor:  color.RGBA{234, 76, 137, 255}, // Dribbble粉红
				CircleColor:    color.RGBA{234, 76, 137, 255}, // Dribbble粉红
				CheckmarkColor: color.White,
			},
		},
		{
			name:    "Slack",
			checked: true, // 默认选中
			style: tools.MaterialCheckboxStyle{
				TileWidth:    112,
				TileHeight:   112,
				IconColor:    color.RGBA{74, 21, 75, 255}, // Slack紫色
				LabelColor:   color.RGBA{74, 21, 75, 255}, // Slack紫色
				BorderColor:  color.RGBA{74, 21, 75, 255}, // Slack紫色
				BgColor:      color.White,
				CornerRadius: 8,
				IconPath:     svgPaths[5],
				// 自定义动画颜色
				HoverColor:     color.RGBA{74, 21, 75, 100}, // Slack紫半透明
				SelectedColor:  color.RGBA{74, 21, 75, 255}, // Slack紫
				CircleColor:    color.RGBA{74, 21, 75, 255}, // Slack紫
				CheckmarkColor: color.White,
			},
		},
	}

	var checkboxWidgets []fyne.CanvasObject
	var checkboxInstances []*tools.MaterialCheckbox

	// 创建第一行复选框（使用不同的自定义动画颜色方案）
	for i, cb := range checkboxes[:3] {
		var style tools.MaterialCheckboxStyle
		switch i {
		case 0:
			// 绿色主题
			style = tools.MaterialCheckboxStyle{
				TileWidth:    112,
				TileHeight:   112,
				IconColor:    color.RGBA{46, 204, 113, 255}, // 绿色
				LabelColor:   color.RGBA{46, 204, 113, 255}, // 绿色
				BorderColor:  color.RGBA{39, 174, 96, 255},  // 深绿色
				BgColor:      color.White,
				CornerRadius: 8,
				IconPath:     svgPaths[0],
				// 自定义动画颜色
				HoverColor:     color.RGBA{46, 204, 113, 100}, // 绿色半透明
				SelectedColor:  color.RGBA{46, 204, 113, 255}, // 绿色
				CircleColor:    color.RGBA{46, 204, 113, 255}, // 绿色
				CheckmarkColor: color.White,
			}
		case 1:
			// 橙色主题
			style = tools.MaterialCheckboxStyle{
				TileWidth:    112,
				TileHeight:   112,
				IconColor:    color.RGBA{243, 156, 18, 255}, // 橙色
				LabelColor:   color.RGBA{243, 156, 18, 255}, // 橙色
				BorderColor:  color.RGBA{211, 84, 0, 255},   // 深橙色
				BgColor:      color.White,
				CornerRadius: 8,
				IconPath:     svgPaths[1],
				// 自定义动画颜色
				HoverColor:     color.RGBA{243, 156, 18, 100}, // 橙色半透明
				SelectedColor:  color.RGBA{243, 156, 18, 255}, // 橙色
				CircleColor:    color.RGBA{243, 156, 18, 255}, // 橙色
				CheckmarkColor: color.White,
			}
		case 2:
			// 紫色主题
			style = tools.MaterialCheckboxStyle{
				TileWidth:    112,
				TileHeight:   112,
				IconColor:    color.RGBA{155, 89, 182, 255}, // 紫色
				LabelColor:   color.RGBA{155, 89, 182, 255}, // 紫色
				BorderColor:  color.RGBA{142, 68, 173, 255}, // 深紫色
				BgColor:      color.White,
				CornerRadius: 8,
				IconPath:     svgPaths[2],
				// 自定义动画颜色
				HoverColor:     color.RGBA{155, 89, 182, 100}, // 紫色半透明
				SelectedColor:  color.RGBA{155, 89, 182, 255}, // 紫色
				CircleColor:    color.RGBA{155, 89, 182, 255}, // 紫色
				CheckmarkColor: color.White,
			}
		}

		checkbox := tools.NewMaterialCheckbox(cb.name, cb.checked, style.TileWidth, style.TileHeight)
		checkbox.SetStyle(style)
		checkbox.SetFontPath("ttf/english.ttf")
		currentName := cb.name
		checkbox.OnChanged = func(checked bool) {
			println("复选框状态变化:", currentName, "checked:", checked)
		}
		checkboxInstances = append(checkboxInstances, checkbox)
		checkboxWidgets = append(checkboxWidgets, checkbox.WrapWithIsolationContainer())
	}

	// 创建第二行复选框
	for i, cb := range checkboxes[3:] {
		var style tools.MaterialCheckboxStyle
		if i == 0 {
			// 青色主题
			style = tools.MaterialCheckboxStyle{
				TileWidth:    112,
				TileHeight:   112,
				IconColor:    color.RGBA{26, 188, 156, 255}, // 青色
				LabelColor:   color.RGBA{26, 188, 156, 255}, // 青色
				BorderColor:  color.RGBA{22, 160, 133, 255}, // 深青色
				BgColor:      color.White,
				CornerRadius: 8,
				IconPath:     svgPaths[3],
				// 自定义动画颜色
				HoverColor:     color.RGBA{26, 188, 156, 100}, // 青色半透明
				SelectedColor:  color.RGBA{26, 188, 156, 255}, // 青色
				CircleColor:    color.RGBA{26, 188, 156, 255}, // 青色
				CheckmarkColor: color.White,
			}
		} else {
			// 使用定义的样式
			style = cb.style
		}

		checkbox := tools.NewMaterialCheckbox(cb.name, cb.checked, style.TileWidth, style.TileHeight)
		checkbox.SetStyle(style)
		checkbox.SetFontPath("ttf/english.ttf")
		currentName := cb.name
		checkbox.OnChanged = func(checked bool) {
			println("复选框状态变化:", currentName, "checked:", checked)
		}
		checkboxInstances = append(checkboxInstances, checkbox)
		checkboxWidgets = append(checkboxWidgets, checkbox.WrapWithIsolationContainer())
	}

	// 创建两行布局
	firstRow := container.NewHBox()
	for i := 0; i < 3; i++ {
		firstRow.Add(checkboxWidgets[i])
		if i < 2 {
			firstRow.Add(widget.NewLabel("")) // 间距
		}
	}

	secondRow := container.NewHBox()
	for i := 3; i < 6; i++ {
		secondRow.Add(checkboxWidgets[i])
		if i < 5 {
			secondRow.Add(widget.NewLabel("")) // 间距
		}
	}

	// 创建演示控件区域
	demoSection := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("CSS风格的复选框演示 - 每个都有自定义动画颜色"),
		widget.NewSeparator(),
		container.NewCenter(firstRow),
		container.NewCenter(secondRow),
	)

	// 创建控制面板
	controlPanel := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("控制面板"),
		widget.NewSeparator(),
	)

	// 添加单选按钮组 - 现在可以自定义动画颜色
	radioGroup := widget.NewRadioGroup([]string{
		"默认样式 (蓝色动画)",
		"红色动画主题",
		"绿色动画主题",
		"紫色动画主题",
		"金色动画主题",
		"彩虹动画主题",
	}, func(selected string) {
		for i, cb := range checkboxInstances {
			var newStyle tools.MaterialCheckboxStyle

			// 获取当前样式的非动画属性
			currentStyle := cb.Style

			// 根据选择应用不同的动画颜色方案
			switch selected {
			case "红色动画主题":
				newStyle = tools.MaterialCheckboxStyle{
					TileWidth:    currentStyle.TileWidth,
					TileHeight:   currentStyle.TileHeight,
					IconColor:    currentStyle.IconColor,
					LabelColor:   currentStyle.LabelColor,
					BorderColor:  currentStyle.BorderColor,
					BgColor:      currentStyle.BgColor,
					ShadowColor:  currentStyle.ShadowColor,
					CornerRadius: currentStyle.CornerRadius,
					IconSize:     currentStyle.IconSize,
					FontSize:     currentStyle.FontSize,
					IconPath:     currentStyle.IconPath,
					// 红色动画
					HoverColor:     color.RGBA{255, 100, 100, 100},
					SelectedColor:  color.RGBA{255, 0, 0, 255},
					CircleColor:    color.RGBA{255, 0, 0, 255},
					CheckmarkColor: color.White,
				}

			case "绿色动画主题":
				newStyle = tools.MaterialCheckboxStyle{
					TileWidth:    currentStyle.TileWidth,
					TileHeight:   currentStyle.TileHeight,
					IconColor:    currentStyle.IconColor,
					LabelColor:   currentStyle.LabelColor,
					BorderColor:  currentStyle.BorderColor,
					BgColor:      currentStyle.BgColor,
					ShadowColor:  currentStyle.ShadowColor,
					CornerRadius: currentStyle.CornerRadius,
					IconSize:     currentStyle.IconSize,
					FontSize:     currentStyle.FontSize,
					IconPath:     currentStyle.IconPath,
					// 绿色动画
					HoverColor:     color.RGBA{100, 255, 100, 100},
					SelectedColor:  color.RGBA{0, 200, 0, 255},
					CircleColor:    color.RGBA{0, 200, 0, 255},
					CheckmarkColor: color.White,
				}

			case "紫色动画主题":
				newStyle = tools.MaterialCheckboxStyle{
					TileWidth:    currentStyle.TileWidth,
					TileHeight:   currentStyle.TileHeight,
					IconColor:    currentStyle.IconColor,
					LabelColor:   currentStyle.LabelColor,
					BorderColor:  currentStyle.BorderColor,
					BgColor:      currentStyle.BgColor,
					ShadowColor:  currentStyle.ShadowColor,
					CornerRadius: currentStyle.CornerRadius,
					IconSize:     currentStyle.IconSize,
					FontSize:     currentStyle.FontSize,
					IconPath:     currentStyle.IconPath,
					// 紫色动画
					HoverColor:     color.RGBA{180, 100, 255, 100},
					SelectedColor:  color.RGBA{150, 0, 255, 255},
					CircleColor:    color.RGBA{150, 0, 255, 255},
					CheckmarkColor: color.White,
				}

			case "金色动画主题":
				newStyle = tools.MaterialCheckboxStyle{
					TileWidth:    currentStyle.TileWidth,
					TileHeight:   currentStyle.TileHeight,
					IconColor:    currentStyle.IconColor,
					LabelColor:   currentStyle.LabelColor,
					BorderColor:  currentStyle.BorderColor,
					BgColor:      currentStyle.BgColor,
					ShadowColor:  currentStyle.ShadowColor,
					CornerRadius: currentStyle.CornerRadius,
					IconSize:     currentStyle.IconSize,
					FontSize:     currentStyle.FontSize,
					IconPath:     currentStyle.IconPath,
					// 金色动画
					HoverColor:     color.RGBA{255, 215, 0, 100},
					SelectedColor:  color.RGBA{255, 215, 0, 255},
					CircleColor:    color.RGBA{255, 215, 0, 255},
					CheckmarkColor: color.White,
				}

			case "彩虹动画主题":
				// 每个复选框不同颜色
				rainbowColors := []color.RGBA{
					{255, 0, 0, 255},   // 红
					{255, 127, 0, 255}, // 橙
					{255, 255, 0, 255}, // 黄
					{0, 255, 0, 255},   // 绿
					{0, 0, 255, 255},   // 蓝
					{75, 0, 130, 255},  // 靛
				}

				if i < len(rainbowColors) {
					newStyle = tools.MaterialCheckboxStyle{
						TileWidth:    currentStyle.TileWidth,
						TileHeight:   currentStyle.TileHeight,
						IconColor:    currentStyle.IconColor,
						LabelColor:   currentStyle.LabelColor,
						BorderColor:  currentStyle.BorderColor,
						BgColor:      currentStyle.BgColor,
						ShadowColor:  currentStyle.ShadowColor,
						CornerRadius: currentStyle.CornerRadius,
						IconSize:     currentStyle.IconSize,
						FontSize:     currentStyle.FontSize,
						IconPath:     currentStyle.IconPath,
						// 彩虹颜色
						HoverColor: color.RGBA{
							rainbowColors[i].R,
							rainbowColors[i].G,
							rainbowColors[i].B,
							100,
						},
						SelectedColor:  rainbowColors[i],
						CircleColor:    rainbowColors[i],
						CheckmarkColor: color.White,
					}
				} else {
					// 如果索引超出，使用默认蓝色
					newStyle = currentStyle
				}

			default: // "默认样式 (蓝色动画)"
				newStyle = tools.MaterialCheckboxStyle{
					TileWidth:    currentStyle.TileWidth,
					TileHeight:   currentStyle.TileHeight,
					IconColor:    currentStyle.IconColor,
					LabelColor:   currentStyle.LabelColor,
					BorderColor:  currentStyle.BorderColor,
					BgColor:      currentStyle.BgColor,
					ShadowColor:  currentStyle.ShadowColor,
					CornerRadius: currentStyle.CornerRadius,
					IconSize:     currentStyle.IconSize,
					FontSize:     currentStyle.FontSize,
					IconPath:     currentStyle.IconPath,
					// 默认蓝色动画
					HoverColor:     tools.CheckboxColorHover,
					SelectedColor:  tools.CheckboxColorSelected,
					CircleColor:    tools.CheckboxColorSelected,
					CheckmarkColor: color.White,
				}
			}

			cb.SetStyle(newStyle)
		}
	})
	radioGroup.SetSelected("默认样式 (蓝色动画)")

	// 添加全选/取消全选按钮
	selectAllBtn := widget.NewButton("全选", func() {
		for _, cb := range checkboxInstances {
			cb.SetChecked(true)
		}
	})

	deselectAllBtn := widget.NewButton("取消全选", func() {
		for _, cb := range checkboxInstances {
			cb.SetChecked(false)
		}
	})

	// 添加重置样式按钮
	resetStyleBtn := widget.NewButton("重置样式", func() {
		radioGroup.SetSelected("默认样式 (蓝色动画)")
		for i, cb := range checkboxInstances {
			// 重置为初始样式
			cb.SetStyle(checkboxes[i].style)
		}
	})

	// 添加单个颜色自定义按钮
	customColorBtn := widget.NewButton("自定义颜色演示", func() {
		// 为每个复选框设置不同的自定义动画颜色
		customColors := []struct {
			hover    color.RGBA
			selected color.RGBA
			circle   color.RGBA
		}{
			// 柔和色调
			{color.RGBA{200, 230, 255, 100}, color.RGBA{100, 180, 255, 255}, color.RGBA{100, 180, 255, 255}}, // 天蓝
			{color.RGBA{255, 230, 200, 100}, color.RGBA{255, 180, 100, 255}, color.RGBA{255, 180, 100, 255}}, // 橙黄
			{color.RGBA{230, 200, 255, 100}, color.RGBA{180, 100, 255, 255}, color.RGBA{180, 100, 255, 255}}, // 粉紫
			{color.RGBA{200, 255, 230, 100}, color.RGBA{100, 255, 180, 255}, color.RGBA{100, 255, 180, 255}}, // 薄荷绿
			{color.RGBA{255, 200, 230, 100}, color.RGBA{255, 100, 180, 255}, color.RGBA{255, 100, 180, 255}}, // 粉红
			{color.RGBA{230, 255, 200, 100}, color.RGBA{180, 255, 100, 255}, color.RGBA{180, 255, 100, 255}}, // 浅绿
		}

		for i, cb := range checkboxInstances {
			if i < len(customColors) {
				currentStyle := cb.Style
				newStyle := tools.MaterialCheckboxStyle{
					TileWidth:      currentStyle.TileWidth,
					TileHeight:     currentStyle.TileHeight,
					IconColor:      currentStyle.IconColor,
					LabelColor:     currentStyle.LabelColor,
					BorderColor:    currentStyle.BorderColor,
					BgColor:        currentStyle.BgColor,
					ShadowColor:    currentStyle.ShadowColor,
					CornerRadius:   currentStyle.CornerRadius,
					IconSize:       currentStyle.IconSize,
					FontSize:       currentStyle.FontSize,
					IconPath:       currentStyle.IconPath,
					HoverColor:     customColors[i].hover,
					SelectedColor:  customColors[i].selected,
					CircleColor:    customColors[i].circle,
					CheckmarkColor: color.White,
				}
				cb.SetStyle(newStyle)
			}
		}
	})

	// 添加选中状态显示
	statusLabel := widget.NewLabel("选中状态: 2/6")

	// 更新状态显示的定时器
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			checkedCount := 0
			for _, cb := range checkboxInstances {
				if cb.IsChecked() {
					checkedCount++
				}
			}
			finalCount := checkedCount
			fyne.Do(func() {
				statusLabel.SetText("选中状态: " +
					strconv.Itoa(finalCount) + "/6")
			})
		}
	}()

	// 添加说明文本 - 更新为反映自定义动画颜色功能
	instructions := widget.NewLabel("说明:\n" +
		"• 点击卡片切换选中状态\n" +
		"• Hover时显示阴影和选中标记轮廓\n" +
		"• 选中时边框、图标、文字变为选中颜色\n" +
		"• 每个复选框都有自定义动画颜色\n" +
		"• 使用右侧控制面板切换动画颜色主题")

	controlPanel.Add(radioGroup)
	controlPanel.Add(container.NewHBox(
		selectAllBtn,
		deselectAllBtn,
	))
	controlPanel.Add(resetStyleBtn)
	controlPanel.Add(customColorBtn)
	controlPanel.Add(widget.NewSeparator())
	controlPanel.Add(statusLabel)
	controlPanel.Add(widget.NewSeparator())
	controlPanel.Add(instructions)

	// 主容器：左侧演示，右侧控制面板
	mainContainer := container.NewHSplit(
		container.NewPadded(demoSection),
		container.NewPadded(controlPanel),
	)
	mainContainer.SetOffset(0.7) // 演示区域占70%

	return container.NewPadded(
		container.NewVBox(
			titleLabel,
			mainContainer,
		),
	)
}
