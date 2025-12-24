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
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// ==================== 样式配置 ====================

// StyleConfig 样式配置
type StyleConfig struct {
	Width            float64
	Height           float64
	CircleSize       float64
	IconSize         float64
	Spacing          float64
	TextOffsetY      float64
	LineHeight       float64
	LineWidth        float64
	IndicatorSize    float64
	IndicatorOffsetY float64
	ButtonAlpha      float64 // 按钮透明度 (0-1)
}

// ColorConfig 颜色配置
type ColorConfig struct {
	Background color.Color
	Normal     color.Color
	Active     color.Color
	Disabled   color.Color
	Line       color.Color
	Text       color.Color
	Icon       color.Color
	Button     color.Color // 按钮颜色
}

// DefaultStyle 默认样式
func DefaultStyle() StyleConfig {
	return StyleConfig{
		Width:            600,
		Height:           120,
		CircleSize:       35,
		IconSize:         28,
		Spacing:          120,
		TextOffsetY:      45,
		LineHeight:       2,
		LineWidth:        2,
		IndicatorSize:    8,
		IndicatorOffsetY: 60,
		ButtonAlpha:      0.1, // 10%透明度
	}
}

// DefaultColors 默认颜色
func DefaultColors() ColorConfig {
	return ColorConfig{
		Background: color.White,
		Normal:     color.RGBA{224, 224, 224, 255},
		Active:     color.RGBA{91, 192, 222, 255},
		Disabled:   color.RGBA{200, 200, 200, 255},
		Line:       color.RGBA{224, 224, 224, 255},
		Text:       color.Black,
		Icon:       color.RGBA{85, 85, 85, 255},
		Button:     color.White, // 白色按钮
	}
}

// ==================== SVG渲染 ====================

func loadSVG(filePath string, width, height int) (image.Image, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("SVG文件不存在: %s", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取SVG失败: %v", err)
	}

	icon, err := oksvg.ReadIconStream(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("解析SVG失败: %v", err)
	}

	icon.SetTarget(0, 0, float64(width), float64(height))
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	rast := rasterx.NewDasher(width, height, scanner)
	icon.Draw(rast, 1.0)

	return img, nil
}

// ==================== Tab项定义 ====================

type TabItem struct {
	ID         string
	Title      string
	IconPath   string
	Content    fyne.CanvasObject
	Enabled    bool
	cachedIcon image.Image
	iconError  error
}

// ==================== StepTabs 组件 ====================

type StepTabs struct {
	widget.BaseWidget

	items     []*TabItem
	current   int
	style     StyleConfig
	colors    ColorConfig
	OnChanged func(int, string)

	// 关键：用透明层处理点击
	mainContainer *fyne.Container
	bgImage       *canvas.Image
	clickOverlay  *transparentClickOverlay

	isDirty bool
}

// 新建StepTabs
func NewStepTabs(items []*TabItem) (*StepTabs, error) {
	if len(items) < 2 {
		return nil, fmt.Errorf("至少需要2个Tab项")
	}

	for i, item := range items {
		if item.IconPath == "" {
			return nil, fmt.Errorf("第%d项没有图标路径", i)
		}
		ext := strings.ToLower(filepath.Ext(item.IconPath))
		if ext != ".svg" {
			return nil, fmt.Errorf("第%d项不是SVG文件: %s", i, item.IconPath)
		}
	}

	t := &StepTabs{
		items:   items,
		current: 0,
		style:   DefaultStyle(),
		colors:  DefaultColors(),
		isDirty: true,
	}

	// 加载所有图标
	t.loadAllIcons()

	// 创建主容器
	t.mainContainer = container.NewWithoutLayout()

	// 创建背景图像（包含按钮）
	t.bgImage = canvas.NewImageFromResource(nil)
	t.bgImage.FillMode = canvas.ImageFillOriginal
	t.mainContainer.Add(t.bgImage)

	// 创建透明点击层
	t.clickOverlay = &transparentClickOverlay{tabs: t}
	t.clickOverlay.ExtendBaseWidget(t.clickOverlay)
	t.mainContainer.Add(t.clickOverlay)

	t.ExtendBaseWidget(t)
	return t, nil
}

// 加载所有图标
func (t *StepTabs) loadAllIcons() {
	for _, item := range t.items {
		img, err := loadSVG(item.IconPath, int(t.style.IconSize), int(t.style.IconSize))
		item.cachedIcon = img
		item.iconError = err
	}
}

// 获取按钮位置（供点击层使用）
func (t *StepTabs) getButtonPositions() []struct {
	Index  int
	Center fyne.Position
	Radius float32
} {
	style := t.style
	itemCount := len(t.items)

	// 取整计算
	circleSize := math.Round(style.CircleSize)
	spacing := math.Round(style.Spacing)

	totalWidth := float64(itemCount)*2*circleSize + float64(itemCount-1)*spacing
	startX := math.Round((style.Width - totalWidth) / 2)
	centerY := math.Round(style.Height / 2)

	var positions []struct {
		Index  int
		Center fyne.Position
		Radius float32
	}

	for i := range t.items {
		x := math.Round(startX + float64(i)*(2*circleSize+spacing) + circleSize)

		// 点击半径：按钮半径（比边框小2像素）
		clickRadius := circleSize - 2

		positions = append(positions, struct {
			Index  int
			Center fyne.Position
			Radius float32
		}{
			Index:  i,
			Center: fyne.NewPos(float32(x), float32(centerY)),
			Radius: float32(clickRadius),
		})
	}

	return positions
}

// 设置样式
func (t *StepTabs) SetStyle(style StyleConfig) {
	t.style = style
	t.isDirty = true
	t.Refresh()
}

// 设置颜色
func (t *StepTabs) SetColors(colors ColorConfig) {
	t.colors = colors
	t.isDirty = true
	t.Refresh()
}

// 选择Tab
func (t *StepTabs) Select(index int) error {
	if index < 0 || index >= len(t.items) {
		return fmt.Errorf("索引越界: %d", index)
	}

	if !t.items[index].Enabled {
		return fmt.Errorf("Tab项已禁用: %d", index)
	}

	t.current = index
	t.isDirty = true

	if t.OnChanged != nil {
		t.OnChanged(index, t.items[index].ID)
	}

	t.Refresh()
	return nil
}

// 获取当前索引
func (t *StepTabs) GetCurrentIndex() int {
	return t.current
}

// 绘制Tab栏（包含按钮）
func (t *StepTabs) drawTabBar() image.Image {
	style := t.style
	colors := t.colors
	items := t.items

	dc := gg.NewContext(int(style.Width), int(style.Height))
	dc.SetColor(colors.Background)
	dc.Clear()

	// ==================== 所有关键尺寸取整 ====================
	circleSize := math.Round(style.CircleSize)
	iconSize := math.Round(style.IconSize)
	spacing := math.Round(style.Spacing)
	indicatorSize := math.Round(style.IndicatorSize)
	textOffsetY := math.Round(style.TextOffsetY)
	indicatorOffsetY := math.Round(style.IndicatorOffsetY)

	itemCount := len(items)

	// 计算总宽度并取整
	totalWidth := float64(itemCount)*2*circleSize + float64(itemCount-1)*spacing
	startX := math.Round((style.Width - totalWidth) / 2)
	centerY := math.Round(style.Height / 2)

	// ==================== 绘制连接线 ====================
	lineY := math.Round(centerY)
	lineStartX := math.Round(startX + circleSize)
	lineEndX := math.Round(startX + totalWidth - circleSize)

	dc.SetColor(colors.Line)
	dc.SetLineWidth(math.Round(style.LineWidth))
	dc.DrawLine(lineStartX, lineY, lineEndX, lineY)
	dc.Stroke()

	// ==================== 绘制每个Tab项 ====================
	for i, item := range items {
		// 计算x坐标并取整
		x := math.Round(startX + float64(i)*(2*circleSize+spacing) + circleSize)

		// 确定颜色
		borderColor := colors.Normal
		textColor := colors.Text

		if !item.Enabled {
			borderColor = colors.Disabled
			textColor = colors.Disabled
		} else if i == t.current {
			borderColor = colors.Active
			textColor = colors.Active
		}

		// 1. 绘制圆形背景
		dc.SetColor(colors.Background)
		dc.DrawCircle(x, centerY, circleSize)
		dc.Fill()

		// 2. 绘制圆形边框（装饰，在外面）
		dc.SetColor(borderColor)
		dc.SetLineWidth(2)
		dc.DrawCircle(x, centerY, circleSize)
		dc.Stroke()

		// 3. 绘制按钮（在边框内部，无边框）
		if item.Enabled {
			// 按钮在边框内部，比边框小2像素
			buttonRadius := circleSize - 2

			// 半透明白色（无边框）
			r, g, b, _ := colors.Button.RGBA()
			alpha := uint8(float64(255) * style.ButtonAlpha)
			buttonColor := color.RGBA{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
				alpha,
			}

			dc.SetColor(buttonColor)
			dc.DrawCircle(x, centerY, buttonRadius)
			dc.Fill()
		}

		// 4. 绘制图标
		if item.cachedIcon != nil {
			iconX := math.Round(x - iconSize/2)
			iconY := math.Round(centerY - iconSize/2)
			dc.DrawImage(item.cachedIcon, int(iconX), int(iconY))
		} else if item.iconError != nil {
			dc.SetColor(color.RGBA{200, 200, 200, 255})
			rectX := math.Round(x - iconSize/2)
			rectY := math.Round(centerY - iconSize/2)
			dc.DrawRectangle(rectX, rectY, iconSize, iconSize)
			dc.Stroke()
			dc.DrawStringAnchored("?", x, centerY, 0.5, 0.5)
		}

		// 5. 绘制标题
		if item.Title != "" {
			dc.SetColor(textColor)
			if err := dc.LoadFontFace("", textOffsetY/3); err == nil {
				textY := math.Round(centerY + textOffsetY)
				dc.DrawStringAnchored(item.Title, x, textY, 0.5, 0.5)
			}
		}

		// 6. 绘制激活指示器
		if i == t.current && item.Enabled {
			indicatorY := math.Round(centerY + circleSize + indicatorOffsetY)

			dc.SetColor(borderColor)
			dc.MoveTo(x-indicatorSize, indicatorY)
			dc.LineTo(x+indicatorSize, indicatorY)
			dc.LineTo(x, indicatorY+indicatorSize)
			dc.ClosePath()
			dc.Fill()
		}
	}

	return dc.Image()
}

// ==================== 透明点击层 ====================

type transparentClickOverlay struct {
	widget.BaseWidget
	tabs *StepTabs
}

func (o *transparentClickOverlay) CreateRenderer() fyne.WidgetRenderer {
	return &transparentClickOverlayRenderer{}
}

type transparentClickOverlayRenderer struct{}

func (r *transparentClickOverlayRenderer) Layout(size fyne.Size) {}
func (r *transparentClickOverlayRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}
func (r *transparentClickOverlayRenderer) Refresh() {}
func (r *transparentClickOverlayRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{}
}
func (r *transparentClickOverlayRenderer) Destroy() {}

func (o *transparentClickOverlay) Tapped(e *fyne.PointEvent) {
	positions := o.tabs.getButtonPositions()

	for _, pos := range positions {
		if !o.tabs.items[pos.Index].Enabled {
			continue
		}

		// 关键：点击检测半径也缩小到90%
		clickRadius := float32(float64(pos.Radius) * 0.9)

		dx := float64(e.Position.X - pos.Center.X)
		dy := float64(e.Position.Y - pos.Center.Y)
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance <= float64(clickRadius) {
			o.tabs.Select(pos.Index)
			return
		}
	}
}

// ==================== 渲染器实现 ====================

type stepTabsRenderer struct {
	tabs    *StepTabs
	objects []fyne.CanvasObject
}

func (t *StepTabs) CreateRenderer() fyne.WidgetRenderer {
	return &stepTabsRenderer{
		tabs:    t,
		objects: []fyne.CanvasObject{t.mainContainer},
	}
}

func (r *stepTabsRenderer) Layout(size fyne.Size) {
	r.tabs.style.Width = float64(size.Width)
	r.tabs.style.Height = float64(size.Height)

	// 设置背景图像大小
	r.tabs.bgImage.Resize(size)
	r.tabs.bgImage.Move(fyne.NewPos(0, 0))

	// 设置点击层大小
	r.tabs.clickOverlay.Resize(size)
	r.tabs.clickOverlay.Move(fyne.NewPos(0, 0))
}

func (r *stepTabsRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(r.tabs.style.Width), float32(r.tabs.style.Height))
}

func (r *stepTabsRenderer) Refresh() {
	if !r.tabs.isDirty {
		return
	}

	img := r.tabs.drawTabBar()
	r.tabs.bgImage.Image = img
	r.tabs.bgImage.Refresh()

	r.tabs.isDirty = false
}

func (r *stepTabsRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *stepTabsRenderer) Destroy() {}

// 更新图标
func (t *StepTabs) UpdateItemIcon(index int, iconPath string) error {
	if index < 0 || index >= len(t.items) {
		return fmt.Errorf("索引越界: %d", index)
	}

	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		return fmt.Errorf("图标文件不存在: %s", iconPath)
	}

	ext := strings.ToLower(filepath.Ext(iconPath))
	if ext != ".svg" {
		return fmt.Errorf("不是SVG文件: %s", iconPath)
	}

	t.items[index].IconPath = iconPath
	img, err := loadSVG(iconPath, int(t.style.IconSize), int(t.style.IconSize))
	t.items[index].cachedIcon = img
	t.items[index].iconError = err

	t.isDirty = true
	t.Refresh()

	return nil
}

// ==================== 物理隔离容器 ====================
// WrapWithIsolationContainer 返回一个高度+2px的Stack包裹当前控件，实现物理隔离，防止canvas.Image刷新污染
func (t *StepTabs) WrapWithIsolationContainer() fyne.CanvasObject {
	// 使用透明矩形作为背景，高度增加2px以提供物理隔离
	bg := canvas.NewRectangle(color.Transparent)
	bg.SetMinSize(fyne.NewSize(float32(t.style.Width), float32(t.style.Height+2)))

	// 使用Stack容器进行包裹
	return container.NewStack(
		bg,                     // 透明背景层，提供物理高度
		container.NewCenter(t), // 居中放置StepTabs
	)
}
