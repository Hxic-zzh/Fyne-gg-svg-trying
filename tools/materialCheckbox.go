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
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/font"
)

// ==================== 颜色常量 ====================
var (
	// 从CSS中提取的颜色（默认值）
	CheckboxColorDefault  = color.RGBA{181, 191, 217, 255} // #b5bfd9
	CheckboxColorSelected = color.RGBA{34, 96, 255, 255}   // #2260ff
	CheckboxColorHover    = color.RGBA{34, 96, 255, 255}   // #2260ff
	CheckboxColorText     = color.RGBA{112, 112, 112, 255} // #707070
	CheckboxColorTextSel  = color.RGBA{34, 96, 255, 255}   // #2260ff

	CheckboxColorShadow = color.RGBA{0, 0, 0, 25}     // 阴影颜色
	CheckboxColorGlow   = color.RGBA{34, 96, 255, 50} // 泛光颜色
)

// ==================== SVG图标缓存（支持重新着色）====================
type svgCacheEntry struct {
	image       image.Image
	size        float64
	baseColor   color.RGBA // 原始颜色
	targetColor color.RGBA // 目标颜色
}

var (
	svgCache   = make(map[string]svgCacheEntry)
	svgCacheMu sync.RWMutex
)

// colorToString 将颜色转换为字符串用于缓存键
func colorToString(c color.Color) string {
	if c == nil {
		return "nil"
	}
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("%d_%d_%d_%d", r>>8, g>>8, b>>8, a>>8)
}

// recolorImage 重新着色图像
func recolorImage(img image.Image, newColor color.Color) image.Image {
	bounds := img.Bounds()
	recolored := image.NewRGBA(bounds)

	// 将 newColor 转换为 RGBA
	var r, g, b uint8
	if rgba, ok := newColor.(color.RGBA); ok {
		r, g, b = rgba.R, rgba.G, rgba.B
	} else {
		// 如果无法转换，使用默认值
		r, g, b = 0, 0, 0
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				// 保留原始Alpha通道，只改变颜色
				recolored.Set(x, y, color.RGBA{
					R: r,
					G: g,
					B: b,
					A: uint8(a >> 8),
				})
			} else {
				recolored.Set(x, y, color.Transparent)
			}
		}
	}

	return recolored
}

// loadAndRenderSVG 从文件路径加载并渲染SVG（支持重新着色）
func loadAndRenderSVG(svgPath string, size float64, targetColor color.Color) (image.Image, error) {
	if svgPath == "" {
		return nil, nil
	}

	// 创建缓存键（包含颜色）
	cacheKey := fmt.Sprintf("%s_%.0f_%s", svgPath, size, colorToString(targetColor))

	svgCacheMu.RLock()
	if entry, exists := svgCache[cacheKey]; exists {
		svgCacheMu.RUnlock()
		return entry.image, nil
	}
	svgCacheMu.RUnlock()

	// 读取SVG文件
	svgData, err := os.ReadFile(svgPath)
	if err != nil {
		return nil, err
	}

	// 解析SVG
	icon, err := oksvg.ReadIconStream(bytes.NewBuffer(svgData))
	if err != nil {
		return nil, err
	}

	// 设置目标尺寸
	icon.SetTarget(0, 0, size, size)

	// 创建图像
	img := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))

	// 创建扫描仪
	scanner := rasterx.NewScannerGV(int(size), int(size), img, img.Bounds())
	raster := rasterx.NewDasher(int(size), int(size), scanner)

	// 绘制SVG
	icon.Draw(raster, 1.0)

	// 重新着色
	var finalImage image.Image = img
	if targetColor != nil {
		finalImage = recolorImage(img, targetColor)
	}

	// 缓存结果
	svgCacheMu.Lock()
	entry := svgCacheEntry{
		image: finalImage,
		size:  size,
	}
	// 尝试转换为 RGBA 存储
	if rgba, ok := targetColor.(color.RGBA); ok {
		entry.targetColor = rgba
	}
	svgCache[cacheKey] = entry
	svgCacheMu.Unlock()

	return finalImage, nil
}

// ==================== 布局参数 ====================
type CheckboxLayout struct {
	TileWidth     float64 // 卡片宽度
	TileHeight    float64 // 卡片高度
	CheckSize     float64 // 选中标记尺寸（左上角小圆点）
	Padding       float64 // 内边距
	BorderRadius  float64 // 圆角半径
	IconSize      float64 // 图标尺寸
	LabelFontSize float64 // 文字大小
	Spacing       float64 // 图标与文字间距
}

// ==================== 样式结构体 ====================
type MaterialCheckboxStyle struct {
	TileWidth    float64     // 卡片宽度
	TileHeight   float64     // 卡片高度
	IconColor    color.Color // 图标颜色
	LabelColor   color.Color // 文字颜色
	BorderColor  color.Color // 边框颜色
	BgColor      color.Color // 背景色
	ShadowColor  color.Color // 阴影颜色
	CornerRadius float64     // 圆角
	IconSize     float64     // 图标尺寸
	FontSize     float64     // 字体大小
	IconPath     string      // SVG文件路径

	// 新增：动画颜色（外部可覆盖）
	HoverColor     color.Color // hover状态颜色
	SelectedColor  color.Color // 选中状态颜色
	CircleColor    color.Color // 圆圈颜色
	CheckmarkColor color.Color // 对勾颜色
}

// ==================== 自定义Checkbox ====================
type MaterialCheckbox struct {
	widget.BaseWidget

	// 核心状态
	checked bool
	hovered bool

	// 内容
	text string

	// 样式配置
	Style  MaterialCheckboxStyle
	layout CheckboxLayout

	// 动画参数（修复版本）
	hoverProgress  float64 // hover动画进度 (0-1)
	checkProgress  float64 // 选中动画进度 (0-1)
	targetHover    float64 // 目标hover进度
	targetCheck    float64 // 目标选中进度
	scaleProgress  float64 // 缩放动画进度 (1.0-1.02)
	targetScale    float64 // 目标缩放
	circleProgress float64 // 圆圈动画进度 (0-1)

	// 渲染组件
	canvasImg *canvas.Image
	container *fyne.Container

	// 实际尺寸
	width, height float64

	// 动画控制
	animTicker *time.Ticker
	animStop   chan bool
	animMutex  sync.Mutex

	// 回调函数
	OnChanged func(bool)

	// 字体缓存
	fontPath       string
	cachedFontFace font.Face
	cachedFontSize float64

	// SVG图标缓存
	svgImage     image.Image
	lastSVGSize  float64
	lastSVGPath  string
	lastSVGColor color.Color
}

// ==================== 创建方法 ====================
func NewMaterialCheckbox(text string, checked bool, tileSize ...float64) *MaterialCheckbox {
	c := &MaterialCheckbox{
		text:          text,
		checked:       checked,
		hovered:       false,
		hoverProgress: 0,
		targetHover:   0,
		targetScale:   1.0,
		scaleProgress: 1.0,
		fontPath:      "ttf/english.ttf",
		animStop:      make(chan bool, 1),
	}

	// 设置选中状态对应的进度
	if checked {
		c.checkProgress = 1.0
		c.targetCheck = 1.0
		c.circleProgress = 1.0
	} else {
		c.checkProgress = 0
		c.targetCheck = 0
		c.circleProgress = 0
	}

	// 根据用户输入设置默认尺寸
	var tileWidth, tileHeight float64 = 112, 112 // 默认7rem = 112px
	if len(tileSize) > 0 && tileSize[0] > 0 {
		tileWidth = tileSize[0]
	}
	if len(tileSize) > 1 && tileSize[1] > 0 {
		tileHeight = tileSize[1]
	}

	// 默认样式（基于CSS设计） - 硬编码颜色作为默认值
	c.Style = MaterialCheckboxStyle{
		TileWidth:    tileWidth,
		TileHeight:   tileHeight,
		IconColor:    CheckboxColorText,    // 默认灰色
		LabelColor:   CheckboxColorText,    // 默认灰色
		BorderColor:  CheckboxColorDefault, // 默认边框颜色
		BgColor:      color.White,
		ShadowColor:  CheckboxColorShadow,
		CornerRadius: 8, // 0.5rem = 8px

		// 设置动画颜色默认值（硬编码）
		HoverColor:     CheckboxColorHover,
		SelectedColor:  CheckboxColorSelected,
		CircleColor:    CheckboxColorSelected, // 圆圈颜色与选中色相同
		CheckmarkColor: color.White,           // 对勾为白色
	}

	// 计算布局
	c.calculateLayout()

	// 创建画布
	c.canvasImg = canvas.NewImageFromResource(nil)
	c.canvasImg.FillMode = canvas.ImageFillContain

	// 创建容器
	c.container = container.NewWithoutLayout(c.canvasImg)

	// 初始化
	c.ExtendBaseWidget(c)

	// 启动动画
	c.startAnimation()

	// 初始渲染
	fyne.Do(func() {
		c.updateCanvas()
	})

	return c
}

// ==================== 布局计算 ====================
func (c *MaterialCheckbox) calculateLayout() {
	// 使用CSS的7rem = 112px作为基础
	tileWidth := c.Style.TileWidth
	tileHeight := c.Style.TileHeight

	// 如果用户没有指定，使用默认值
	if tileWidth <= 0 {
		tileWidth = 112
	}
	if tileHeight <= 0 {
		tileHeight = 112
	}

	// 设置布局参数（基于CSS设计）
	c.layout.TileWidth = tileWidth
	c.layout.TileHeight = tileHeight
	c.layout.CheckSize = 20 // 左上角选中标记尺寸 1.25rem = 20px
	c.layout.Padding = 4    // 0.25rem = 4px
	c.layout.BorderRadius = c.Style.CornerRadius

	// 图标尺寸：3rem = 48px
	if c.Style.IconSize > 0 {
		c.layout.IconSize = c.Style.IconSize
	} else {
		c.layout.IconSize = 48
	}

	// 字体大小：根据卡片高度比例计算
	if c.Style.FontSize > 0 {
		c.layout.LabelFontSize = c.Style.FontSize
	} else {
		c.layout.LabelFontSize = 14
	}

	c.layout.Spacing = 8 // 图标与文字间距

	// 整体控件尺寸 = 卡片高度 + 文字区域
	c.width = tileWidth
	c.height = tileHeight + c.layout.LabelFontSize + c.layout.Spacing
}

// ==================== SVG图标方法 ====================
func (c *MaterialCheckbox) updateSVGIcon() {
	if c.Style.IconPath == "" {
		c.svgImage = nil
		return
	}

	// 计算图标颜色 - 使用样式中的SelectedColor
	var iconColor color.Color = c.Style.IconColor
	if c.checked {
		iconColor = c.getSelectedColor()
	}

	// 如果路径、尺寸和颜色没有变化，不需要重新渲染
	if c.svgImage != nil &&
		c.lastSVGSize == c.layout.IconSize &&
		c.lastSVGPath == c.Style.IconPath &&
		colorToString(c.lastSVGColor) == colorToString(iconColor) {
		return
	}

	// 加载并渲染SVG
	img, err := loadAndRenderSVG(c.Style.IconPath, c.layout.IconSize, iconColor)
	if err != nil {
		c.svgImage = nil
		return
	}

	c.svgImage = img
	c.lastSVGSize = c.layout.IconSize
	c.lastSVGPath = c.Style.IconPath
	c.lastSVGColor = iconColor
}

// ==================== 设置方法 ====================
func (c *MaterialCheckbox) SetIconPath(svgPath string) {
	c.Style.IconPath = svgPath
	c.svgImage = nil // 清除缓存，强制重新渲染
	c.lastSVGPath = ""
	c.lastSVGColor = nil
	fyne.Do(func() {
		c.updateCanvas()
		c.canvasImg.Refresh()
	})
}

func (c *MaterialCheckbox) SetStyle(style MaterialCheckboxStyle) {
	// 保留用户输入的尺寸
	if style.TileWidth > 0 {
		c.Style.TileWidth = style.TileWidth
	}
	if style.TileHeight > 0 {
		c.Style.TileHeight = style.TileHeight
	}

	// 更新其他样式属性
	if style.IconColor != nil {
		c.Style.IconColor = style.IconColor
	}
	if style.LabelColor != nil {
		c.Style.LabelColor = style.LabelColor
	}
	if style.BorderColor != nil {
		c.Style.BorderColor = style.BorderColor
	}
	if style.BgColor != nil {
		c.Style.BgColor = style.BgColor
	}
	if style.ShadowColor != nil {
		c.Style.ShadowColor = style.ShadowColor
	}
	if style.CornerRadius > 0 {
		c.Style.CornerRadius = style.CornerRadius
	}
	if style.IconSize > 0 {
		c.Style.IconSize = style.IconSize
	}
	if style.FontSize > 0 {
		c.Style.FontSize = style.FontSize
	}
	if style.IconPath != "" {
		c.Style.IconPath = style.IconPath
		c.svgImage = nil // 清除缓存
		c.lastSVGPath = ""
		c.lastSVGColor = nil
	}

	// 更新动画颜色属性（如果提供了则覆盖默认值）
	if style.HoverColor != nil {
		c.Style.HoverColor = style.HoverColor
	}
	if style.SelectedColor != nil {
		c.Style.SelectedColor = style.SelectedColor
	}
	if style.CircleColor != nil {
		c.Style.CircleColor = style.CircleColor
	}
	if style.CheckmarkColor != nil {
		c.Style.CheckmarkColor = style.CheckmarkColor
	}

	// 重新计算布局
	c.calculateLayout()

	// 触发更新
	fyne.Do(func() {
		c.updateCanvas()
		c.canvasImg.Refresh()
	})
}

func (c *MaterialCheckbox) SetFontPath(path string) {
	c.fontPath = path
	c.cachedFontFace = nil // 清除字体缓存
	fyne.Do(func() {
		c.updateCanvas()
		c.canvasImg.Refresh()
	})
}

// ==================== 动画系统（修复版本）====================
func (c *MaterialCheckbox) startAnimation() {
	c.animTicker = time.NewTicker(16 * time.Millisecond) // ~60fps

	go func() {
		for {
			select {
			case <-c.animTicker.C:
				c.updateAnimation()
			case <-c.animStop:
				if c.animTicker != nil {
					c.animTicker.Stop()
					c.animTicker = nil
				}
				return
			}
		}
	}()
}

func (c *MaterialCheckbox) stopAnimation() {
	if c.animStop != nil {
		select {
		case c.animStop <- true:
		default:
		}
	}
}

func (c *MaterialCheckbox) updateAnimation() {
	c.animMutex.Lock()
	defer c.animMutex.Unlock()

	needsUpdate := false

	// 1. 选中动画
	oldCheck := c.checkProgress
	if c.checkProgress < c.targetCheck {
		c.checkProgress += 0.2
		if c.checkProgress > c.targetCheck {
			c.checkProgress = c.targetCheck
		}
	} else if c.checkProgress > c.targetCheck {
		c.checkProgress -= 0.2
		if c.checkProgress < c.targetCheck {
			c.checkProgress = c.targetCheck
		}
	}
	if c.absFloat(oldCheck-c.checkProgress) > 0.001 {
		needsUpdate = true
	}

	// 2. Hover动画
	oldHover := c.hoverProgress
	if c.hoverProgress < c.targetHover {
		c.hoverProgress += 0.3 // 进入时较快
		if c.hoverProgress > c.targetHover {
			c.hoverProgress = c.targetHover
		}
	} else if c.hoverProgress > c.targetHover {
		c.hoverProgress -= 0.15 // 退出时较慢
		if c.hoverProgress < c.targetHover {
			c.hoverProgress = c.targetHover
		}
	}
	if c.absFloat(oldHover-c.hoverProgress) > 0.001 {
		needsUpdate = true
	}

	// 3. 圆圈动画进度（CSS中的:before伪元素）
	oldCircle := c.circleProgress
	if c.hovered || c.checked {
		// hover或选中时圆圈完全显示
		targetCircle := 1.0
		if c.circleProgress < targetCircle {
			c.circleProgress += 0.25
			if c.circleProgress > targetCircle {
				c.circleProgress = targetCircle
			}
		}
	} else {
		// 未hover且未选中时圆圈消失
		if c.circleProgress > 0 {
			c.circleProgress -= 0.15
			if c.circleProgress < 0 {
				c.circleProgress = 0
			}
		}
	}
	if c.absFloat(oldCircle-c.circleProgress) > 0.001 {
		needsUpdate = true
	}

	// 4. 缩放动画（hover时有轻微放大效果）
	oldScale := c.scaleProgress
	if c.hoverProgress > 0 {
		c.targetScale = 1.02 // CSS中的轻微放大效果
	} else {
		c.targetScale = 1.0
	}

	if c.scaleProgress < c.targetScale {
		c.scaleProgress += 0.1
		if c.scaleProgress > c.targetScale {
			c.scaleProgress = c.targetScale
		}
	} else if c.scaleProgress > c.targetScale {
		c.scaleProgress -= 0.1
		if c.scaleProgress < c.targetScale {
			c.scaleProgress = c.targetScale
		}
	}
	if c.absFloat(oldScale-c.scaleProgress) > 0.001 {
		needsUpdate = true
	}

	// 更新UI
	if needsUpdate {
		fyne.Do(func() {
			c.updateCanvas()
			c.canvasImg.Refresh()
		})
	}
}

// ==================== 画布绘制（修复选中标记动画 - CSS风格）====================
func (c *MaterialCheckbox) updateCanvas() {
	// 重新计算布局
	c.calculateLayout()

	// 更新SVG图标
	c.updateSVGIcon()

	// 确保最小尺寸
	minSize := 40.0
	if c.width < minSize {
		c.width = minSize
	}
	if c.height < minSize {
		c.height = minSize
	}

	// 创建画布
	dc := gg.NewContext(int(c.width), int(c.height))
	if dc == nil {
		return
	}

	// 清空背景
	dc.SetColor(color.Transparent)
	dc.Clear()

	// ==================== 绘制卡片 ====================

	// 应用缩放变换
	dc.Translate(c.width/2, c.height/2)
	dc.Scale(c.scaleProgress, c.scaleProgress)
	dc.Translate(-c.width/2, -c.height/2)

	// 卡片位置（居中）
	cardX := (c.width - c.layout.TileWidth) / 2
	cardY := float64(0)
	cardWidth := c.layout.TileWidth
	cardHeight := c.layout.TileHeight

	// 绘制卡片阴影（hover时显示，像CSS的box-shadow）
	if c.hoverProgress > 0 {
		shadowAlpha := uint8(25 * c.hoverProgress)
		shadowColor := c.colorWithAlpha(c.getShadowColor(), shadowAlpha)

		// 绘制多个阴影层实现模糊效果
		for i := 0; i < 3; i++ {
			offset := float64(i) * 0.5
			dc.SetColor(c.colorWithAlpha(shadowColor, uint8(float64(shadowAlpha)*0.7)))
			dc.DrawRoundedRectangle(
				cardX+offset,
				cardY+offset,
				cardWidth,
				cardHeight,
				c.layout.BorderRadius,
			)
			dc.Fill()
		}
	}

	// 绘制卡片背景
	bgColor := c.getBgColor()
	dc.SetColor(bgColor)
	dc.DrawRoundedRectangle(cardX, cardY, cardWidth, cardHeight, c.layout.BorderRadius)
	dc.Fill()

	// ==================== 绘制边框 ====================

	// 边框颜色：选中时为选中颜色，hover时逐渐变为hover颜色
	var borderColor color.Color
	if c.checked {
		borderColor = c.getSelectedColor()
	} else {
		// 根据hover进度插值
		borderColor = c.interpolateColor(c.getBorderColor(), c.getHoverColor(), c.hoverProgress)
	}

	// 边框宽度：hover或选中时稍微变粗
	borderWidth := 2.0
	if c.checked || c.hoverProgress > 0 {
		borderWidth = 2.5
	}

	dc.SetLineWidth(borderWidth)
	dc.SetColor(borderColor)
	dc.DrawRoundedRectangle(cardX, cardY, cardWidth, cardHeight, c.layout.BorderRadius)
	dc.Stroke()

	// ==================== 绘制左上角选中标记（CSS :before 伪元素风格）===================

	checkMarkX := cardX + c.layout.Padding
	checkMarkY := cardY + c.layout.Padding
	checkMarkSize := c.layout.CheckSize
	centerX := checkMarkX + checkMarkSize/2
	centerY := checkMarkY + checkMarkSize/2

	// CSS效果：圆圈从scale(0)到scale(1)的缩放动画
	circleScale := 0.5 + 0.5*c.circleProgress // 从0.5缩放到1.0
	currentRadius := checkMarkSize / 2 * circleScale

	// 只有hover或选中时才绘制圆圈
	if c.circleProgress > 0 {
		// 圆圈边框颜色：选中时为选中颜色，hover时为hover颜色，否则为边框颜色
		var circleBorderColor color.Color
		if c.checked {
			circleBorderColor = c.getSelectedColor()
		} else {
			circleBorderColor = c.interpolateColor(c.getBorderColor(), c.getHoverColor(), c.hoverProgress)
		}

		// 圆圈填充：选中时为圆圈颜色，hover时为白色
		var circleFillColor color.Color = color.White
		if c.checked {
			circleFillColor = c.getCircleColor()
		}

		// 绘制圆圈填充
		dc.SetColor(circleFillColor)
		dc.DrawCircle(centerX, centerY, currentRadius)
		dc.Fill()

		// 绘制圆圈边框
		dc.SetLineWidth(2.0)
		dc.SetColor(circleBorderColor)
		dc.DrawCircle(centerX, centerY, currentRadius)
		dc.Stroke()

		// 如果选中，绘制对勾（使用CheckmarkColor）
		if c.checked && c.checkProgress > 0.3 {
			// 对勾尺寸
			checkmarkSize := checkMarkSize * 0.5 * c.checkProgress

			// 绘制对勾
			dc.SetColor(c.getCheckmarkColor())
			dc.SetLineWidth(2.0)

			// 对勾路径（从左上到右下的勾）
			x1, y1 := centerX-checkmarkSize*0.3, centerY
			x2, y2 := centerX-checkmarkSize*0.1, centerY+checkmarkSize*0.3
			x3, y3 := centerX+checkmarkSize*0.35, centerY-checkmarkSize*0.25

			dc.MoveTo(x1, y1)
			dc.LineTo(x2, y2)
			dc.LineTo(x3, y3)
			dc.Stroke()
		}
	}

	// ==================== 绘制SVG图标 ====================

	if c.svgImage != nil {
		// 图标居中在卡片内
		iconX := cardX + (cardWidth-c.layout.IconSize)/2
		iconY := cardY + (cardHeight-c.layout.IconSize)/2

		// 图标颜色过渡
		var iconColor color.Color = c.getIconColor()
		if c.checked {
			iconColor = c.getSelectedColor()
		} else if c.hoverProgress > 0 {
			iconColor = c.interpolateColor(c.getIconColor(), c.getHoverColor(), c.hoverProgress)
		}

		// 重新着色并绘制图标
		recoloredIcon := recolorImage(c.svgImage, iconColor)
		dc.DrawImage(recoloredIcon, int(iconX), int(iconY))
	}

	// ==================== 绘制文字标签 ====================

	if c.text != "" {
		// 设置字体
		if c.fontPath != "" && c.layout.LabelFontSize > 0 {
			if c.cachedFontFace == nil || c.cachedFontSize != c.layout.LabelFontSize {
				face, err := gg.LoadFontFace(c.fontPath, c.layout.LabelFontSize)
				if err == nil {
					c.cachedFontFace = face
					c.cachedFontSize = c.layout.LabelFontSize
				}
			}
			if c.cachedFontFace != nil {
				dc.SetFontFace(c.cachedFontFace)
			}
		}

		// 文字颜色：选中时为选中颜色，hover时逐渐变为hover颜色
		var textColor color.Color
		if c.checked {
			textColor = c.getSelectedColor()
		} else {
			textColor = c.interpolateColor(c.getLabelColor(), c.getHoverColor(), c.hoverProgress)
		}
		dc.SetColor(textColor)

		// 文字位置（卡片下方居中）
		textX := c.width / 2
		textY := cardY + cardHeight + c.layout.Spacing + c.layout.LabelFontSize/2

		// 居中对齐
		dc.DrawStringAnchored(c.text, textX, textY, 0.5, 0.5)
	}

	// 更新画布图像
	c.canvasImg.Image = dc.Image()
}

// ==================== 交互事件（修复版本）====================
func (c *MaterialCheckbox) Tapped(_ *fyne.PointEvent) {
	c.animMutex.Lock()
	c.checked = !c.checked

	if c.checked {
		c.targetCheck = 1.0
	} else {
		c.targetCheck = 0
	}
	c.animMutex.Unlock()

	// 通知回调
	if c.OnChanged != nil {
		c.OnChanged(c.checked)
	}

	// 触发重绘
	fyne.Do(func() {
		c.updateCanvas()
		c.canvasImg.Refresh()
	})
}

func (c *MaterialCheckbox) TappedSecondary(_ *fyne.PointEvent) {
	// 右键点击
}

func (c *MaterialCheckbox) MouseIn(_ *desktop.MouseEvent) {
	c.animMutex.Lock()
	c.hovered = true
	c.targetHover = 1.0
	c.animMutex.Unlock()

	// 立即触发一次更新
	fyne.Do(func() {
		c.updateCanvas()
		c.canvasImg.Refresh()
	})
}

func (c *MaterialCheckbox) MouseOut() {
	c.animMutex.Lock()
	c.hovered = false
	c.targetHover = 0
	c.animMutex.Unlock()

	// 立即触发一次更新
	fyne.Do(func() {
		c.updateCanvas()
		c.canvasImg.Refresh()
	})
}

func (c *MaterialCheckbox) MouseMoved(_ *desktop.MouseEvent) {
	// 可以添加悬停跟踪
}

func (c *MaterialCheckbox) FocusGained() {
	// 焦点效果（可选）
}

func (c *MaterialCheckbox) FocusLost() {
	// 焦点效果（可选）
}

// ==================== 公开方法 ====================
func (c *MaterialCheckbox) SetChecked(checked bool) {
	c.animMutex.Lock()
	if c.checked != checked {
		c.checked = checked

		if checked {
			c.targetCheck = 1.0
		} else {
			c.targetCheck = 0
		}

		if c.OnChanged != nil {
			c.OnChanged(c.checked)
		}
	}
	c.animMutex.Unlock()

	fyne.Do(func() {
		c.updateCanvas()
		c.canvasImg.Refresh()
	})
}

func (c *MaterialCheckbox) IsChecked() bool {
	return c.checked
}

func (c *MaterialCheckbox) SetText(text string) {
	c.text = text
	fyne.Do(func() {
		c.updateCanvas()
		c.canvasImg.Refresh()
	})
}

// ==================== 渲染器实现 ====================
func (c *MaterialCheckbox) CreateRenderer() fyne.WidgetRenderer {
	return &materialCheckboxRenderer{c}
}

type materialCheckboxRenderer struct {
	checkbox *MaterialCheckbox
}

func (r *materialCheckboxRenderer) Destroy() {
	r.checkbox.stopAnimation()
}

func (r *materialCheckboxRenderer) Layout(size fyne.Size) {
	r.checkbox.calculateLayout()
	r.checkbox.canvasImg.Resize(fyne.NewSize(
		float32(r.checkbox.width),
		float32(r.checkbox.height),
	))
	fyne.Do(func() {
		r.checkbox.updateCanvas()
	})
}

func (r *materialCheckboxRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		float32(r.checkbox.width),
		float32(r.checkbox.height),
	)
}

func (r *materialCheckboxRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.checkbox.canvasImg}
}

func (r *materialCheckboxRenderer) Refresh() {
	fyne.Do(func() {
		r.checkbox.updateCanvas()
		r.checkbox.canvasImg.Refresh()
	})
}

// ==================== 工具函数 ====================
// 线性插值
func (c *MaterialCheckbox) lerp(current, target, factor float64) float64 {
	return current + (target-current)*factor
}

// 绝对值
func (c *MaterialCheckbox) absFloat(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// 颜色插值
func (c *MaterialCheckbox) interpolateColor(c1, c2 color.Color, factor float64) color.Color {
	if factor <= 0 {
		return c1
	}
	if factor >= 1 {
		return c2
	}

	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return color.RGBA{
		R: uint8(float64(uint8(r1>>8))*(1-factor) + float64(uint8(r2>>8))*factor),
		G: uint8(float64(uint8(g1>>8))*(1-factor) + float64(uint8(g2>>8))*factor),
		B: uint8(float64(uint8(b1>>8))*(1-factor) + float64(uint8(b2>>8))*factor),
		A: uint8(float64(uint8(a1>>8))*(1-factor) + float64(uint8(a2>>8))*factor),
	}
}

// 设置颜色透明度
func (c *MaterialCheckbox) colorWithAlpha(c1 color.Color, alpha uint8) color.Color {
	r, g, b, _ := c1.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), alpha}
}

// ==================== 物理隔离容器 ====================
// WrapWithIsolationContainer 返回一个高度+2px的Stack包裹当前控件，实现物理隔离，防止canvas.Image刷新污染
func (c *MaterialCheckbox) WrapWithIsolationContainer() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.Transparent)
	bg.SetMinSize(fyne.NewSize(float32(c.width), float32(c.height+2)))
	return container.NewStack(
		bg,
		container.NewCenter(c),
	)
}

// ==================== 样式获取器（支持nil检查）====================
// 这些方法确保即使外部样式为nil，也能返回有效的默认颜色

func (c *MaterialCheckbox) getIconColor() color.Color {
	if c.Style.IconColor != nil {
		return c.Style.IconColor
	}
	return CheckboxColorText
}

func (c *MaterialCheckbox) getLabelColor() color.Color {
	if c.Style.LabelColor != nil {
		return c.Style.LabelColor
	}
	return CheckboxColorText
}

func (c *MaterialCheckbox) getBorderColor() color.Color {
	if c.Style.BorderColor != nil {
		return c.Style.BorderColor
	}
	return CheckboxColorDefault
}

func (c *MaterialCheckbox) getBgColor() color.Color {
	if c.Style.BgColor != nil {
		return c.Style.BgColor
	}
	return color.White
}

func (c *MaterialCheckbox) getShadowColor() color.Color {
	if c.Style.ShadowColor != nil {
		return c.Style.ShadowColor
	}
	return CheckboxColorShadow
}

func (c *MaterialCheckbox) getHoverColor() color.Color {
	if c.Style.HoverColor != nil {
		return c.Style.HoverColor
	}
	return CheckboxColorHover
}

func (c *MaterialCheckbox) getSelectedColor() color.Color {
	if c.Style.SelectedColor != nil {
		return c.Style.SelectedColor
	}
	return CheckboxColorSelected
}

func (c *MaterialCheckbox) getCircleColor() color.Color {
	if c.Style.CircleColor != nil {
		return c.Style.CircleColor
	}
	return CheckboxColorSelected
}

func (c *MaterialCheckbox) getCheckmarkColor() color.Color {
	if c.Style.CheckmarkColor != nil {
		return c.Style.CheckmarkColor
	}
	return color.White
}
