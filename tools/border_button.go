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
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
)

// BorderButton 边框按钮（模仿CSS default按钮样式，更像开关）
type BorderButton struct {
	widget.BaseWidget
	OnToggle func(active bool) // 开关状态改变回调

	// 按钮配置
	Text          string
	Width, Height float32

	// 颜色配置
	DefaultColor   color.RGBA // 默认背景色
	DefaultContour color.RGBA // 默认轮廓色
	DefaultText    color.RGBA // 默认文字色
	HoverColor     color.RGBA // 悬停背景色
	PressedColor   color.RGBA // 按下背景色
	PressedContour color.RGBA // 按下轮廓色
	PressedText    color.RGBA // 按下文字色
	ActiveColor    color.RGBA // 激活/常亮状态背景色
	ActiveContour  color.RGBA // 激活状态轮廓色
	ActiveText     color.RGBA // 激活状态文字色

	// 状态
	isActive  bool // 是否激活（常亮状态）
	isHovered bool
	isPressed bool

	// 样式配置
	BorderRadius       float64 // 圆角半径
	PaddingLeft        int     // 左内边距
	PaddingRight       int     // 右内边距
	PaddingTop         int     // 上内边距
	PaddingBottom      int     // 下内边距
	ContourScale       float64 // 轮廓比例（相对于按钮尺寸，默认0.85）
	ContourLineWidth   float64 // 轮廓线宽度
	ContourWidthScale  float64 // 轮廓宽比例（相对于按钮宽度，默认0.85）
	ContourHeightScale float64 // 轮廓高比例（相对于按钮高度，默认0.85）

	// 字体配置
	UseGGFont     bool        // 是否使用gg字体渲染
	GGFontType    string      // 字体类型：chinese/english
	GGFontSize    float64     // 字体大小
	GGFontColor   color.Color // 字体颜色
	GGFontOffsetX float64     // 字体X偏移
	GGFontOffsetY float64     // 字体Y偏移

	// 动画控制
	stopAnimation chan bool
}

// 实现必要的接口
var _ fyne.Widget = (*BorderButton)(nil)
var _ fyne.Tappable = (*BorderButton)(nil)
var _ desktop.Hoverable = (*BorderButton)(nil)

// BorderButtonStyle CSS-like样式配置
type BorderButtonStyle struct {
	// 默认状态（初始状态）
	DefaultColor   color.RGBA // 默认背景色
	DefaultContour color.RGBA // 默认轮廓色
	DefaultText    color.RGBA // 默认文字色

	// 悬停状态
	HoverColor color.RGBA // 悬停背景色

	// 按下状态（点击时）
	PressedColor   color.RGBA // 按下背景色
	PressedContour color.RGBA // 按下轮廓色
	PressedText    color.RGBA // 按下文字色

	// 激活状态（常亮）
	ActiveColor   color.RGBA // 激活背景色
	ActiveContour color.RGBA // 激活轮廓色
	ActiveText    color.RGBA // 激活文字色

	// 尺寸和形状
	BorderRadius       float64 // 圆角半径
	ContourScale       float64 // 轮廓比例（相对于按钮尺寸，默认0.85）
	ContourLineWidth   float64 // 轮廓线宽度
	ContourWidthScale  float64 // 轮廓宽比例（相对于按钮宽度，默认0.85）
	ContourHeightScale float64 // 轮廓高比例（相对于按钮高度，默认0.85）

	// 内边距
	PaddingLeft   int // 左内边距
	PaddingRight  int // 右内边距
	PaddingTop    int // 上内边距
	PaddingBottom int // 下内边距

	// 字体配置
	UseGGFont     bool        // 是否使用gg字体渲染
	GGFontType    string      // 字体类型：chinese/english
	GGFontSize    float64     // 字体大小
	GGFontColor   color.Color // 字体颜色
	GGFontOffsetX float64     // 字体X偏移
	GGFontOffsetY float64     // 字体Y偏移
}

// NewBorderButton 创建默认样式的边框按钮（开关式）
func NewBorderButton(onToggle func(active bool), text string) *BorderButton {
	// 默认状态（初始状态）
	defaultColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}   // 浅灰背景
	defaultContour := color.RGBA{R: 180, G: 180, B: 180, A: 255} // 灰色轮廓
	defaultText := color.RGBA{R: 51, G: 51, B: 51, A: 255}       // #333文字

	// 悬停状态
	hoverColor := color.RGBA{R: 220, G: 220, B: 220, A: 255} // 稍深灰色

	// 按下状态（点击时）
	pressedColor := color.RGBA{R: 200, G: 200, B: 200, A: 255}   // 更深灰色
	pressedContour := color.RGBA{R: 150, G: 150, B: 150, A: 255} // 深灰色轮廓
	pressedText := color.RGBA{R: 30, G: 30, B: 30, A: 255}       // 更深的文字

	// 激活状态（常亮）
	activeColor := color.RGBA{R: 200, G: 230, B: 255, A: 255}   // 浅蓝色常亮
	activeContour := color.RGBA{R: 100, G: 180, B: 255, A: 255} // 蓝色轮廓
	activeText := color.RGBA{R: 0, G: 100, B: 200, A: 255}      // 深蓝色文字

	return NewBorderButtonWithAllColors(onToggle, text,
		defaultColor, defaultContour, defaultText,
		hoverColor,
		pressedColor, pressedContour, pressedText,
		activeColor, activeContour, activeText)
}

// NewInverseBorderButton 创建反色样式的边框按钮
func NewInverseBorderButton(onToggle func(active bool), text string) *BorderButton {
	// 默认状态（初始状态）- 深色版
	defaultColor := color.RGBA{R: 60, G: 60, B: 60, A: 255}      // 深灰背景
	defaultContour := color.RGBA{R: 100, G: 100, B: 100, A: 255} // 灰色轮廓
	defaultText := color.RGBA{R: 220, G: 220, B: 220, A: 255}    // 浅灰色文字

	// 悬停状态
	hoverColor := color.RGBA{R: 80, G: 80, B: 80, A: 255} // 稍亮灰色

	// 按下状态（点击时）
	pressedColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}   // 更亮灰色
	pressedContour := color.RGBA{R: 120, G: 120, B: 120, A: 255} // 亮灰色轮廓
	pressedText := color.RGBA{R: 240, G: 240, B: 240, A: 255}    // 更亮的文字

	// 激活状态（常亮）
	activeColor := color.RGBA{R: 100, G: 150, B: 255, A: 255}  // 蓝色常亮
	activeContour := color.RGBA{R: 60, G: 120, B: 220, A: 255} // 深蓝色轮廓
	activeText := color.RGBA{R: 255, G: 255, B: 255, A: 255}   // 白色文字

	return NewBorderButtonWithAllColors(onToggle, text,
		defaultColor, defaultContour, defaultText,
		hoverColor,
		pressedColor, pressedContour, pressedText,
		activeColor, activeContour, activeText)
}

// NewTransparentBorderButton 创建透明样式的边框按钮
func NewTransparentBorderButton(onToggle func(active bool), text string) *BorderButton {
	// 默认状态（初始状态）- 透明版
	defaultColor := color.RGBA{R: 255, G: 255, B: 255, A: 0}     // 透明背景
	defaultContour := color.RGBA{R: 200, G: 200, B: 200, A: 100} // 半透明轮廓
	defaultText := color.RGBA{R: 80, G: 80, B: 80, A: 255}       // 灰色文字

	// 悬停状态
	hoverColor := color.RGBA{R: 240, G: 240, B: 240, A: 150} // 半透明白色

	// 按下状态（点击时）
	pressedColor := color.RGBA{R: 225, G: 225, B: 225, A: 180}   // 半透明深白色
	pressedContour := color.RGBA{R: 170, G: 170, B: 170, A: 150} // 半透明深灰色轮廓
	pressedText := color.RGBA{R: 50, G: 50, B: 50, A: 255}       // 深灰色文字

	// 激活状态（常亮）
	activeColor := color.RGBA{R: 200, G: 230, B: 255, A: 180}   // 半透明蓝色常亮
	activeContour := color.RGBA{R: 100, G: 180, B: 255, A: 200} // 半透明蓝色轮廓
	activeText := color.RGBA{R: 0, G: 100, B: 200, A: 255}      // 深蓝色文字

	return NewBorderButtonWithAllColors(onToggle, text,
		defaultColor, defaultContour, defaultText,
		hoverColor,
		pressedColor, pressedContour, pressedText,
		activeColor, activeContour, activeText)
}

// NewBorderButtonWithStyle 使用样式配置创建按钮
func NewBorderButtonWithStyle(onToggle func(active bool), text string, style BorderButtonStyle) *BorderButton {
	btn := &BorderButton{
		OnToggle: onToggle,
		Text:     text,
		Width:    120,
		Height:   36,

		// 状态颜色
		DefaultColor:   style.DefaultColor,
		DefaultContour: style.DefaultContour,
		DefaultText:    style.DefaultText,
		HoverColor:     style.HoverColor,
		PressedColor:   style.PressedColor,
		PressedContour: style.PressedContour,
		PressedText:    style.PressedText,
		ActiveColor:    style.ActiveColor,
		ActiveContour:  style.ActiveContour,
		ActiveText:     style.ActiveText,

		isActive:  false,
		isHovered: false,
		isPressed: false,

		// 尺寸和形状
		BorderRadius:       style.BorderRadius,
		ContourScale:       style.ContourScale,
		ContourLineWidth:   style.ContourLineWidth,
		ContourWidthScale:  style.ContourWidthScale,
		ContourHeightScale: style.ContourHeightScale,

		// 内边距
		PaddingLeft:   style.PaddingLeft,
		PaddingRight:  style.PaddingRight,
		PaddingTop:    style.PaddingTop,
		PaddingBottom: style.PaddingBottom,

		// 字体配置
		UseGGFont:     style.UseGGFont,
		GGFontType:    style.GGFontType,
		GGFontSize:    style.GGFontSize,
		GGFontColor:   style.GGFontColor,
		GGFontOffsetX: style.GGFontOffsetX,
		GGFontOffsetY: style.GGFontOffsetY,

		stopAnimation: make(chan bool),
	}

	// 设置默认值
	if btn.BorderRadius <= 0 {
		btn.BorderRadius = 8
	}
	if btn.ContourWidthScale <= 0 {
		btn.ContourWidthScale = 0.85 // 默认85%
	}
	if btn.ContourHeightScale <= 0 {
		btn.ContourHeightScale = 0.85 // 默认85%
	}
	if btn.ContourLineWidth <= 0 {
		btn.ContourLineWidth = 2.0
	}
	if btn.PaddingLeft <= 0 {
		btn.PaddingLeft = 12
	}
	if btn.PaddingRight <= 0 {
		btn.PaddingRight = 12
	}
	if btn.PaddingTop <= 0 {
		btn.PaddingTop = 6
	}
	if btn.PaddingBottom <= 0 {
		btn.PaddingBottom = 6
	}
	if btn.GGFontType == "" {
		btn.GGFontType = "chinese"
	}
	if btn.GGFontSize <= 0 {
		btn.GGFontSize = 14
	}
	if btn.GGFontColor == nil {
		btn.GGFontColor = btn.DefaultText
	}

	// 设置颜色默认值
	if btn.DefaultColor.A == 0 && btn.DefaultColor.R == 0 && btn.DefaultColor.G == 0 && btn.DefaultColor.B == 0 {
		btn.DefaultColor = color.RGBA{R: 240, G: 240, B: 240, A: 255}
	}
	if btn.DefaultContour.A == 0 && btn.DefaultContour.R == 0 && btn.DefaultContour.G == 0 && btn.DefaultContour.B == 0 {
		btn.DefaultContour = color.RGBA{R: 180, G: 180, B: 180, A: 255}
	}
	if btn.DefaultText.A == 0 && btn.DefaultText.R == 0 && btn.DefaultText.G == 0 && btn.DefaultText.B == 0 {
		btn.DefaultText = color.RGBA{R: 51, G: 51, B: 51, A: 255}
	}
	if btn.HoverColor.A == 0 && btn.HoverColor.R == 0 && btn.HoverColor.G == 0 && btn.HoverColor.B == 0 {
		btn.HoverColor = color.RGBA{R: 220, G: 220, B: 220, A: 255}
	}
	if btn.PressedColor.A == 0 && btn.PressedColor.R == 0 && btn.PressedColor.G == 0 && btn.PressedColor.B == 0 {
		btn.PressedColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	}
	if btn.PressedContour.A == 0 && btn.PressedContour.R == 0 && btn.PressedContour.G == 0 && btn.PressedContour.B == 0 {
		btn.PressedContour = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	}
	if btn.PressedText.A == 0 && btn.PressedText.R == 0 && btn.PressedText.G == 0 && btn.PressedText.B == 0 {
		btn.PressedText = color.RGBA{R: 30, G: 30, B: 30, A: 255}
	}
	if btn.ActiveColor.A == 0 && btn.ActiveColor.R == 0 && btn.ActiveColor.G == 0 && btn.ActiveColor.B == 0 {
		btn.ActiveColor = color.RGBA{R: 200, G: 230, B: 255, A: 255}
	}
	if btn.ActiveContour.A == 0 && btn.ActiveContour.R == 0 && btn.ActiveContour.G == 0 && btn.ActiveContour.B == 0 {
		btn.ActiveContour = color.RGBA{R: 100, G: 180, B: 255, A: 255}
	}
	if btn.ActiveText.A == 0 && btn.ActiveText.R == 0 && btn.ActiveText.G == 0 && btn.ActiveText.B == 0 {
		btn.ActiveText = color.RGBA{R: 0, G: 100, B: 200, A: 255}
	}

	btn.ExtendBaseWidget(btn)
	return btn
}

// NewBorderButtonWithAllColors 直接指定所有颜色创建按钮
func NewBorderButtonWithAllColors(onToggle func(active bool), text string,
	defaultColor, defaultContour, defaultText,
	hoverColor,
	pressedColor, pressedContour, pressedText,
	activeColor, activeContour, activeText color.RGBA) *BorderButton {

	style := BorderButtonStyle{
		DefaultColor:       defaultColor,
		DefaultContour:     defaultContour,
		DefaultText:        defaultText,
		HoverColor:         hoverColor,
		PressedColor:       pressedColor,
		PressedContour:     pressedContour,
		PressedText:        pressedText,
		ActiveColor:        activeColor,
		ActiveContour:      activeContour,
		ActiveText:         activeText,
		BorderRadius:       8,
		ContourScale:       0.85,
		ContourLineWidth:   2.0,
		ContourWidthScale:  0.85,
		ContourHeightScale: 0.85,
		PaddingLeft:        12,
		PaddingRight:       12,
		PaddingTop:         6,
		PaddingBottom:      6,
		UseGGFont:          false,
		GGFontType:         "chinese",
		GGFontSize:         14,
		GGFontColor:        defaultText,
	}

	return NewBorderButtonWithStyle(onToggle, text, style)
}

// 可选参数配置函数
func WithBorderButtonStyle(style BorderButtonStyle) func(*BorderButton) {
	return func(btn *BorderButton) {
		if style.DefaultColor.A != 0 || style.DefaultColor.R != 0 || style.DefaultColor.G != 0 || style.DefaultColor.B != 0 {
			btn.DefaultColor = style.DefaultColor
		}
		if style.DefaultContour.A != 0 || style.DefaultContour.R != 0 || style.DefaultContour.G != 0 || style.DefaultContour.B != 0 {
			btn.DefaultContour = style.DefaultContour
		}
		if style.DefaultText.A != 0 || style.DefaultText.R != 0 || style.DefaultText.G != 0 || style.DefaultText.B != 0 {
			btn.DefaultText = style.DefaultText
		}
		if style.HoverColor.A != 0 || style.HoverColor.R != 0 || style.HoverColor.G != 0 || style.HoverColor.B != 0 {
			btn.HoverColor = style.HoverColor
		}
		if style.PressedColor.A != 0 || style.PressedColor.R != 0 || style.PressedColor.G != 0 || style.PressedColor.B != 0 {
			btn.PressedColor = style.PressedColor
		}
		if style.PressedContour.A != 0 || style.PressedContour.R != 0 || style.PressedContour.G != 0 || style.PressedContour.B != 0 {
			btn.PressedContour = style.PressedContour
		}
		if style.PressedText.A != 0 || style.PressedText.R != 0 || style.PressedText.G != 0 || style.PressedText.B != 0 {
			btn.PressedText = style.PressedText
		}
		if style.ActiveColor.A != 0 || style.ActiveColor.R != 0 || style.ActiveColor.G != 0 || style.ActiveColor.B != 0 {
			btn.ActiveColor = style.ActiveColor
		}
		if style.ActiveContour.A != 0 || style.ActiveContour.R != 0 || style.ActiveContour.G != 0 || style.ActiveContour.B != 0 {
			btn.ActiveContour = style.ActiveContour
		}
		if style.ActiveText.A != 0 || style.ActiveText.R != 0 || style.ActiveText.G != 0 || style.ActiveText.B != 0 {
			btn.ActiveText = style.ActiveText
		}
		if style.BorderRadius > 0 {
			btn.BorderRadius = style.BorderRadius
		}
		if style.ContourScale > 0 {
			btn.ContourScale = style.ContourScale
		}
		if style.ContourLineWidth > 0 {
			btn.ContourLineWidth = style.ContourLineWidth
		}
		if style.PaddingLeft > 0 {
			btn.PaddingLeft = style.PaddingLeft
		}
		if style.PaddingRight > 0 {
			btn.PaddingRight = style.PaddingRight
		}
		if style.PaddingTop > 0 {
			btn.PaddingTop = style.PaddingTop
		}
		if style.PaddingBottom > 0 {
			btn.PaddingBottom = style.PaddingBottom
		}
		if style.UseGGFont {
			btn.UseGGFont = style.UseGGFont
			btn.GGFontType = style.GGFontType
			btn.GGFontSize = style.GGFontSize
			btn.GGFontColor = style.GGFontColor
			btn.GGFontOffsetX = style.GGFontOffsetX
			btn.GGFontOffsetY = style.GGFontOffsetY
		}
	}
}

// SetActive 设置按钮激活状态
func (btn *BorderButton) SetActive(active bool) {
	if btn.isActive != active {
		btn.isActive = active
		btn.Refresh()
		if btn.OnToggle != nil {
			btn.OnToggle(active)
		}
	}
}

// Toggle 切换按钮状态
func (btn *BorderButton) Toggle() {
	btn.SetActive(!btn.isActive)
}

// IsActive 获取按钮激活状态
func (btn *BorderButton) IsActive() bool {
	return btn.isActive
}

// SetAllColors 设置所有颜色
func (btn *BorderButton) SetAllColors(
	defaultColor, defaultContour, defaultText,
	hoverColor,
	pressedColor, pressedContour, pressedText,
	activeColor, activeContour, activeText color.RGBA) {

	btn.DefaultColor = defaultColor
	btn.DefaultContour = defaultContour
	btn.DefaultText = defaultText
	btn.HoverColor = hoverColor
	btn.PressedColor = pressedColor
	btn.PressedContour = pressedContour
	btn.PressedText = pressedText
	btn.ActiveColor = activeColor
	btn.ActiveContour = activeContour
	btn.ActiveText = activeText

	btn.Refresh()
}

// SetContourScale 设置轮廓比例
func (btn *BorderButton) SetContourScale(scale float64) {
	if scale > 0 && scale <= 1.0 {
		btn.ContourScale = scale
		btn.Refresh()
	}
}

// SetContourLineWidth 设置轮廓线宽度
func (btn *BorderButton) SetContourLineWidth(width float64) {
	if width > 0 {
		btn.ContourLineWidth = width
		btn.Refresh()
	}
}

// Destroy 销毁按钮
func (btn *BorderButton) Destroy() {
	select {
	case btn.stopAnimation <- true:
	default:
	}
}

// SetSize 设置按钮尺寸
func (btn *BorderButton) SetSize(width, height float32) {
	btn.Width = width
	btn.Height = height
	btn.Refresh()
}

// SetText 设置按钮文字
func (btn *BorderButton) SetText(text string) {
	btn.Text = text
	btn.Refresh()
}

// Hover事件：鼠标移入
func (btn *BorderButton) MouseIn(*desktop.MouseEvent) {
	btn.isHovered = true
	btn.Refresh()
}

// Hover事件：鼠标移出
func (btn *BorderButton) MouseOut() {
	btn.isHovered = false
	btn.Refresh()
}

// 鼠标移动
func (btn *BorderButton) MouseMoved(*desktop.MouseEvent) {}

// Tapped 处理点击事件
func (btn *BorderButton) Tapped(e *fyne.PointEvent) {
	btn.isPressed = true
	btn.Refresh()

	// 切换激活状态
	btn.Toggle()

	// 短暂显示按下效果
	go func() {
		time.Sleep(100 * time.Millisecond)
		fyne.Do(func() {
			btn.isPressed = false
			btn.Refresh()
		})
	}()
}

// CreateRenderer 创建渲染器
func (btn *BorderButton) CreateRenderer() fyne.WidgetRenderer {
	mainCanvas := canvas.NewImageFromImage(nil)
	mainCanvas.FillMode = canvas.ImageFillOriginal

	// 计算最小尺寸（基于内边距）
	minWidth := float32(btn.PaddingLeft + btn.PaddingRight + 20)  // 最小20像素内容宽度
	minHeight := float32(btn.PaddingTop + btn.PaddingBottom + 20) // 最小20像素内容高度

	// 确保不小于指定的宽高
	if btn.Width < minWidth {
		btn.Width = minWidth
	}
	if btn.Height < minHeight {
		btn.Height = minHeight
	}

	container := container.NewWithoutLayout(mainCanvas)
	container.Resize(fyne.NewSize(btn.Width, btn.Height))

	renderer := &borderButtonRenderer{
		btn:        btn,
		mainCanvas: mainCanvas,
		container:  container,
		objects:    []fyne.CanvasObject{container},
	}

	// 立即刷新一次，确保初始显示
	renderer.Refresh()

	return renderer
}

// borderButtonRenderer 渲染器
type borderButtonRenderer struct {
	btn        *BorderButton
	mainCanvas *canvas.Image
	container  *fyne.Container
	objects    []fyne.CanvasObject
}

func (r *borderButtonRenderer) Destroy() {
	r.btn.Destroy()
}

func (r *borderButtonRenderer) Layout(size fyne.Size) {
	r.btn.Width = size.Width
	r.btn.Height = size.Height

	// 更新画布尺寸
	r.mainCanvas.Resize(size)

	// 容器尺寸
	r.container.Resize(size)

	// 布局改变后也需要刷新
	r.Refresh()
}

func (r *borderButtonRenderer) MinSize() fyne.Size {
	// 最小尺寸基于内边距
	minWidth := float32(r.btn.PaddingLeft + r.btn.PaddingRight + 20)
	minHeight := float32(r.btn.PaddingTop + r.btn.PaddingBottom + 20)

	// 确保不小于指定的最小宽高
	if minWidth < r.btn.Width {
		minWidth = r.btn.Width
	}
	if minHeight < r.btn.Height {
		minHeight = r.btn.Height
	}

	return fyne.NewSize(minWidth, minHeight)
}

func (r *borderButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// drawButton 绘制按钮（模仿CSS default样式，开关式）
func (r *borderButtonRenderer) drawButton() *gg.Context {
	width := float64(r.btn.Width)
	height := float64(r.btn.Height)

	dc := gg.NewContext(int(width), int(height))

	// 根据按钮状态选择颜色
	var bgColor, textColor color.RGBA
	var showContour bool
	var contourColor color.RGBA

	if r.btn.isActive {
		// 激活状态：常亮颜色，显示轮廓
		bgColor = r.btn.ActiveColor
		textColor = r.btn.ActiveText
		showContour = true
		contourColor = r.btn.ActiveContour
	} else if r.btn.isPressed {
		// 按下状态（但未激活）- 不显示轮廓
		bgColor = r.btn.PressedColor
		textColor = r.btn.PressedText
		showContour = false
	} else if r.btn.isHovered {
		// 悬停状态 - 不显示轮廓
		bgColor = r.btn.HoverColor
		textColor = r.btn.DefaultText
		showContour = false
	} else {
		// 默认状态 - 不显示轮廓
		bgColor = r.btn.DefaultColor
		textColor = r.btn.DefaultText
		showContour = false
	}

	// 1. 先绘制背景（整个按钮区域）
	dc.SetColor(bgColor)
	dc.DrawRoundedRectangle(0, 0, width, height, r.btn.BorderRadius)
	dc.Fill()

	// 2. 只有在激活状态才绘制独立的外围轮廓
	if showContour {
		// 计算轮廓尺寸（按钮尺寸的比例，宽高分别单独控制）
		contourWidth := width * r.btn.ContourWidthScale
		contourHeight := height * r.btn.ContourHeightScale

		// 计算轮廓位置（精确居中）
		contourX := (width - contourWidth) / 2
		contourY := (height - contourHeight) / 2

		// 计算轮廓圆角（按比例缩小，取宽高比例的平均值）
		contourRadius := r.btn.BorderRadius * ((r.btn.ContourWidthScale + r.btn.ContourHeightScale) / 2)

		// 绘制独立的外围轮廓（没有填充，只有线条）
		dc.SetColor(contourColor)
		dc.SetLineWidth(r.btn.ContourLineWidth)
		dc.DrawRoundedRectangle(
			contourX,
			contourY,
			contourWidth,
			contourHeight,
			contourRadius,
		)
		dc.Stroke() // 只画线，不填充
	}

	// 3. 绘制文字
	centerX := width / 2
	centerY := height / 2

	// 使用指定的文字颜色
	var finalTextColor color.Color
	if r.btn.UseGGFont && r.btn.GGFontColor != nil {
		finalTextColor = r.btn.GGFontColor
	} else {
		finalTextColor = textColor
	}

	// 尝试加载字体
	fontPath := "ttf/chinese.ttf"
	if r.btn.GGFontType == "english" {
		fontPath = "ttf/english.ttf"
	}

	if err := dc.LoadFontFace(fontPath, r.btn.GGFontSize); err == nil {
		dc.SetColor(finalTextColor)
		dc.DrawStringAnchored(r.btn.Text,
			centerX+r.btn.GGFontOffsetX,
			centerY+r.btn.GGFontOffsetY,
			0.5, 0.5)
	} else {
		// 如果字体加载失败，使用默认绘制
		dc.SetColor(textColor)
		dc.DrawStringAnchored(r.btn.Text, centerX, centerY, 0.5, 0.5)
	}

	return dc
}

func (r *borderButtonRenderer) Refresh() {
	// 绘制按钮图像
	dc := r.drawButton()
	r.mainCanvas.Image = dc.Image()
	r.mainCanvas.Refresh()
}

// ==================== 物理隔离容器 ====================
// WrapWithIsolationContainer 返回一个高度+2px的Stack包裹当前控件，实现物理隔离，防止canvas.Image刷新污染
func (btn *BorderButton) WrapWithIsolationContainer() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.Transparent)
	bg.SetMinSize(fyne.NewSize(float32(btn.Width), float32(btn.Height+2)))
	return container.NewStack(
		bg,
		container.NewCenter(btn),
	)
}
