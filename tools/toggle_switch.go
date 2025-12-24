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
	"math"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
)

// ==================== 颜色定义 ====================
var (
	ColorOn       = color.RGBA{R: 3, G: 169, B: 244, A: 255}   // #03a9f4
	ColorOff      = color.RGBA{R: 244, G: 67, B: 54, A: 255}   // #f44336
	ColorBgOn     = color.RGBA{R: 235, G: 247, B: 252, A: 255} // #ebf7fc
	ColorBgOff    = color.RGBA{R: 252, G: 235, B: 235, A: 255} // #fcebeb
	ColorText     = color.RGBA{R: 255, G: 255, B: 255, A: 255} // 白色
	ColorTextDark = color.RGBA{R: 78, G: 78, B: 78, A: 255}    // #4e4e4e
)

// ==================== 类型定义 ====================
type SwitchEffect int

const (
	EffectSlide SwitchEffect = iota
	EffectTwoBallSwap
	EffectVerticalSwap
	EffectProjectionFlip
	EffectRectSlide
	EffectColorBlockSwap
	// ... 其他效果
)

// ==================== 自定义配置结构体 ====================
type SwitchConfig struct {
	YesLabel      string
	NoLabel       string
	FontPath      string
	YesColor      color.RGBA
	NoColor       color.RGBA
	YesBgColor    color.RGBA
	NoBgColor     color.RGBA
	TextColor     color.RGBA
	TextDarkColor color.RGBA
	YesValue      bool // 新增：左侧布尔值
	NoValue       bool // 新增：右侧布尔值
}

// ==================== 默认配置 ====================
var DefaultConfig = SwitchConfig{
	YesLabel:      "YES",
	NoLabel:       "NO",
	FontPath:      "ttf/toggle_switch.ttf",
	YesColor:      ColorOn,
	NoColor:       ColorOff,
	YesBgColor:    ColorBgOn,
	NoBgColor:     ColorBgOff,
	TextColor:     ColorText,
	TextDarkColor: ColorTextDark,
	YesValue:      true,  // 默认左true
	NoValue:       false, // 默认右false
}

// ==================== 开关控件 ====================
type drawParams struct {
	width        float64
	height       float64
	padding      float64
	knobHeight   float64
	knobPadding  float64
	knobDiameter float64
	knobRadius   float64
	fontSize     float64
	leftLimit    float64
	rightLimit   float64
	centerY      float64
}

type ToggleSwitch struct {
	widget.BaseWidget

	// 状态
	Checked   bool
	OnChanged func(value bool) // value为自定义布尔值
	Disabled  bool

	// 精确尺寸（匹配CSS）
	Width      float64
	Height     float64
	Effect     SwitchEffect
	Animation  bool
	ShowLabels bool

	// 自定义配置
	Config SwitchConfig

	// 交互状态
	isHovered bool
	isPressed bool

	// 动画值
	animValue float64 // 0.0-1.0

	perspective      float64 // 透视距离（像素）
	animationRunning int32   // 0:未运行 1:运行中

	drawCache drawParams // 新增：缓存绘制参数
}

// 计算并缓存绘制参数
func (s *ToggleSwitch) calcDrawParams() {
	w := s.Width
	h := s.Height
	padding := w / 18.5
	knobHeight := w / 7.4
	knobPadding := w / 8.22
	knobTotalHeight := knobHeight + knobPadding*2
	knobDiameter := knobTotalHeight
	knobRadius := knobDiameter / 2
	fontSize := w / 7.4
	leftLimit := w / 18.5
	rightLimit := w - knobDiameter - leftLimit
	centerY := h / 2

	s.drawCache = drawParams{
		width:        w,
		height:       h,
		padding:      padding,
		knobHeight:   knobHeight,
		knobPadding:  knobPadding,
		knobDiameter: knobDiameter,
		knobRadius:   knobRadius,
		fontSize:     fontSize,
		leftLimit:    leftLimit,
		rightLimit:   rightLimit,
		centerY:      centerY,
	}
}

// ==================== 构造函数 ====================
func NewToggleSwitch(checked bool) *ToggleSwitch {
	s := &ToggleSwitch{
		Checked:     checked,
		Effect:      EffectSlide, // 第一个效果
		Width:       74,          // 精确匹配CSS
		Height:      36,          // 精确匹配CSS
		Animation:   true,
		ShowLabels:  true,
		Config:      DefaultConfig,
		animValue:   0.0,
		perspective: 60.0, // 和CSS的60px对应
	}

	if checked {
		s.animValue = 1.0
	}

	s.calcDrawParams() // 新增：初始化时计算参数

	s.ExtendBaseWidget(s)

	return s
}

// ==================== 颜色混合函数 ====================
func mixColor(c1, c2 color.RGBA, t float64) color.RGBA {
	clamp := func(v float64) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v)
	}
	return color.RGBA{
		R: clamp(float64(c1.R)*(1-t) + float64(c2.R)*t),
		G: clamp(float64(c1.G)*(1-t) + float64(c2.G)*t),
		B: clamp(float64(c1.B)*(1-t) + float64(c2.B)*t),
		A: clamp(float64(c1.A)*(1-t) + float64(c2.A)*t),
	}
}

// ==================== 动画系统 ====================
func (s *ToggleSwitch) startAnimation() {
	if !atomic.CompareAndSwapInt32(&s.animationRunning, 0, 1) {
		return // 已有动画在跑
	}
	go s.animationLoop()
}

func (s *ToggleSwitch) animationLoop() {
	defer atomic.StoreInt32(&s.animationRunning, 0)
	for {
		time.Sleep(16 * time.Millisecond)
		target := 0.0
		if s.Checked {
			target = 1.0
		}
		if math.Abs(s.animValue-target) < 0.01 {
			s.animValue = target
			fyne.Do(func() { s.Refresh() })
			break // 动画完成，退出循环
		}
		// 动画推进逻辑
		needsRefresh := false
		if s.Effect == EffectTwoBallSwap || s.Effect == EffectVerticalSwap {
			diff := target - s.animValue
			s.animValue += diff * 0.15
			needsRefresh = true
		} else if s.Effect == EffectProjectionFlip {
			diff := target - s.animValue
			s.animValue += diff * 0.25
			needsRefresh = true
		} else if s.Effect == EffectRectSlide || s.Effect == EffectColorBlockSwap {
			diff := target - s.animValue
			s.animValue += diff * 0.4
			needsRefresh = true
		} else {
			diff := target - s.animValue
			s.animValue += diff * 0.4
			needsRefresh = true
		}
		if needsRefresh {
			fyne.Do(func() { s.Refresh() })
		}
	}
}

// ==================== 交互接口 ====================
func (s *ToggleSwitch) Tapped(e *fyne.PointEvent) {
	if s.Disabled {
		return
	}

	s.Checked = !s.Checked
	s.startAnimation() // 切换时启动动画

	if s.OnChanged != nil {
		s.OnChanged(s.Value()) // 传递自定义布尔值
	}
}

func (s *ToggleSwitch) MouseIn(*desktop.MouseEvent) {
	s.isHovered = true
	s.Refresh()
}

func (s *ToggleSwitch) MouseOut() {
	s.isHovered = false
	s.isPressed = false
	s.Refresh()
}

func (s *ToggleSwitch) MouseDown(*desktop.MouseEvent) {
	s.isPressed = true
	s.Refresh()
}

func (s *ToggleSwitch) MouseUp(*desktop.MouseEvent) {
	s.isPressed = false
	s.Refresh()
}

// ==================== 渲染器 ====================
type toggleSwitchRenderer struct {
	s       *ToggleSwitch
	image   *canvas.Image
	objects []fyne.CanvasObject
}

func (s *ToggleSwitch) CreateRenderer() fyne.WidgetRenderer {
	img := canvas.NewImageFromImage(nil)
	img.FillMode = canvas.ImageFillOriginal

	return &toggleSwitchRenderer{
		s:       s,
		image:   img,
		objects: []fyne.CanvasObject{img},
	}
}

func (r *toggleSwitchRenderer) Destroy() {}

func (r *toggleSwitchRenderer) Layout(size fyne.Size) {
	r.image.Resize(size)
	r.Refresh()
}

func (r *toggleSwitchRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(r.s.Width), float32(r.s.Height))
}

func (r *toggleSwitchRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// ==================== 核心绘制方法 ====================
func (r *toggleSwitchRenderer) Refresh() {
	dc := r.drawSwitch()
	r.image.Image = dc.Image()
	r.image.Refresh()
}

func (r *toggleSwitchRenderer) drawSwitch() *gg.Context {
	width := int(r.s.Width)
	height := int(r.s.Height)

	dc := gg.NewContext(width, height)
	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	switch r.s.Effect {
	case EffectSlide:
		r.drawEffect1Exact(dc)
	case EffectTwoBallSwap:
		r.drawEffect2TwoBallSwap(dc)
	case EffectVerticalSwap:
		r.drawEffect3VerticalSwap(dc)
	case EffectProjectionFlip:
		r.drawEffect5ProjectionFlip(dc)
	case EffectRectSlide:
		r.drawEffect10RectSlide(dc)
	case EffectColorBlockSwap:
		r.drawEffect12ColorBlockSwap(dc)
	default:
		r.drawEffect1Exact(dc)
	}

	return dc
}

// ==================== 精确实现第一个效果 ====================
func (r *toggleSwitchRenderer) drawEffect1Exact(dc *gg.Context) {
	p := r.s.drawCache
	width := p.width
	height := p.height

	// 1. 绘制背景
	bgColor := r.s.Config.YesBgColor
	if r.s.Checked {
		bgColor = r.s.Config.NoBgColor
	}

	bgRadius := math.Min(height/2, 25.0)
	dc.SetColor(bgColor)
	dc.DrawRoundedRectangle(0, 0, width, height, bgRadius)
	dc.Fill()

	// 2. 滑块尺寸
	knobDiameter := p.knobDiameter
	knobRadius := p.knobRadius

	// 3. 滑块可移动范围
	leftLimit := p.leftLimit
	rightLimit := p.rightLimit

	// 4. 根据animValue计算位置
	currentX := leftLimit + (rightLimit-leftLimit)*r.s.animValue

	// 垂直居中
	knobY := (height - knobDiameter) / 2

	// 5. 滑块颜色
	knobColor := r.s.Config.YesColor
	if r.s.Checked {
		knobColor = r.s.Config.NoColor
	}

	// 绘制滑块
	dc.SetColor(knobColor)
	centerX := currentX + knobRadius
	centerY := knobY + knobRadius
	dc.DrawCircle(centerX, centerY, knobRadius)
	dc.Fill()

	// 6. 绘制文字
	dc.SetColor(r.s.Config.TextColor)

	// 加载字体并自适应大小
	fontSize := p.fontSize
	err := dc.LoadFontFace(r.s.Config.FontPath, fontSize)
	if err != nil {
		dc.LoadFontFace("", fontSize)
	}

	// 文字内容
	text := r.s.Config.YesLabel
	if r.s.Checked {
		text = r.s.Config.NoLabel
	}

	// 计算文字位置
	textX := currentX + width/18.0
	if r.s.Checked {
		textX = currentX + width/12.0 // NO
	}

	_, textHeight := dc.MeasureString(text)
	textY := knobY + p.knobPadding + textHeight*0.98

	dc.DrawString(text, textX, textY)
}

// ==================== 修正的双球交换效果（正确字体） ====================
func (r *toggleSwitchRenderer) drawEffect2TwoBallSwap(dc *gg.Context) {
	p := r.s.drawCache
	width := p.width
	height := p.height

	// 1. 绘制背景
	bgColor := r.s.Config.YesBgColor
	if r.s.Checked {
		bgColor = r.s.Config.NoBgColor
	}
	dc.SetColor(bgColor)
	bgRadius := math.Min(height/2, 25.0)
	dc.DrawRoundedRectangle(0, 0, width, height, bgRadius)
	dc.Fill()

	// 2. 计算关键尺寸
	knobDiameter := width / 2.642 // 20px
	knobRadius := knobDiameter / 2
	padding := p.padding
	outOffset := width / 2.642 // 28px

	// 3. 使用animValue作为动画进度
	t := r.s.animValue

	// 4. 修正的位置计算
	// YES球：Checked时向左移出，未Checked时在左边
	yesVisiblePos := padding
	yesHiddenPos := -outOffset - knobDiameter
	yesX := yesVisiblePos
	if r.s.Checked {
		yesX = yesVisiblePos + (yesHiddenPos-yesVisiblePos)*t
	} else {
		yesX = yesHiddenPos + (yesVisiblePos-yesHiddenPos)*(1.0-t)
	}

	// NO球：Checked时在右边，未Checked时向右移出
	noVisiblePos := width - knobDiameter - padding
	noHiddenPos := width + outOffset
	noX := noHiddenPos
	if r.s.Checked {
		noX = noHiddenPos + (noVisiblePos-noHiddenPos)*t
	} else {
		noX = noVisiblePos + (noHiddenPos-noVisiblePos)*(1.0-t)
	}

	centerY := p.centerY

	// 5. 加载字体
	fontSize := p.fontSize
	err := dc.LoadFontFace(r.s.Config.FontPath, fontSize)
	if err != nil {
		dc.LoadFontFace("", fontSize)
	}

	// 6. 绘制YES球和文字
	if yesX+knobDiameter > -knobDiameter && yesX < width+knobDiameter {
		dc.SetColor(r.s.Config.YesColor)
		dc.DrawCircle(yesX+knobRadius, centerY, knobRadius)
		dc.Fill()

		// YES文字
		dc.SetColor(r.s.Config.TextColor)
		yesText := r.s.Config.YesLabel
		textWidth, textHeight := dc.MeasureString(yesText)
		textX := yesX + (knobDiameter-textWidth)/2
		textY := centerY + textHeight*0.3
		dc.DrawString(yesText, textX, textY)
	}

	// 7. 绘制NO球和文字
	if noX+knobDiameter > -knobDiameter && noX < width+knobDiameter {
		dc.SetColor(r.s.Config.NoColor)
		dc.DrawCircle(noX+knobRadius, centerY, knobRadius)
		dc.Fill()

		// NO文字
		dc.SetColor(r.s.Config.TextColor)
		noText := r.s.Config.NoLabel
		textWidth, textHeight := dc.MeasureString(noText)
		textX := noX + (knobDiameter-textWidth)/2
		textY := centerY + textHeight*0.3
		dc.DrawString(noText, textX, textY)
	}
}

// ==================== 双球纵向移动效果====================
func (r *toggleSwitchRenderer) drawEffect3VerticalSwap(dc *gg.Context) {
	p := r.s.drawCache
	width := p.width
	height := p.height

	// 1. 绘制背景
	bgColor := r.s.Config.YesBgColor
	if r.s.Checked {
		bgColor = r.s.Config.NoBgColor
	}
	dc.SetColor(bgColor)
	bgRadius := math.Min(height/2, 25.0)
	dc.DrawRoundedRectangle(0, 0, width, height, bgRadius)
	dc.Fill()

	// 2. 计算关键尺寸
	knobDiameter := width / 2.642 // 20px
	knobRadius := knobDiameter / 2
	padding := p.padding
	outOffset := width / 2.642 // 28px

	// 3. 使用animValue作为动画进度
	t := r.s.animValue

	// 4. 位置计算（YES在左边，NO在右边）
	yesVisibleY := padding
	yesHiddenY := -outOffset - knobDiameter
	yesX := padding // 左边4px
	yesY := yesVisibleY
	if r.s.Checked {
		yesY = yesVisibleY + (yesHiddenY-yesVisibleY)*t
	} else {
		yesY = yesHiddenY + (yesVisibleY-yesHiddenY)*(1.0-t)
	}

	noVisibleY := padding
	noHiddenY := -outOffset - knobDiameter
	noX := width - knobDiameter - padding // 右边4px
	noY := noHiddenY
	if r.s.Checked {
		noY = noHiddenY + (noVisibleY-noHiddenY)*t
	} else {
		noY = noVisibleY + (noHiddenY-noVisibleY)*(1.0-t)
	}

	// 5. 加载字体
	fontSize := p.fontSize
	err := dc.LoadFontFace(r.s.Config.FontPath, fontSize)
	if err != nil {
		dc.LoadFontFace("", fontSize)
	}

	// 6. 绘制YES球和文字
	if yesY+knobDiameter > -knobDiameter && yesY < height+knobDiameter {
		dc.SetColor(r.s.Config.YesColor)
		dc.DrawCircle(yesX+knobRadius, yesY+knobRadius, knobRadius)
		dc.Fill()

		// YES文字
		dc.SetColor(r.s.Config.TextColor)
		yesText := r.s.Config.YesLabel
		textWidth, textHeight := dc.MeasureString(yesText)
		textX := yesX + (knobDiameter-textWidth)/2
		textY := yesY + knobRadius + textHeight*0.3
		dc.DrawString(yesText, textX, textY)
	}

	// 7. 绘制NO球和文字
	if noY+knobDiameter > -knobDiameter && noY < height+knobDiameter {
		dc.SetColor(r.s.Config.NoColor)
		dc.DrawCircle(noX+knobRadius, noY+knobRadius, knobRadius)
		dc.Fill()

		// NO文字
		dc.SetColor(r.s.Config.TextColor)
		noText := r.s.Config.NoLabel
		textWidth, textHeight := dc.MeasureString(noText)
		textX := noX + (knobDiameter-textWidth)/2
		textY := noY + knobRadius + textHeight*0.3
		dc.DrawString(noText, textX, textY)
	}
}

// ==================== 第五个效果：精确模拟CSS ====================
func (r *toggleSwitchRenderer) drawEffect5ProjectionFlip(dc *gg.Context) {
	width := float64(dc.Width())
	height := float64(dc.Height())

	// 1. 绘制背景（背景也要旋转）
	// CSS中背景会旋转-180度
	bgAngle := -r.s.animValue * math.Pi
	bgColor := mixColor(r.s.Config.YesBgColor, r.s.Config.NoBgColor, r.s.animValue)

	// 保存当前状态
	dc.Push()

	// 平移坐标系到背景中心
	dc.Translate(width/2, height/2)
	// 旋转背景
	dc.Rotate(bgAngle)
	// 平移回去
	dc.Translate(-width/2, -height/2)

	dc.SetColor(bgColor)
	bgRadius := math.Min(height/2, 25.0)
	dc.DrawRoundedRectangle(0, 0, width, height, bgRadius)
	dc.Fill()

	// 恢复状态
	dc.Pop()

	// 2. 计算关键尺寸
	knobDiameter := width / 2.642 // 20px
	knobRadius := knobDiameter / 2

	// 3. 滑块位置：YES在左边4px，NO在右边4px
	// 所以滑块在左右两边移动，但同时旋转
	padding := width / 18.5 // 4px
	yesX := padding
	noX := width - knobDiameter - padding
	currentX := yesX + (noX-yesX)*r.s.animValue
	knobY := (height - knobDiameter) / 2

	// 4. 绘制滑块（滑块旋转+180度）
	knobAngle := r.s.animValue * math.Pi

	dc.Push()
	// 平移坐标系到滑块中心
	dc.Translate(currentX+knobRadius, knobY+knobRadius)
	// 旋转滑块
	dc.Rotate(knobAngle)

	// 根据角度选择颜色
	knobColor := mixColor(r.s.Config.YesColor, r.s.Config.NoColor, r.s.animValue)
	dc.SetColor(knobColor)
	dc.DrawCircle(0, 0, knobRadius)
	dc.Fill()

	// 恢复状态（在绘制文字前恢复，这样文字不会旋转）
	dc.Pop()

	// 5. 绘制文字（文字不旋转，保持正向）
	if r.s.ShowLabels {
		dc.SetColor(r.s.Config.TextColor)
		fontSize := width / 7.4
		err := dc.LoadFontFace(r.s.Config.FontPath, fontSize)
		if err != nil {
			dc.LoadFontFace("", fontSize)
		}

		text := r.s.Config.YesLabel
		if r.s.animValue > 0.5 {
			text = r.s.Config.NoLabel
		}

		textWidth, textHeight := dc.MeasureString(text)
		textX := currentX + (knobDiameter-textWidth)/2
		textY := knobY + knobRadius + textHeight*0.3
		dc.DrawString(text, textX, textY)
	}
}

// ==================== 第十个效果：优化版（精确对齐） ====================
func (r *toggleSwitchRenderer) drawEffect10RectSlide(dc *gg.Context) {
	p := r.s.drawCache
	width := p.width
	height := p.height

	// 1. 绘制背景
	bgColor := mixColor(r.s.Config.YesBgColor, r.s.Config.NoBgColor, r.s.animValue)
	dc.SetColor(bgColor)
	dc.DrawRoundedRectangle(0, 0, width, height, 2.0)
	dc.Fill()

	// 2. 正方形滑块
	knobSize := width / 3.7 // 20px
	padding := p.padding

	// 3. 加载字体
	fontSize := p.fontSize
	err := dc.LoadFontFace(r.s.Config.FontPath, fontSize)
	if err != nil {
		dc.LoadFontFace("", fontSize)
	}

	// 4. 测量文字
	yesText := r.s.Config.YesLabel
	noText := r.s.Config.NoLabel
	yesWidth, textHeight := dc.MeasureString(yesText)
	noWidth, _ := dc.MeasureString(noText)

	// 5. 计算关键位置
	leftLimit := padding
	rightLimit := width - knobSize - padding
	currentX := leftLimit + (rightLimit-leftLimit)*r.s.animValue
	knobY := (height - knobSize) / 2

	// 文字垂直位置（居中于滑块）
	textY := knobY + knobSize/2 + textHeight*0.3

	// 固定背景文字位置（左右两边，在滑块区域内居中）
	leftTextX := leftLimit + (knobSize-yesWidth)/2
	rightTextX := rightLimit + (knobSize-noWidth)/2

	// 6. 绘制背景文字（深色）
	dc.SetColor(r.s.Config.TextDarkColor)
	dc.DrawString(yesText, leftTextX, textY) // 左边的YES
	dc.DrawString(noText, rightTextX, textY) // 右边的NO

	// 7. 绘制正方形滑块
	knobColor := mixColor(r.s.Config.YesColor, r.s.Config.NoColor, r.s.animValue)
	dc.SetColor(knobColor)
	dc.DrawRoundedRectangle(currentX, knobY, knobSize, knobSize, 2.0)
	dc.Fill()

	// 8. 绘制滑块上的文字（白色）
	dc.SetColor(r.s.Config.TextColor)

	// 根据滑块位置决定显示哪个文字
	var sliderText string
	var sliderTextX float64

	if r.s.animValue < 0.5 {
		sliderText = r.s.Config.YesLabel
		sliderTextX = currentX + (knobSize-yesWidth)/2
	} else {
		sliderText = r.s.Config.NoLabel
		sliderTextX = currentX + (knobSize-noWidth)/2
	}

	// 绘制滑块文字
	dc.DrawString(sliderText, sliderTextX, textY)
}

// ==================== 第十二个效果：纯正方形滑块（无滑块文字） ====================
func (r *toggleSwitchRenderer) drawEffect12ColorBlockSwap(dc *gg.Context) {
	p := r.s.drawCache
	width := p.width
	height := p.height

	// 1. 绘制背景（矩形，2px圆角）
	bgColor := mixColor(r.s.Config.YesBgColor, r.s.Config.NoBgColor, r.s.animValue)
	dc.SetColor(bgColor)
	dc.DrawRoundedRectangle(0, 0, width, height, 2.0)
	dc.Fill()

	// 2. 正方形滑块
	sliderSize := width / 3.7 // 20px正方形
	padding := p.padding

	// 3. 滑块位置动画
	t := r.s.animValue
	leftPos := padding
	rightPos := width - sliderSize - padding
	sliderX := leftPos + (rightPos-leftPos)*t
	sliderY := (height - sliderSize) / 2

	// 4. 绘制固定文字
	if r.s.ShowLabels {
		fontSize := p.fontSize
		err := dc.LoadFontFace(r.s.Config.FontPath, fontSize)
		if err != nil {
			dc.LoadFontFace("", fontSize)
		}

		dc.SetColor(r.s.Config.TextDarkColor)

		yesText := r.s.Config.YesLabel
		noText := r.s.Config.NoLabel
		_, textHeight := dc.MeasureString(yesText)
		noWidth, _ := dc.MeasureString(noText)

		textY := height/2 + textHeight*0.3

		// YES在左边固定位置
		dc.DrawString(yesText, padding, textY)
		// NO在右边固定位置
		dc.DrawString(noText, width-padding-noWidth, textY)
	}

	// 5. 绘制正方形滑块（颜色随位置变化）
	sliderColor := mixColor(r.s.Config.YesColor, r.s.Config.NoColor, t)
	dc.SetColor(sliderColor)
	dc.DrawRoundedRectangle(sliderX, sliderY, sliderSize, sliderSize, 2.0)
	dc.Fill()
}

// ==================== 实用函数 ====================

// 设置透视距离
func (s *ToggleSwitch) SetPerspective(perspective float64) *ToggleSwitch {
	s.perspective = perspective
	s.Refresh()
	return s
}

// 获取当前效果的描述
func (s *ToggleSwitch) GetEffectName() string {
	switch s.Effect {
	case EffectSlide:
		return "Slide Effect"
	case EffectTwoBallSwap:
		return "Two Ball Swap"
	case EffectVerticalSwap:
		return "Vertical Swap"
	case EffectProjectionFlip:
		return "3D Projection Flip"
	case EffectRectSlide:
		return "Rect Slide"
	case EffectColorBlockSwap:
		return "Color Block Swap"
	default:
		return "Unknown Effect"
	}
}

// 设置自定义配置
func (s *ToggleSwitch) SetConfig(config SwitchConfig) *ToggleSwitch {
	s.Config = config
	s.calcDrawParams() // 配置变更时重新计算
	s.Refresh()
	return s
}

// 设置YES标签
func (s *ToggleSwitch) SetYesLabel(label string) *ToggleSwitch {
	s.Config.YesLabel = label
	s.Refresh()
	return s
}

// 设置NO标签
func (s *ToggleSwitch) SetNoLabel(label string) *ToggleSwitch {
	s.Config.NoLabel = label
	s.Refresh()
	return s
}

// 设置字体路径
func (s *ToggleSwitch) SetFontPath(path string) *ToggleSwitch {
	s.Config.FontPath = path
	s.Refresh()
	return s
}

// 设置YES颜色
func (s *ToggleSwitch) SetYesColor(color color.RGBA) *ToggleSwitch {
	s.Config.YesColor = color
	s.Refresh()
	return s
}

// 设置NO颜色
func (s *ToggleSwitch) SetNoColor(color color.RGBA) *ToggleSwitch {
	s.Config.NoColor = color
	s.Refresh()
	return s
}

// 设置YES背景颜色
func (s *ToggleSwitch) SetYesBgColor(color color.RGBA) *ToggleSwitch {
	s.Config.YesBgColor = color
	s.Refresh()
	return s
}

// 设置NO背景颜色
func (s *ToggleSwitch) SetNoBgColor(color color.RGBA) *ToggleSwitch {
	s.Config.NoBgColor = color
	s.Refresh()
	return s
}

// 设置文字颜色
func (s *ToggleSwitch) SetTextColor(color color.RGBA) *ToggleSwitch {
	s.Config.TextColor = color
	s.Refresh()
	return s
}

// 设置深色文字颜色
func (s *ToggleSwitch) SetTextDarkColor(color color.RGBA) *ToggleSwitch {
	s.Config.TextDarkColor = color
	s.Refresh()
	return s
}

// ==================== 新增布尔值自定义接口 ====================
func (s *ToggleSwitch) SetYesValue(val bool) *ToggleSwitch {
	s.Config.YesValue = val
	s.Refresh()
	return s
}

func (s *ToggleSwitch) SetNoValue(val bool) *ToggleSwitch {
	s.Config.NoValue = val
	s.Refresh()
	return s
}

// 获取当前值（返回自定义布尔值）
func (s *ToggleSwitch) Value() bool {
	if s.Checked {
		return s.Config.NoValue
	}
	return s.Config.YesValue
}

func (s *ToggleSwitch) SetChecked(checked bool) {
	s.Checked = checked
	s.Refresh()
}

func (s *ToggleSwitch) GetChecked() bool {
	return s.Checked
}

func (s *ToggleSwitch) Toggle() {
	s.Tapped(nil)
}

func (s *ToggleSwitch) SetOnChanged(f func(bool)) *ToggleSwitch {
	s.OnChanged = f
	return s
}

func (s *ToggleSwitch) SetEffect(effect SwitchEffect) *ToggleSwitch {
	s.Effect = effect
	s.Refresh()
	return s
}

func (s *ToggleSwitch) SetSize(width, height float64) *ToggleSwitch {
	s.Width = width
	s.Height = height
	s.calcDrawParams() // 尺寸变更时重新计算
	s.Refresh()
	return s
}

// ==================== 物理隔离容器 ====================
// WrapWithIsolationContainer 返回一个高度+2px的Stack包裹当前控件，实现物理隔离，防止canvas.Image刷新污染
func (s *ToggleSwitch) WrapWithIsolationContainer() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.Transparent)
	bg.SetMinSize(fyne.NewSize(float32(s.Width), float32(s.Height+2)))
	return container.NewStack(
		bg,
		container.NewCenter(s),
	)
}
