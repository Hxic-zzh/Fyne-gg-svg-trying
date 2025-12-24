/*
	                                                   $$\     $$\
	                                                   $$ |    \__|
	$$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\  $$$$$$\   $$\  $$$$$$\  $$$$$$$\

$$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\  \____$$\ \_$$  _|  $$ |$$  __$$\ $$  __$$\
$$ /  $$ |$$ /  $$ |$$$$$$$$ |$$ |  \__| $$$$$$$ |  $$ |    $$ |$$ /  $$ |$$ |  $$ |
$$ |  $$ |$$ |  $$ |$$   ____|$$ |      $$  __$$ |  $$ |$$\ $$ |$$ |  $$ |$$ |  $$ |
\$$$$$$  |$$$$$$$  |\$$$$$$$\ $$ |      \$$$$$$$ |  \$$$$  |$$ |\$$$$$$  |$$ |  $$ |

	\______/ $$  ____/  \_______|\__|       \_______|   \____/ \__| \______/ \__|  \__|
	         $$ |
	         $$ |
	         \__|
*/
package tools

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// InputTransparentTheme 只让输入框透明的主题
type InputTransparentTheme struct {
	defaultTheme fyne.Theme
}

func NewInputTransparentTheme() fyne.Theme {
	return &InputTransparentTheme{
		defaultTheme: theme.DefaultTheme(),
	}
}

func (t *InputTransparentTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// 只修改输入框相关颜色，其他保持不变
	switch name {
	case theme.ColorNameInputBackground:
		return color.Transparent
	case theme.ColorNameInputBorder:
		return color.Transparent
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 153, G: 153, B: 153, A: 180}
	default:
		return t.defaultTheme.Color(name, variant)
	}
}

func (t *InputTransparentTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.defaultTheme.Font(style)
}

func (t *InputTransparentTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.defaultTheme.Icon(name)
}

func (t *InputTransparentTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.defaultTheme.Size(name)
}
