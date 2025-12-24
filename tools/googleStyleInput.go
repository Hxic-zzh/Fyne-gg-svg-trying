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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// ==================== 颜色常量 ====================
var (
	colorGray        = color.RGBA{153, 153, 153, 255} // #999999
	colorHighlight   = color.RGBA{82, 100, 174, 255}  // #5264AE
	colorUnderline   = color.RGBA{82, 100, 174, 255}  // #5264AE
	colorPulse       = color.RGBA{82, 100, 174, 255}  // #5264AE
	colorBorder      = color.RGBA{117, 117, 117, 255} // #757575
	colorText        = color.Black
	colorTransparent = color.RGBA{0, 0, 0, 0}
)

// ==================== 布局参数结构体 ====================
type LayoutParams struct {
	// 输入框参数（核心）
	EntryX      float64 // 输入框X位置
	EntryY      float64 // 输入框Y位置
	EntryWidth  float64 // 输入框宽度
	EntryHeight float64 // 输入框高度

	// 文字参数（基于输入框计算）
	LabelFontSize float64 // 标签字体大小（基于输入框高度）
	TextFontSize  float64 // 输入文字大小（基于输入框高度）
	LabelX        float64 // 标签X位置（与输入框对齐）
	LabelStartY   float64 // 标签初始Y位置（在输入框内）
	LabelFloatY   float64 // 标签上浮Y位置

	// 下划线参数（基于输入框计算）
	UnderlineY      float64 // 下划线Y位置（输入框底部）
	UnderlineHeight float64 // 下划线粗细
}

// ==================== 自定义Entry类型 ====================
type materialEntryWrapper struct {
	widget.Entry
	parent *MaterialEntry
}

// 重写焦点方法
func (e *materialEntryWrapper) FocusGained() {
	e.Entry.FocusGained()
	if e.parent != nil && !e.parent.focused {
		e.parent.focused = true
		e.parent.pulseActive = true
		e.parent.pulseProgress = 0
		e.parent.requestAnimationUpdate()
	}
}

func (e *materialEntryWrapper) FocusLost() {
	e.Entry.FocusLost()
	if e.parent != nil && e.parent.focused {
		e.parent.focused = false
		e.parent.requestAnimationUpdate()
	}
}

// ==================== MaterialEntry 控件 ====================
type MaterialEntry struct {
	widget.BaseWidget

	// 用户自定义输入框尺寸（以输入框为主）
	EntryWidth  float64 // 输入框宽度
	EntryHeight float64 // 输入框高度

	// CSS-like 属性
	Style MaterialEntryStyle // 新增：css-like风格属性

	// 核心组件
	container *fyne.Container       // 容器：canvas + entry
	canvasImg *canvas.Image         // 底层：绘制Material样式
	entry     *materialEntryWrapper // 上层：自定义Entry

	// 布局参数
	layoutParams LayoutParams

	// 样式状态（基于布尔值判断）
	placeholder string
	focused     bool // 自己维护焦点状态
	hasText     bool
	lastFocused bool // 用于检测焦点变化

	// 动画参数
	labelY         float64 // 标签Y位置
	labelScale     float64 // 标签缩放
	labelColor     color.Color
	underlineWidth float64 // 下划线宽度 (0-1)
	pulseProgress  float64 // 脉冲进度 (0-1)
	pulseActive    bool    // 是否激活脉冲

	// 尺寸
	width, height float64

	// 字体
	fontPath string

	// 回调
	OnChanged   func(string)
	OnSubmitted func(string)

	// 动画控制
	animTicker *time.Ticker
	animStop   chan bool

	// ============ 新增：主题透明控制 ============
	UseThemeTransparency bool        // 是否使用主题透明度
	CustomBackground     color.Color // 自定义背景色（覆盖主题）
	CustomBorder         color.Color // 自定义边框色

	// ============ 新增：透明效果配置 ============
	GlassEffect  bool    // 毛玻璃效果
	BlurRadius   float64 // 模糊半径
	CornerRadius float64 // 圆角半径

	// ============ 新增：主题颜色缓存 ============
	themeBackground color.Color
	themeBorder     color.Color
	themeText       color.Color

	// ============ 字体缓存优化 ============
	cachedFontPath      string
	cachedFontSize      float64
	cachedFontFace      font.Face
	cachedLabelFontSize float64
	cachedLabelFontFace font.Face
}

// ==================== CSS-like 风格结构体 ====================
type MaterialEntryStyle struct {
	Width           float64     // 输入框宽度
	Height          float64     // 输入框高度
	FontSize        float64     // 字体大小
	LabelColor      color.Color // 标签颜色
	TextColor       color.Color // 输入文字颜色
	BorderColor     color.Color // 边框颜色
	BgColor         color.Color // 背景色
	Radius          float64     // 圆角
	UnderlineColor  color.Color // 下划线颜色
	UnderlineHeight float64     // 下划线粗细
	// 可扩展更多css-like属性
}

// ==================== 计算布局参数 ====================
func (m *MaterialEntry) calculateLayout() {
	// 优先用Style里的宽高，否则用EntryWidth/EntryHeight，否则用默认
	var ew, eh float64
	if m.Style.Width > 0 {
		ew = m.Style.Width
	} else if m.EntryWidth > 0 {
		ew = m.EntryWidth
	} else {
		ew = 300
	}
	if m.Style.Height > 0 {
		eh = m.Style.Height
	} else if m.EntryHeight > 0 {
		eh = m.EntryHeight
	} else {
		eh = 40
	}
	paddingY := eh * 0.5
	paddingX := 5.0

	m.layoutParams.EntryWidth = ew
	m.layoutParams.EntryHeight = eh
	m.layoutParams.EntryX = paddingX
	m.layoutParams.EntryY = paddingY

	// 字体大小优先用Style
	if m.Style.FontSize > 0 {
		m.layoutParams.TextFontSize = m.Style.FontSize
	} else {
		m.layoutParams.TextFontSize = eh * 0.45
	}
	m.layoutParams.LabelFontSize = m.layoutParams.TextFontSize * 0.78

	m.layoutParams.LabelX = m.layoutParams.EntryX + 5
	m.layoutParams.LabelStartY = m.layoutParams.EntryY + eh*0.75
	m.layoutParams.LabelFloatY = m.layoutParams.EntryY - m.layoutParams.LabelFontSize*0.7

	m.layoutParams.UnderlineY = m.layoutParams.EntryY + eh + 1
	if m.Style.UnderlineHeight > 0 {
		m.layoutParams.UnderlineHeight = m.Style.UnderlineHeight
	} else {
		m.layoutParams.UnderlineHeight = 4
	}

	m.width = ew + paddingX*2
	m.height = eh + paddingY + m.layoutParams.LabelFontSize + 8
}

// ==================== 设置CSS-like风格方法 ====================
func (m *MaterialEntry) SetStyle(style MaterialEntryStyle) {
	m.Style = style
	m.EntryWidth = style.Width
	m.EntryHeight = style.Height
	// 禁止自定义圆角和背景色
	// m.CornerRadius = style.Radius // 不允许自定义圆角
	// if style.BgColor != nil {
	// 	m.CustomBackground = style.BgColor
	// 	m.UseThemeTransparency = false
	// }
	if style.BorderColor != nil {
		m.CustomBorder = style.BorderColor
	}
	// 其它属性可扩展
	m.calculateLayout()
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

// ==================== 创建方法 ====================
func NewMaterialEntry(placeholder string, entrySize ...float64) *MaterialEntry {
	var ew, eh float64
	if len(entrySize) > 0 {
		ew = entrySize[0]
	}
	if len(entrySize) > 1 {
		eh = entrySize[1]
	}
	m := &MaterialEntry{
		placeholder:          placeholder,
		labelY:               0,
		labelScale:           1.0,
		labelColor:           colorGray,
		underlineWidth:       0,
		pulseProgress:        0,
		pulseActive:          false,
		EntryWidth:           ew,
		EntryHeight:          eh,
		fontPath:             "ttf/english.ttf",
		animStop:             make(chan bool, 1),
		UseThemeTransparency: true,
		GlassEffect:          false,
		BlurRadius:           5.0,
		CornerRadius:         8.0,
		// 初始主题颜色
		themeBackground: color.Transparent,
		themeBorder:     color.RGBA{R: 200, G: 200, B: 200, A: 100},
		themeText:       color.RGBA{R: 33, G: 33, B: 33, A: 255},
	}

	// 1. 计算布局参数（预运算）
	m.calculateLayout()
	m.labelY = m.layoutParams.LabelStartY

	// 2. 创建底层画布
	m.canvasImg = canvas.NewImageFromResource(nil)
	m.canvasImg.FillMode = canvas.ImageFillContain

	// 3. 创建自定义的Entry
	m.entry = &materialEntryWrapper{parent: m}
	m.entry.ExtendBaseWidget(m.entry)
	m.entry.SetPlaceHolder("") // 不显示占位符，我们自己绘制

	// 设置回调
	m.entry.OnChanged = func(text string) {
		m.hasText = len(text) > 0
		m.requestAnimationUpdate()

		if m.OnChanged != nil {
			m.OnChanged(text)
		}
	}

	m.entry.OnSubmitted = func(text string) {
		if m.OnSubmitted != nil {
			m.OnSubmitted(text)
		}
	}

	// 4. 创建容器
	m.container = container.NewWithoutLayout(m.canvasImg, m.entry)

	// 5. 初始化
	m.ExtendBaseWidget(m)

	// 获取当前主题颜色
	m.updateThemeColors()

	// 在主线程中初始更新画布
	fyne.Do(func() {
		m.updateCanvas()
	})

	// 6. 启动动画
	m.startAnimation()

	return m
}

// ==================== 透明效果方法 ====================
// EnableThemeTransparency 启用/禁用主题透明
func (m *MaterialEntry) EnableThemeTransparency(enable bool) {
	m.UseThemeTransparency = enable
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

// SetCustomBackground 设置自定义背景（覆盖主题）
func (m *MaterialEntry) SetCustomBackground(color color.Color) {
	m.CustomBackground = color
	m.UseThemeTransparency = false
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

// SetCustomBorder 设置自定义边框
func (m *MaterialEntry) SetCustomBorder(color color.Color) {
	m.CustomBorder = color
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

// SetGlassEffect 设置毛玻璃效果
func (m *MaterialEntry) SetGlassEffect(enable bool, blurRadius float64) {
	m.GlassEffect = enable
	m.BlurRadius = blurRadius
	if enable {
		// 毛玻璃效果：半透明白色背景
		m.CustomBackground = color.RGBA{R: 255, G: 255, B: 255, A: 180}
		m.CustomBorder = color.RGBA{R: 255, G: 255, B: 255, A: 100}
		m.UseThemeTransparency = false
	}
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

// SetTransparent 设置为完全透明
func (m *MaterialEntry) SetTransparent(transparent bool) {
	if transparent {
		m.CustomBackground = color.Transparent
		m.CustomBorder = color.Transparent
		m.UseThemeTransparency = false
	} else {
		m.UseThemeTransparency = true
	}
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

// SetCornerRadius 设置圆角半径
func (m *MaterialEntry) SetCornerRadius(radius float64) {
	m.CornerRadius = radius
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

// updateThemeColors 更新主题颜色
func (m *MaterialEntry) updateThemeColors() {
	// 获取当前主题的颜色
	currentTheme := fyne.CurrentApp().Settings().Theme()

	// 获取输入框相关颜色
	m.themeBackground = currentTheme.Color(theme.ColorNameInputBackground, theme.VariantLight)
	m.themeBorder = currentTheme.Color(theme.ColorNameInputBorder, theme.VariantLight)
	m.themeText = currentTheme.Color(theme.ColorNameForeground, theme.VariantLight)
}

// ==================== 颜色 getter ====================
func (m *MaterialEntry) getLabelColor(focusedOrFloat bool) color.Color {
	if m.Style.LabelColor != nil {
		return m.Style.LabelColor
	}
	if focusedOrFloat {
		return colorHighlight
	}
	return colorGray
}

func (m *MaterialEntry) getUnderlineColor() color.Color {
	if m.Style.UnderlineColor != nil {
		return m.Style.UnderlineColor
	}
	return colorUnderline
}

func (m *MaterialEntry) getBorderColor() color.Color {
	if m.Style.BorderColor != nil {
		return m.Style.BorderColor
	}
	if !m.UseThemeTransparency && m.CustomBorder != nil {
		return m.CustomBorder
	}
	return m.themeBorder
}

func (m *MaterialEntry) getTextColor() color.Color {
	if m.Style.TextColor != nil {
		return m.Style.TextColor
	}
	return colorText
}

func (m *MaterialEntry) getPulseColor() color.Color {
	if m.Style.UnderlineColor != nil {
		return m.Style.UnderlineColor // 脉冲色通常与高亮下划线色一致
	}
	return colorPulse
}

// ==================== 动画控制 ====================
func (m *MaterialEntry) startAnimation() {
	// 使用较慢的动画刷新率，避免过度刷新
	m.animTicker = time.NewTicker(33 * time.Millisecond) // ~30fps

	go func() {
		for {
			select {
			case <-m.animTicker.C:
				// 只在动画协程中更新数据状态
				m.updateAnimation()
			case <-m.animStop:
				if m.animTicker != nil {
					m.animTicker.Stop()
					m.animTicker = nil
				}
				return
			}
		}
	}()
}

func (m *MaterialEntry) stopAnimation() {
	if m.animStop != nil {
		select {
		case m.animStop <- true:
		default:
		}
	}
}

func (m *MaterialEntry) requestAnimationUpdate() {
	// 标记需要更新，实际更新会在动画循环中进行
}

// ==================== 动画更新 ====================
func (m *MaterialEntry) updateAnimation() {
	needsUpdate := false

	// 更新脉冲动画
	if m.pulseActive {
		m.pulseProgress += 0.03 // 减慢脉冲速度
		if m.pulseProgress >= 1.0 {
			m.pulseActive = false
			m.pulseProgress = 0
		}
		needsUpdate = true
	}

	// 更新标签和下划线动画
	shouldFloat := m.hasText || m.focused
	targetLabelY := m.layoutParams.LabelStartY
	targetLabelScale := 1.0
	targetLabelColor := m.getLabelColor(false)
	targetUnderlineWidth := 0.0

	if shouldFloat {
		targetLabelY = m.layoutParams.LabelFloatY
		targetLabelScale = 0.78 // CSS: font-size从18px变为14px (14/18=0.78)
		targetLabelColor = m.getLabelColor(true)
	}

	if m.focused {
		targetUnderlineWidth = 1.0
	}

	// 平滑过渡
	oldLabelY := m.labelY
	oldLabelScale := m.labelScale
	oldUnderlineWidth := m.underlineWidth
	oldLabelColor := m.labelColor

	m.labelY = lerp(m.labelY, targetLabelY, 0.15) // 减慢动画速度
	m.labelScale = lerp(m.labelScale, targetLabelScale, 0.15)
	m.labelColor = interpolateColor(m.labelColor, targetLabelColor, 0.15)
	m.underlineWidth = lerp(m.underlineWidth, targetUnderlineWidth, 0.2)

	// ====== 优化：下划线动画不做阈值判断，保证动画完整 ======
	const floatThreshold = 0.1
	colorChanged := func(a, b color.Color) bool {
		r1, g1, b1, a1 := a.RGBA()
		r2, g2, b2, a2 := b.RGBA()
		return absInt(int(r1)-int(r2)) > 8 || absInt(int(g1)-int(g2)) > 8 || absInt(int(b1)-int(b2)) > 8 || absInt(int(a1)-int(a2)) > 8
	}

	if m.pulseActive ||
		absFloat(oldLabelY-m.labelY) > floatThreshold ||
		absFloat(oldLabelScale-m.labelScale) > floatThreshold ||
		colorChanged(oldLabelColor, m.labelColor) ||
		oldUnderlineWidth != m.underlineWidth { // underlineWidth 只要有变化就刷新
		needsUpdate = true
	}

	// 在主线程中执行UI刷新
	if needsUpdate {
		fyne.Do(func() {
			m.updateCanvas()
			m.canvasImg.Refresh()
		})
	}
}

func absFloat(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func absInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// ==================== 画布绘制 ====================
func (m *MaterialEntry) updateCanvas() {
	// 确保最小尺寸
	if m.width < 200 {
		m.width = 300
	}
	if m.height < 60 {
		m.height = 80
	}

	// 重新计算布局（如果尺寸变化）
	m.calculateLayout()

	// 创建画布上下文
	dc := gg.NewContext(int(m.width), int(m.height))
	if dc == nil {
		return
	}

	// 清空背景（完全透明）
	dc.SetColor(color.Transparent)
	dc.Clear()

	// ============ 绘制输入框背景 ============
	var bgColor color.Color
	if !m.UseThemeTransparency && m.CustomBackground != nil {
		bgColor = m.CustomBackground
	} else {
		bgColor = m.themeBackground
	}

	// 只绘制非完全透明的背景
	if _, _, _, a := bgColor.RGBA(); a > 0 {
		dc.SetColor(bgColor)
		dc.DrawRoundedRectangle(
			m.layoutParams.EntryX,
			m.layoutParams.EntryY,
			m.layoutParams.EntryWidth,
			m.layoutParams.EntryHeight,
			m.CornerRadius,
		)
		dc.Fill()
	}

	// ============ 字体缓存优化 ============
	// 输入字体缓存
	if m.fontPath != "" && m.layoutParams.TextFontSize > 0 {
		if m.cachedFontFace == nil || m.cachedFontPath != m.fontPath || m.cachedFontSize != m.layoutParams.TextFontSize {
			face, err := gg.LoadFontFace(m.fontPath, m.layoutParams.TextFontSize)
			if err == nil {
				m.cachedFontFace = face
				m.cachedFontPath = m.fontPath
				m.cachedFontSize = m.layoutParams.TextFontSize
			} else {
				m.cachedFontFace = nil
			}
		}
		dc.SetFontFace(m.cachedFontFace)
	} else {
		dc.SetFontFace(nil)
	}

	// 绘制底部静态下划线
	dc.SetColor(m.getBorderColor())
	dc.SetLineWidth(1)
	dc.DrawLine(
		m.layoutParams.EntryX,
		m.layoutParams.UnderlineY,
		m.layoutParams.EntryX+m.layoutParams.EntryWidth,
		m.layoutParams.UnderlineY,
	)
	dc.Stroke()

	// 绘制高亮下划线（动态线）
	if m.underlineWidth > 0 {
		underlineWidth := m.layoutParams.EntryWidth * m.underlineWidth * 0.5 // 每边50%
		centerX := m.layoutParams.EntryX + (m.layoutParams.EntryWidth / 2)

		dc.SetColor(m.getUnderlineColor())
		dc.SetLineWidth(m.layoutParams.UnderlineHeight)

		// 左边部分（从中间向左展开）
		dc.DrawLine(centerX-underlineWidth, m.layoutParams.UnderlineY, centerX, m.layoutParams.UnderlineY)
		// 右边部分（从中间向右展开）
		dc.DrawLine(centerX, m.layoutParams.UnderlineY, centerX+underlineWidth, m.layoutParams.UnderlineY)
		dc.Stroke()
	}

	// 绘制脉冲效果
	if m.pulseProgress > 0 {
		pulseAlpha := (1 - m.pulseProgress) * 0.3
		pulseWidth := 100 * (1 - m.pulseProgress)

		dc.SetColor(colorWithAlpha(m.getPulseColor(), uint8(pulseAlpha*255)))
		dc.SetLineWidth(3)

		pulseStartX := m.layoutParams.EntryX
		dc.DrawLine(
			pulseStartX,
			m.layoutParams.UnderlineY,
			pulseStartX+pulseWidth,
			m.layoutParams.UnderlineY,
		)
		dc.Stroke()
	}

	// 绘制浮动标签
	dc.SetColor(m.labelColor)

	// 标签字体缓存
	if m.fontPath != "" && m.layoutParams.LabelFontSize*m.labelScale > 0 {
		labelFontSize := m.layoutParams.LabelFontSize * m.labelScale
		if m.cachedLabelFontFace == nil || m.cachedFontPath != m.fontPath || m.cachedLabelFontSize != labelFontSize {
			face, err := gg.LoadFontFace(m.fontPath, labelFontSize)
			if err == nil {
				m.cachedLabelFontFace = face
				m.cachedFontPath = m.fontPath
				m.cachedLabelFontSize = labelFontSize
			} else {
				m.cachedLabelFontFace = nil
			}
		}
		dc.SetFontFace(m.cachedLabelFontFace)
	} else {
		dc.SetFontFace(nil)
	}

	// 标签位置（已经过平滑过渡）
	dc.DrawString(m.placeholder, m.layoutParams.LabelX, m.labelY)

	// （已移除：canvas不再绘制输入内容，实际内容只由entry控件渲染）

	// 更新画布图像
	m.canvasImg.Image = dc.Image()
}

// ==================== 渲染器实现 ====================
func (m *MaterialEntry) CreateRenderer() fyne.WidgetRenderer {
	return &materialEntryRenderer{m}
}

type materialEntryRenderer struct {
	entry *MaterialEntry
}

func (r *materialEntryRenderer) Destroy() {
	r.entry.stopAnimation()
}

func (r *materialEntryRenderer) Layout(size fyne.Size) {
	// 以输入框为主，画布和Entry自适应
	r.entry.calculateLayout()
	r.entry.canvasImg.Resize(fyne.NewSize(float32(r.entry.width), float32(r.entry.height)))
	r.entry.entry.Resize(fyne.NewSize(
		float32(r.entry.layoutParams.EntryWidth),
		float32(r.entry.layoutParams.EntryHeight),
	))
	r.entry.entry.Move(fyne.NewPos(
		float32(r.entry.layoutParams.EntryX),
		float32(r.entry.layoutParams.EntryY),
	))
	// 在主线程中更新画布
	fyne.Do(func() {
		r.entry.updateCanvas()
	})
}

func (r *materialEntryRenderer) MinSize() fyne.Size {
	// 返回整体控件最小尺寸
	return fyne.NewSize(float32(r.entry.width), float32(r.entry.height))
}

func (r *materialEntryRenderer) Objects() []fyne.CanvasObject {
	return r.entry.container.Objects
}

func (r *materialEntryRenderer) Refresh() {
	// 更新主题颜色
	r.entry.updateThemeColors()

	// 在主线程中更新画布
	fyne.Do(func() {
		r.entry.updateCanvas()
		r.entry.canvasImg.Refresh()
	})
}

// ==================== 工具函数 ====================
func lerp(current, target, factor float64) float64 {
	return current + (target-current)*factor
}

func interpolateColor(c1, c2 color.Color, factor float64) color.Color {
	if factor <= 0 {
		return c1
	}
	if factor >= 1 {
		return c2
	}

	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return color.RGBA{
		uint8(float64(uint8(r1>>8)) + (float64(uint8(r2>>8))-float64(uint8(r1>>8)))*factor),
		uint8(float64(uint8(g1>>8)) + (float64(uint8(g2>>8))-float64(uint8(g1>>8)))*factor),
		uint8(float64(uint8(b1>>8)) + (float64(uint8(b2>>8))-float64(uint8(b1>>8)))*factor),
		uint8(float64(uint8(a1>>8)) + (float64(uint8(a2>>8))-float64(uint8(a1>>8)))*factor),
	}
}

func colorWithAlpha(c color.Color, alpha uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), alpha}
}

// ==================== 公开方法 ====================
func (m *MaterialEntry) SetText(text string) {
	m.entry.SetText(text)
	m.hasText = len(text) > 0

	// 在主线程中更新UI
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

func (m *MaterialEntry) GetText() string {
	return m.entry.Text
}

func (m *MaterialEntry) Clear() {
	m.entry.SetText("")
	m.hasText = false

	// 在主线程中更新UI
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

func (m *MaterialEntry) SetPlaceholder(placeholder string) {
	m.placeholder = placeholder

	// 在主线程中更新UI
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

func (m *MaterialEntry) SetFontPath(path string) {
	m.fontPath = path

	// 在主线程中更新UI
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

func (m *MaterialEntry) FocusGained() {
	if m.entry != nil {
		fyne.CurrentApp().Driver().CanvasForObject(m).Focus(m.entry)
	}
}

func (m *MaterialEntry) FocusLost() {
	// 焦点丢失由自定义entry处理
}

func (m *MaterialEntry) Focused() bool {
	return m.focused
}

func (m *MaterialEntry) TypedRune(r rune) {
	if m.entry != nil {
		m.entry.TypedRune(r)
	}
}

func (m *MaterialEntry) TypedKey(ev *fyne.KeyEvent) {
	if m.entry != nil {
		m.entry.TypedKey(ev)
	}
}

func (m *MaterialEntry) MinSize() fyne.Size {
	return fyne.NewSize(300, 80)
}

func (m *MaterialEntry) Resize(size fyne.Size) {
	m.BaseWidget.Resize(size)
	m.container.Resize(size)

	// 在主线程中更新UI
	fyne.Do(func() {
		m.updateCanvas()
		m.canvasImg.Refresh()
	})
}

// ==================== 鼠标交互 ====================
func (m *MaterialEntry) Tapped(_ *fyne.PointEvent) {
	if m.entry != nil {
		fyne.CurrentApp().Driver().CanvasForObject(m).Focus(m.entry)
	}
}

func (m *MaterialEntry) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// ==================== 物理隔离容器 ====================
// WrapWithIsolationContainer 返回一个高度+2px的Stack包裹当前控件，实现物理隔离，防止canvas.Image刷新污染
func (m *MaterialEntry) WrapWithIsolationContainer() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.Transparent)
	bg.SetMinSize(fyne.NewSize(float32(m.width), float32(m.height+2)))
	return container.NewStack(
		bg,
		container.NewCenter(m),
	)
}

func (m *MaterialEntry) GetEntry() *widget.Entry {
	if m.entry != nil {
		return &m.entry.Entry
	}
	return nil
}
