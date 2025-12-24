package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"2025-12-18-ggAndPng/tools"
)

// BuildSteapTabPage 构建步骤标签页演示界面
func BuildSteapTabPage() fyne.CanvasObject {
	// 创建Tab项
	items := []*tools.TabItem{
		{
			ID:       "step1",
			Title:    "第一步",
			IconPath: "svg/1.svg",
			Content: container.NewCenter(
				widget.NewLabelWithStyle("第一步内容区域\n这是第一个步骤的页面",
					fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			),
			Enabled: true,
		},
		{
			ID:       "step2",
			Title:    "第二步",
			IconPath: "svg/2.svg",
			Content: container.NewCenter(
				widget.NewLabelWithStyle("第二步内容区域\n这是第二个步骤的页面\n包含更多信息",
					fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			),
			Enabled: true,
		},
		{
			ID:       "step3",
			Title:    "第三步",
			IconPath: "svg/3.svg",
			Content: container.NewCenter(
				widget.NewLabelWithStyle("第三步内容区域\n这是最后一个步骤的页面\n已完成所有步骤",
					fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			),
			Enabled: true,
		},
	}

	// 验证文件是否存在（创建演示前检查）
	for _, item := range items {
		if _, err := os.Stat(item.IconPath); os.IsNotExist(err) {
			return container.NewCenter(
				widget.NewLabel(fmt.Sprintf("SVG文件不存在: %s\n请创建svg目录并放入1.svg, 2.svg, 3.svg", item.IconPath)),
			)
		}
	}

	// 创建Tab组件
	stepTabs, err := tools.NewStepTabs(items)
	if err != nil {
		return container.NewCenter(
			widget.NewLabel("创建StepTabs失败: " + err.Error()),
		)
	}

	// 使用物理隔离容器包裹StepTabs
	wrappedStepTabs := stepTabs.WrapWithIsolationContainer()

	// 创建内容容器（使用Stack实现页面切换）
	contentStack := container.NewStack()
	if len(items) > 0 {
		contentStack.Add(items[0].Content)
	}

	// 状态显示标签
	statusLabel := widget.NewLabel("当前步骤: 第一步")
	statusLabel.Alignment = fyne.TextAlignCenter

	// Tab切换回调
	stepTabs.OnChanged = func(index int, id string) {
		// 更新内容
		contentStack.RemoveAll()
		contentStack.Add(items[index].Content)
		contentStack.Refresh()

		// 更新状态显示
		statusLabel.SetText(fmt.Sprintf("当前步骤: %s (ID: %s)", items[index].Title, id))
	}

	// 自定义样式
	customStyle := tools.DefaultStyle()
	customStyle.Width = 700
	customStyle.Height = 130
	customStyle.CircleSize = 32
	customStyle.IconSize = 26
	customStyle.Spacing = 130
	customStyle.TextOffsetY = 42
	customStyle.IndicatorOffsetY = 58

	stepTabs.SetStyle(customStyle)

	// 自定义颜色
	customColors := tools.DefaultColors()
	customColors.Active = color.RGBA{52, 152, 219, 255}      // 蓝色主题
	customColors.Normal = color.RGBA{189, 195, 199, 255}     // 浅灰色
	customColors.Line = color.RGBA{236, 240, 241, 255}       // 更浅的灰色线条
	customColors.Background = color.RGBA{250, 250, 250, 255} // 浅灰背景

	stepTabs.SetColors(customColors)

	// 创建简单的控制面板
	controlPanel := container.NewVBox(
		widget.NewLabelWithStyle("步骤指示器演示",
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		statusLabel,
	)

	// 组装完整界面 - 使用包裹后的StepTabs
	mainContent := container.NewBorder(
		controlPanel,    // 顶部：控制面板
		wrappedStepTabs, // 底部：使用物理隔离后的步骤指示器
		nil,             // 左边
		nil,             // 右边
		contentStack,    // 中间：内容区域
	)

	return container.NewPadded(mainContent)
}
