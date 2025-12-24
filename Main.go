// main.go
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"2025-12-18-ggAndPng/tools"
)

func main() {
	// 创建应用
	myApp := app.NewWithID("com.example.customui")
	myWindow := myApp.NewWindow("自定义UI演示")
	myWindow.Resize(fyne.NewSize(1600, 900))

	// 设置透明输入主题
	myApp.Settings().SetTheme(tools.NewInputTransparentTheme())

	// 创建菜单管理器
	menuManager := NewMenuManager(myWindow)

	// 设置窗口内容
	myWindow.SetContent(menuManager.GetContent())

	myWindow.ShowAndRun()
}
