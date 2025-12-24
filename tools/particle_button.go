// 该文件专门写操作逻辑
/*


                                                    $$\     $$\
                                                    $$ |    \__|
 $$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\  $$$$$$\   $$\  $$$$$$$\  $$$$$$$\
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
	"math/rand"
	"time"

	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
	imagedraw "golang.org/x/image/draw"
)

// Particle 粒子结构体
type Particle struct {
	X, Y        float64 // 位置
	VX, VY      float64 // 速度
	StartSize   float64 // 起始大小
	CurrentSize float64 // 当前大小
	MaxSize     float64 // 最大大小
	Color       color.RGBA
	Life        float64 // 生命周期 (0-1)
	Decay       float64 // 衰减速度
	GrowthRate  float64 // 生长速度（先增大后减小）
	IsGrowing   bool    // 是否在生长阶段
}

// ParticleButton 自定义粒子按钮
type ParticleButton struct {
	widget.BaseWidget
	OnClick func()

	// 按钮配置
	Text   string
	Width  float32
	Height float32

	// 颜色配置
	BaseColor      color.RGBA   // 基础颜色
	GradientTop    color.RGBA   // 渐变顶部颜色
	GradientBottom color.RGBA   // 渐变底部颜色
	ShadowColor    color.RGBA   // 阴影颜色
	ParticleColors []color.RGBA // 粒子颜色数组

	// 粒子系统
	particles   []*Particle
	isAnimating bool
	isPressed   bool
	isHovered   bool

	// 按钮位置和尺寸
	width, height float32
	offsetY       float32 // 垂直偏移（用于悬停/按下效果）

	// 粒子动画控制
	// animationTicker *time.Ticker // 已废弃
	stopAnimation chan bool

	// 画布边框配置
	CanvasBorderWidth int  // 画布边框宽度，像素，默认1
	ShowCanvasBorder  bool // 是否显示画布边框，默认true

	// 字体渲染相关
	UseGGFont     bool        // 是否使用gg库字体渲染文本，默认false（用Fyne）
	GGFontType    string      // gg字体类型："chinese"或"english"，默认chinese
	GGFontSize    float64     // gg字体字号，默认20
	GGFontColor   color.Color // gg字体颜色，默认白色
	GGFontOffsetX float64     // gg字体X偏移，默认0
	GGFontOffsetY float64     // gg字体Y偏移，默认0

	// 自动变色功能
	AutoColorful bool // 启用后点击自动变色

	// 是否启用粒子特效（用户可动态修改）
	EnableParticle bool // true=有粒子，false=无粒子

	// 画布整体偏移
	CanvasOffsetX float64 // 整体X偏移
	CanvasOffsetY float64 // 整体Y偏移

	// 高分辨率修复开关
	HighDPI bool // 是否启用高分辨率修复，默认false
}

// 实现Mouseable接口以支持hover和按下下沉动画
var _ fyne.Widget = (*ParticleButton)(nil)
var _ fyne.Tappable = (*ParticleButton)(nil)
var _ desktop.Hoverable = (*ParticleButton)(nil)

// css-like 样式结构体
// 用户可集中配置所有样式属性
type ParticleButtonStyle struct {
	BaseColor     color.RGBA
	CanvasBorder  int
	UseGGFont     bool
	GGFontType    string
	GGFontSize    float64
	GGFontColor   color.Color
	GGFontOffsetX float64
	GGFontOffsetY float64
	AutoColorful  bool
	CanvasOffsetX float64
	CanvasOffsetY float64
}

// NewParticleButton 创建粒子按钮（默认绿色）
func NewParticleButton(onClick func()) *ParticleButton {
	baseColor := color.RGBA{R: 143, G: 196, B: 0, A: 255} // #8fc400
	return NewParticleButtonWithColor(onClick, baseColor, "按钮")
}

// NewParticleButtonWithColor 创建带自定义颜色的粒子按钮
func NewParticleButtonWithColor(onClick func(), baseColor color.RGBA, text string, opts ...func(*ParticleButton)) *ParticleButton {
	btn := &ParticleButton{
		OnClick:           onClick,
		Text:              text,
		Width:             120, // 默认宽度
		Height:            40,  // 默认高度
		BaseColor:         baseColor,
		particles:         make([]*Particle, 0),
		isAnimating:       false,
		isPressed:         false,
		isHovered:         false,
		width:             120,
		height:            40,
		offsetY:           0,
		stopAnimation:     make(chan bool),
		CanvasBorderWidth: 0,           // 默认无边框
		ShowCanvasBorder:  false,       // 默认不显示
		UseGGFont:         false,       // 默认用Fyne字体
		GGFontType:        "chinese",   // 默认chinese.ttf
		GGFontSize:        20,          // 默认字号
		GGFontColor:       color.White, // 默认白色
		GGFontOffsetX:     0,
		GGFontOffsetY:     0,
		AutoColorful:      false, // 默认不启用自动变色
		EnableParticle:    true,  // 默认启用粒子特效
	}

	// 应用可选参数
	for _, opt := range opts {
		opt(btn)
	}

	btn.generateColorsFromBase()
	btn.ExtendBaseWidget(btn)
	go btn.startAnimationLoop()
	return btn
}

// css-like 构造函数
func NewParticleButtonWithStyle(onClick func(), text string, style ParticleButtonStyle) *ParticleButton {
	btn := &ParticleButton{
		OnClick:           onClick,
		Text:              text,
		Width:             120,
		Height:            40,
		BaseColor:         style.BaseColor,
		particles:         make([]*Particle, 0),
		isAnimating:       false,
		isPressed:         false,
		isHovered:         false,
		width:             120,
		height:            40,
		offsetY:           0,
		stopAnimation:     make(chan bool),
		CanvasBorderWidth: style.CanvasBorder,
		ShowCanvasBorder:  style.CanvasBorder > 0,
		UseGGFont:         style.UseGGFont,
		GGFontType:        style.GGFontType,
		GGFontSize:        style.GGFontSize,
		GGFontColor:       style.GGFontColor,
		GGFontOffsetX:     style.GGFontOffsetX,
		GGFontOffsetY:     style.GGFontOffsetY,
		AutoColorful:      style.AutoColorful,
		CanvasOffsetX:     style.CanvasOffsetX,
		CanvasOffsetY:     style.CanvasOffsetY,
		EnableParticle:    true, // 默认启用粒子特效
	}
	if btn.GGFontType == "" {
		btn.GGFontType = "chinese"
	}
	if btn.GGFontSize == 0 {
		btn.GGFontSize = 20
	}
	if btn.GGFontColor == nil {
		btn.GGFontColor = color.White
	}
	btn.generateColorsFromBase()
	btn.ExtendBaseWidget(btn)
	go btn.startAnimationLoop()
	return btn
}

// 边框选项辅助函数
func WithCanvasBorder(width int) func(*ParticleButton) {
	return func(btn *ParticleButton) {
		btn.ShowCanvasBorder = true
		btn.CanvasBorderWidth = width
	}
}

// 可选参数：设置使用gg字体及类型、字号、颜色、偏移
func WithGGFont(fontType string) func(*ParticleButton) {
	return func(btn *ParticleButton) {
		btn.UseGGFont = true
		if fontType == "english" {
			btn.GGFontType = "english"
		} else {
			btn.GGFontType = "chinese"
		}
	}
}

func WithGGFontSize(size float64) func(*ParticleButton) {
	return func(btn *ParticleButton) {
		btn.GGFontSize = size
	}
}

func WithGGFontColor(c color.Color) func(*ParticleButton) {
	return func(btn *ParticleButton) {
		btn.GGFontColor = c
	}
}

func WithGGFontOffset(x, y float64) func(*ParticleButton) {
	return func(btn *ParticleButton) {
		btn.GGFontOffsetX = x
		btn.GGFontOffsetY = y
	}
}

// 可选参数：自动变色
func WithAutoColorful(enable bool) func(*ParticleButton) {
	return func(btn *ParticleButton) {
		btn.AutoColorful = enable
	}
}

// 可选参数：整体偏移
func WithCanvasOffset(x, y float64) func(*ParticleButton) {
	return func(btn *ParticleButton) {
		btn.CanvasOffsetX = x
		btn.CanvasOffsetY = y
	}
}

// generateColorsFromBase 从基础颜色生成其他所需颜色
func (btn *ParticleButton) generateColorsFromBase() {
	r, g, b, a := btn.BaseColor.R, btn.BaseColor.G, btn.BaseColor.B, btn.BaseColor.A

	// 计算渐变顶部颜色（比基础色亮一些）
	btn.GradientTop = color.RGBA{
		R: min(255, r+30),
		G: min(255, g+30),
		B: min(255, b+30),
		A: a,
	}

	// 计算渐变底部颜色（比基础色暗一些）
	btn.GradientBottom = color.RGBA{
		R: max(0, r-30),
		G: max(0, g-30),
		B: max(0, b-30),
		A: a,
	}

	// 计算阴影颜色（比基础色暗很多）
	btn.ShadowColor = color.RGBA{
		R: max(0, r-50),
		G: max(0, g-50),
		B: max(0, b-50),
		A: a,
	}

	// 生成粒子颜色数组（基于基础颜色的不同明暗度）
	btn.ParticleColors = []color.RGBA{
		// 基础色
		btn.BaseColor,
		// 稍微亮一点
		{R: min(255, r+20), G: min(255, g+20), B: min(255, b+20), A: a},
		// 稍微暗一点
		{R: max(0, r-20), G: max(0, g-20), B: max(0, b-20), A: a},
		// 更亮
		{R: min(255, r+40), G: min(255, g+40), B: min(255, b+40), A: a},
		// 更暗
		{R: max(0, r-40), G: max(0, g-40), B: max(0, b-40), A: a},
		// 中等亮度
		{R: min(255, r+10), G: min(255, g+10), B: min(255, b+10), A: a},
	}
}

// startAnimationLoop 启动粒子动画循环
func (btn *ParticleButton) startAnimationLoop() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			btn.UpdateParticles()
			if btn.isAnimating {
				fyne.Do(func() {
					btn.Refresh()
				})
			}
		case <-btn.stopAnimation:
			return
		}
	}
}

// Destroy 销毁按钮时停止动画
func (btn *ParticleButton) Destroy() {
	select {
	case btn.stopAnimation <- true:
	default:
	}
}

// SetSize 设置按钮尺寸
func (btn *ParticleButton) SetSize(width, height float32) {
	btn.Width = width
	btn.Height = height
	btn.width = width
	btn.height = height
	btn.Refresh()
}

// SetText 设置按钮文字
func (btn *ParticleButton) SetText(text string) {
	btn.Text = text
	btn.Refresh()
}

// SetBaseColor 设置基础颜色并重新生成其他颜色
func (btn *ParticleButton) SetBaseColor(baseColor color.RGBA) {
	btn.BaseColor = baseColor
	btn.generateColorsFromBase()
	btn.Refresh()
}

// SetColors 直接设置所有颜色
func (btn *ParticleButton) SetColors(base, top, bottom, shadow color.RGBA, particles []color.RGBA) {
	btn.BaseColor = base
	btn.GradientTop = top
	btn.GradientBottom = bottom
	btn.ShadowColor = shadow
	if particles != nil {
		btn.ParticleColors = particles
	}
	btn.Refresh()
}

// SetHighDPI 设置是否启用高分辨率修复
func (btn *ParticleButton) SetHighDPI(enable bool) {
	btn.HighDPI = enable
	btn.Refresh()
}

// CreateParticles 创建新粒子
func (btn *ParticleButton) CreateParticles(x, y float64, count int) {
	if !btn.EnableParticle {
		return
	}
	for i := 0; i < count; i++ {
		// 随机方向角度
		angle := rand.Float64() * 2 * math.Pi

		// 适中速度
		speed := 1.5 + rand.Float64()*3.0

		// 计算速度分量
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		// 随机起始大小
		startSize := 2.0 + rand.Float64()*4.0

		// 随机最大大小
		maxSize := 10 + rand.Float64()*12

		// 随机颜色（从粒子颜色数组中选取）
		colorIdx := rand.Intn(len(btn.ParticleColors))

		particle := &Particle{
			X:           x,
			Y:           y,
			VX:          vx,
			VY:          vy,
			StartSize:   startSize,
			CurrentSize: startSize,
			MaxSize:     maxSize,
			Color:       btn.ParticleColors[colorIdx],
			Life:        1.0,
			Decay:       0.01 + rand.Float64()*0.02,
			GrowthRate:  0.2 + rand.Float64()*0.3,
			IsGrowing:   true,
		}

		btn.particles = append(btn.particles, particle)
	}
	btn.isAnimating = true
}

// UpdateParticles 更新粒子状态
func (btn *ParticleButton) UpdateParticles() {
	if !btn.EnableParticle {
		btn.particles = nil
		btn.isAnimating = false
		return
	}
	aliveParticles := make([]*Particle, 0)

	for _, p := range btn.particles {
		p.X += p.VX
		p.Y += p.VY

		// 轻微阻力
		p.VX *= 0.98
		p.VY *= 0.98

		// 添加重力
		p.VY += 0.15

		if p.IsGrowing {
			p.CurrentSize += p.GrowthRate
			if p.CurrentSize >= p.MaxSize {
				p.CurrentSize = p.MaxSize
				p.IsGrowing = false
			}
		} else {
			p.Life -= p.Decay
			p.CurrentSize = math.Max(0.5, p.CurrentSize*p.Life)
		}

		if p.Life > 0.05 && p.CurrentSize > 0.5 {
			aliveParticles = append(aliveParticles, p)
		}
	}

	btn.particles = aliveParticles
	btn.isAnimating = len(btn.particles) > 0
}

// DrawParticles 绘制粒子到指定上下文
func (btn *ParticleButton) DrawParticles(dc *gg.Context, offsetX, offsetY float64) {
	for _, p := range btn.particles {
		// 根据生命周期调整透明度
		alpha := uint8(p.Life * 255)
		colorWithAlpha := color.RGBA{
			R: p.Color.R,
			G: p.Color.G,
			B: p.Color.B,
			A: alpha,
		}
		dc.SetColor(colorWithAlpha)
		dc.DrawCircle(p.X+offsetX, p.Y+offsetY, p.CurrentSize)
		dc.Fill()
	}
}

// Hover事件：鼠标移入
func (btn *ParticleButton) MouseIn(*desktop.MouseEvent) {
	btn.isHovered = true
	if !btn.isPressed {
		btn.offsetY = 2 // hover时下沉一半
		btn.Refresh()
	}
}

// Hover事件：鼠标移出
func (btn *ParticleButton) MouseOut() {
	btn.isHovered = false
	if !btn.isPressed {
		btn.offsetY = 0 // 恢复原状
		btn.Refresh()
	}
}

// 鼠标移动（可选实现，不处理）
func (btn *ParticleButton) MouseMoved(*desktop.MouseEvent) {}

// Tapped 处理点击事件
func (btn *ParticleButton) Tapped(e *fyne.PointEvent) {
	btn.isPressed = true
	if btn.isHovered {
		btn.offsetY = 4 // hover时点击，下沉到最大
	} else {
		btn.offsetY = 4 // 非hover点击也下沉
	}
	btn.Refresh()

	// 在按钮中心位置创建粒子
	btn.CreateParticles(float64(btn.width/2), float64(btn.height/2)+float64(btn.offsetY), 25+rand.Intn(10))

	// 自动变色功能
	if btn.AutoColorful {
		newColor := color.RGBA{
			R: uint8(50 + rand.Intn(150)),
			G: uint8(50 + rand.Intn(150)),
			B: uint8(50 + rand.Intn(150)),
			A: 255,
		}
		btn.SetBaseColor(newColor)
	}

	if btn.OnClick != nil {
		btn.OnClick()
	}

	// 延迟恢复按钮状态
	go func() {
		time.Sleep(150 * time.Millisecond)
		fyne.Do(func() {
			btn.isPressed = false
			if btn.isHovered {
				btn.offsetY = 2 // hover时松开，回到hover下沉
			} else {
				btn.offsetY = 0 // 非hover松开，回到原状
			}
			btn.Refresh()
		})
	}()
}

// IsAnimating 获取动画状态
func (btn *ParticleButton) IsAnimating() bool {
	return btn.isAnimating
}

// CreateRenderer 创建渲染器
func (btn *ParticleButton) CreateRenderer() fyne.WidgetRenderer {
	mainCanvas := canvas.NewImageFromImage(nil)
	mainCanvas.FillMode = canvas.ImageFillOriginal

	text := canvas.NewText(btn.Text, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 14

	// Stack比画布大2px，主画布和文本都放Stack里
	stack := container.NewWithoutLayout(mainCanvas, text)
	// Stack初始尺寸
	stack.Resize(fyne.NewSize(btn.Width+particleButtonStackExtra, btn.Height+particleButtonStackExtra+particleButtonPadding*2))

	renderer := &particleButtonRenderer{
		btn:        btn,
		mainCanvas: mainCanvas,
		text:       text,
		stack:      stack,
		objects:    []fyne.CanvasObject{stack},
	}

	// 关键修复：立即刷新一次，确保初始显示
	renderer.Refresh()

	return renderer
}

// particleButtonRenderer 渲染器
type particleButtonRenderer struct {
	btn        *ParticleButton
	mainCanvas *canvas.Image
	text       *canvas.Text
	stack      *fyne.Container // 新增：物理隔离Stack
	objects    []fyne.CanvasObject
}

const (
	particleButtonPadding    = 8 // 统一的padding，画布与按钮四周的留白
	particleButtonStackExtra = 2 // Stack比画布大2px
)

func (r *particleButtonRenderer) Destroy() {
	r.btn.Destroy()
}

func (r *particleButtonRenderer) Layout(size fyne.Size) {
	// Stack比按钮画布大2px
	canvasW := size.Width - particleButtonStackExtra
	canvasH := size.Height - particleButtonStackExtra
	r.btn.width = canvasW
	r.btn.height = canvasH
	r.mainCanvas.Resize(fyne.NewSize(canvasW, canvasH))

	// 文本居中
	textSize := r.text.MinSize()
	textX := (canvasW - textSize.Width) / 2
	textY := (canvasH - textSize.Height) / 2
	r.text.Move(fyne.NewPos(textX, textY))
	r.text.Resize(textSize)

	// Stack本身自动布局
	r.stack.Resize(size)

	// 布局改变后也需要刷新
	r.Refresh()
}

func (r *particleButtonRenderer) MinSize() fyne.Size {
	// Stack比按钮大2px
	return fyne.NewSize(r.btn.Width+particleButtonStackExtra, r.btn.Height+particleButtonStackExtra+particleButtonPadding*2)
}

func (r *particleButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.stack}
}

// drawCSSButton 绘制CSS按钮（只绘制背景和粒子）
func (r *particleButtonRenderer) drawCSSButton() *gg.Context {
	// 高分辨率修复：scale=2.0，否则1.0
	scale := 1.0
	if r.btn.HighDPI {
		scale = 2.0
	}
	width := int(float64(r.btn.width) * scale)
	height := int(float64(r.btn.height) * scale)
	padding := int(float64(particleButtonPadding) * scale)

	// 创建足够大的画布以容纳阴影和粒子
	canvasWidth := width + padding*2
	canvasHeight := height + padding*2

	dc := gg.NewContext(canvasWidth, canvasHeight)

	// 透明背景（按钮是独立的）
	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	// 先整体平移，仅影响主体内容
	dc.Push()
	dc.Translate(r.btn.CanvasOffsetX*scale, r.btn.CanvasOffsetY*scale)

	// 按钮在画布中的位置（padding留白）
	buttonX := float64(padding)
	buttonY := float64(padding)

	// 根据按钮状态设置偏移
	buttonOffsetY := float64(r.btn.offsetY) * scale

	// 绘制按钮阴影（多层CSS阴影效果）
	if r.btn.isPressed {
		// 按下状态：没有外阴影，只有内阴影
		dc.SetRGBA(1, 1, 1, 0.2)
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY, float64(width), 2*scale, 7*scale)
		dc.Fill()

	} else if r.btn.isHovered {
		// 悬停状态：较少阴影
		// 黑色外阴影（2px模糊）
		dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 180})
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY+2*scale, float64(width), float64(height), 8*scale)
		dc.Fill()

		// 两层阴影
		for i := 0; i < 2; i++ {
			dc.SetColor(r.btn.ShadowColor)
			dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY+float64(i+1)*scale, float64(width), float64(height), 8*scale)
			dc.Stroke()
		}

		// 内阴影（白色高光）
		dc.SetRGBA(1, 1, 1, 0.2)
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY, float64(width), 2*scale, 7*scale)
		dc.Fill()

	} else {
		// 正常状态：完整的多层阴影
		// 黑色外阴影
		dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 180})
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY+4*scale, float64(width), float64(height), 8*scale)
		dc.Fill()

		// 四层阴影
		for i := 0; i < 4; i++ {
			dc.SetColor(r.btn.ShadowColor)
			dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY+float64(i+1)*scale, float64(width), float64(height), 8*scale)
			dc.Stroke()
		}

		// 内阴影（白色高光）
		dc.SetRGBA(1, 1, 1, 0.2)
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY, float64(width), 2*scale, 7*scale)
		dc.Fill()
	}

	// 绘制按钮主体（渐变）
	gradient := gg.NewLinearGradient(
		buttonX, buttonY+buttonOffsetY,
		buttonX, buttonY+buttonOffsetY+float64(height),
	)
	gradient.AddColorStop(0, r.btn.GradientTop)
	gradient.AddColorStop(1, r.btn.GradientBottom)

	dc.SetFillStyle(gradient)
	dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY, float64(width), float64(height), 8*scale)
	dc.Fill()

	// 绘制粒子（在按钮之上）
	if len(r.btn.particles) > 0 {
		// 粒子坐标和大小也要乘以scale
		for _, p := range r.btn.particles {
			alpha := uint8(p.Life * 255)
			colorWithAlpha := color.RGBA{
				R: p.Color.R,
				G: p.Color.G,
				B: p.Color.B,
				A: alpha,
			}
			dc.SetColor(colorWithAlpha)
			dc.DrawCircle(p.X*scale+buttonX, p.Y*scale+buttonY+buttonOffsetY, p.CurrentSize*scale)
			dc.Fill()
		}
	}

	// ====== 用gg库绘制文字（支持自定义字体/字号/颜色/偏移） ======
	if r.btn.UseGGFont {
		fontPath := "ttf/chinese.ttf"
		if r.btn.GGFontType == "english" {
			fontPath = "ttf/english.ttf"
		}
		fontSize := r.btn.GGFontSize * scale
		if fontSize <= 0 {
			fontSize = 20 * scale
		}
		fontColor := r.btn.GGFontColor
		if fontColor == nil {
			fontColor = color.White
		}
		if err := dc.LoadFontFace(fontPath, fontSize); err == nil {
			dc.SetColor(fontColor)
			centerX := buttonX + float64(width)/2 + r.btn.GGFontOffsetX*scale
			centerY := buttonY + buttonOffsetY + float64(height)/2 + r.btn.GGFontOffsetY*scale
			dc.DrawStringAnchored(r.btn.Text, centerX, centerY, 0.5, 0.5)
		}
	}
	dc.Pop() // 恢复，不影响边框

	// 画布最外层边框（可选，始终包裹原始区域，不受偏移影响）
	if r.btn.ShowCanvasBorder && r.btn.CanvasBorderWidth > 0 {
		dc.SetLineWidth(float64(r.btn.CanvasBorderWidth) * scale)
		dc.SetColor(color.Black) // 你可以改成可配置
		dc.DrawRectangle(0.5*float64(r.btn.CanvasBorderWidth)*scale, 0.5*float64(r.btn.CanvasBorderWidth)*scale, float64(canvasWidth)-float64(r.btn.CanvasBorderWidth)*scale, float64(canvasHeight)-float64(r.btn.CanvasBorderWidth)*scale)
		dc.Stroke()
	}

	// 如果高分辨率，缩放回原始尺寸（高质量）
	if scale > 1.0 {
		img := dc.Image()
		dstW := int(r.btn.width) + particleButtonPadding*2
		dstH := int(r.btn.height) + particleButtonPadding*2
		dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
		imagedraw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), imagedraw.Over, nil)
		resized := gg.NewContext(dstW, dstH)
		resized.DrawImage(dst, 0, 0)
		return resized
	}
	return dc
}

func (r *particleButtonRenderer) Refresh() {
	// 更新文本内容
	r.text.Text = r.btn.Text

	// 如果用gg字体，隐藏Fyne文本，否则显示
	if r.btn.UseGGFont {
		if col, ok := r.text.Color.(color.RGBA); ok {
			col.A = 0
			r.text.Color = col
		}
	} else {
		if col, ok := r.text.Color.(color.RGBA); ok {
			col.A = 255
			r.text.Color = col
		} else {
			r.text.Color = color.RGBA{255, 255, 255, 255}
		}
	}

	// 绘制按钮图像
	dc := r.drawCSSButton()
	r.mainCanvas.Image = dc.Image()
	r.mainCanvas.Refresh()
	r.text.Refresh()
}

// 辅助函数
func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
