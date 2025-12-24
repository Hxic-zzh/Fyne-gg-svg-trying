// main_menu.go
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MenuManager 菜单管理器
type MenuManager struct {
	window      fyne.Window
	currentPage fyne.CanvasObject
	content     *fyne.Container
	menu        *widget.List
}

// NewMenuManager 创建菜单管理器
func NewMenuManager(window fyne.Window) *MenuManager {
	m := &MenuManager{
		window:  window,
		content: container.NewStack(),
	}

	m.createMenu()
	return m
}

// 菜单项结构
type MenuItem struct {
	Title    string
	Icon     fyne.Resource
	PageFunc func() fyne.CanvasObject
}

// 创建菜单
func (m *MenuManager) createMenu() {
	// 定义菜单项
	menuItems := []MenuItem{
		{
			Title:    "主界面",
			Icon:     theme.HomeIcon(),
			PageFunc: BuildMainPage,
		},
		{
			Title:    "开关演示",
			Icon:     theme.CheckButtonCheckedIcon(),
			PageFunc: BuildTogglePage,
		},
		{
			Title:    "输入框演示",
			Icon:     theme.DocumentCreateIcon(),
			PageFunc: BuildGoogleInputPage,
		},
		{
			Title:    "复选框演示",
			Icon:     theme.CheckButtonIcon(),
			PageFunc: BuildCheckboxPage,
		},
		{
			Title:    "步骤标签页演示",
			Icon:     theme.NavigateNextIcon(),
			PageFunc: BuildSteapTabPage,
		},
		{
			Title:    "关于",
			Icon:     theme.InfoIcon(),
			PageFunc: BuildAboutPage,
		},
		{
			Title:    "性能测试",
			Icon:     theme.SettingsIcon(),
			PageFunc: BuildBenchmarkPage,
		},
	}

	// 创建菜单列表
	m.menu = widget.NewList(
		func() int {
			return len(menuItems)
		},
		func() fyne.CanvasObject {
			icon := widget.NewIcon(nil)
			label := widget.NewLabel("Template")
			return container.NewHBox(icon, label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			hbox := obj.(*fyne.Container)
			icon := hbox.Objects[0].(*widget.Icon)
			label := hbox.Objects[1].(*widget.Label)

			icon.SetResource(menuItems[id].Icon)
			label.SetText(menuItems[id].Title)
		},
	)

	// 设置菜单选中事件
	m.menu.OnSelected = func(id widget.ListItemID) {
		if id < len(menuItems) {
			m.SwitchToPage(menuItems[id].PageFunc())
		}
	}

	// 设置初始选中项
	m.menu.Select(0)
}

// SwitchToPage 切换到指定页面
func (m *MenuManager) SwitchToPage(page fyne.CanvasObject) {
	m.content.Objects = []fyne.CanvasObject{page}
	m.content.Refresh()
	m.currentPage = page
}

// GetContent 获取完整界面内容
func (m *MenuManager) GetContent() fyne.CanvasObject {
	// 创建菜单标题
	menuTitle := widget.NewLabelWithStyle("导航菜单", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	menuTitle.Resize(fyne.NewSize(200, 40))

	// 创建菜单容器
	menuContainer := container.NewBorder(
		menuTitle,
		container.NewVBox(
			widget.NewSeparator(),
			container.NewHBox(
				widget.NewIcon(theme.ComputerIcon()),
				widget.NewLabel("v1.0.0"),
			),
		),
		nil,
		nil,
		container.NewBorder(nil, nil, nil, nil, m.menu),
	)

	// 创建分割布局：左侧菜单，右侧内容
	split := container.NewHSplit(
		menuContainer,
		m.content,
	)
	split.SetOffset(0.18) // 菜单占18%宽度

	return container.NewPadded(split)
}
