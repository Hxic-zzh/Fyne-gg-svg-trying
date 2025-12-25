# 自定义 UI 控件库 API 速查

> 基于 Fyne 的纯 Go 实现，零 CGo 依赖，可直接 `go run` / `go build`。
> 暂时展示页，过年去完善
---
go.mod
```
module 2025-12-18-ggAndPng

go 1.25.4

require github.com/shirou/gopsutil/v3 v3.24.5

require (
	fyne.io/fyne/v2 v2.7.1 // indirect
	fyne.io/systray v1.11.1-0.20250603113521-ca66a66d8b58 // indirect
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/fredbi/uri v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fyne-io/gl-js v0.2.0 // indirect
	github.com/fyne-io/glfw-js v0.3.0 // indirect
	github.com/fyne-io/image v0.1.1 // indirect
	github.com/fyne-io/oksvg v0.2.0 // indirect
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240506104042-037f3cc74f2a // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-text/render v0.2.0 // indirect
	github.com/go-text/typesetting v0.2.1 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/hack-pad/go-indexeddb v0.3.2 // indirect
	github.com/hack-pad/safejs v0.1.0 // indirect
	github.com/jeandeaual/go-locale v0.0.0-20250612000132-0ef82f21eade // indirect
	github.com/jsummers/gobmp v0.0.0-20230614200233-a9de23ed2e25 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/rymdport/portal v0.4.2 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yuin/goldmark v1.7.8 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/image v0.34.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

```

## 1. 粒子按钮 ParticleButton
| API | 签名 | 说明 |
|---|---|---|
| 构造 | `NewParticleButton(onClick func()) *ParticleButton` | 绿色默认按钮 |
| 构造 | `NewParticleButtonWithStyle(onClick func(), text string, style ParticleButtonStyle) *ParticleButton` | 全样式一次性配置 |
| 设置文字 | `SetText(text string)` | 运行时改文字 |
| 设置尺寸 | `SetSize(w, h float32)` | 单位：像素 |
| 设置颜色 | `SetBaseColor(c color.RGBA)` | 自动推导渐变、阴影、粒子色 |
| 粒子开关 | `EnableParticle = true / false` | 运行时关闭特效 |
| 自动变色 | `AutoColorful = true` | 每次点击随机换色 |
| 字体 | `UseGGFont=true` + `GGFontType="chinese"/"english"` | 使用内置 TTF 渲染 |
| 回调 | `OnClick func()` | 点击事件 |
| 隔离容器 | `WrapWithIsolationContainer() fyne.CanvasObject` | 防刷新污染，推荐外层使用 |

示范：
```
import (
    "image/color"
    "2025-12-18-ggAndPng/tools"
    "fyne.io/fyne/v2"
)

// ...

style := tools.ParticleButtonStyle{
    BaseColor:     color.RGBA{R: 0, G: 180, B: 255, A: 255}, // 主色调
    CanvasBorder:  2,                                         // 边框宽度
    UseGGFont:     true,                                      // 启用gg字体
    GGFontType:    "english",                                 // 字体类型
    GGFontSize:    24,                                        // 字号
    GGFontColor:   color.RGBA{255, 255, 255, 255},            // 字体颜色
    GGFontOffsetX: 0,                                         // X偏移
    GGFontOffsetY: 0,                                         // Y偏移
    AutoColorful:  true,                                      // 自动变色
    CanvasOffsetX: 0,                                         // 画布X偏移
    CanvasOffsetY: 0,                                         // 画布Y偏移
}

btn := tools.NewParticleButtonWithStyle(
    func() { println("按钮被点击") },
    "自定义按钮",
    style,
)

// 添加到窗口
myWindow.SetContent(btn)

```


---

## 2. 边框按钮 BorderButton
| API | 签名 | 说明 |
|---|---|---|
| 构造 | `NewBorderButton(onToggle func(active bool), text string) *BorderButton` | 默认浅灰主题 |
| 构造 | `NewBorderButtonWithStyle(onToggle func(active bool), text string, style BorderButtonStyle) *BorderButton` | 全样式一次性配置 |
| 反色/透明 | `NewInverseBorderButton(...) / NewTransparentBorderButton(...)` | 预制两种风格 |
| 设置激活 | `SetActive(bool)` / `Toggle()` | 开关状态 |
| 查询状态 | `IsActive() bool` | 读取当前状态 |
| 设置文字 | `SetText(string)` | 运行时改文字 |
| 设置尺寸 | `SetSize(w, h float32)` | 单位：像素 |
| 轮廓缩放 | `SetContourScale(0.0~1.0)` | 外框相对尺寸 |
| 隔离容器 | `WrapWithIsolationContainer() fyne.CanvasObject` | 同上 |


示范：
```
import (
    "image/color"
    "2025-12-18-ggAndPng/tools"
    "fyne.io/fyne/v2"
)

// ...

style := tools.BorderButtonStyle{
    DefaultColor:   color.RGBA{240, 240, 240, 255},
    DefaultContour: color.RGBA{180, 180, 180, 255},
    DefaultText:    color.RGBA{51, 51, 51, 255},
    HoverColor:     color.RGBA{220, 220, 220, 255},
    PressedColor:   color.RGBA{200, 200, 200, 255},
    PressedContour: color.RGBA{150, 150, 150, 255},
    PressedText:    color.RGBA{30, 30, 30, 255},
    ActiveColor:    color.RGBA{200, 230, 255, 255},
    ActiveContour:  color.RGBA{100, 180, 255, 255},
    ActiveText:     color.RGBA{0, 100, 200, 255},
    BorderRadius:   12,
    GGFontType:     "chinese",
    GGFontSize:     18,
    GGFontColor:    color.RGBA{0, 100, 200, 255},
    UseGGFont:      true,
}

btn := tools.NewBorderButtonWithStyle(
    func(active bool) { println("按钮状态:", active) },
    "自定义边框按钮",
    style,
)

// 推荐用隔离容器包裹
content := btn.WrapWithIsolationContainer()

myWindow.SetContent(content)
```

---

## 3. Material 风格输入框 MaterialEntry
| API | 签名 | 说明 |
|---|---|---|
| 构造 | `NewMaterialEntry(placeholder string, width, height float64) *MaterialEntry` | 可省略尺寸 |
| 设置样式 | `SetStyle(style MaterialEntryStyle)` | 一次性配置颜色、字号、圆角等 |
| 设置文字 | `SetText(string)` / `GetText() string` | 读写 |
| 清空 | `Clear()` | 置空并触发动画 |
| 占位符 | `SetPlaceholder(string)` | 浮动标签文字 |
| 字体 | `SetFontPath("ttf/xxx.ttf")` | 支持中英文 |
| 背景 | `SetCustomBackground(color.Color)` | 覆盖主题 |
| 毛玻璃 | `SetGlassEffect(true, blurRadius float64)` | 半透+模糊 |
| 完全透明 | `SetTransparent(true)` | 背景与边框全透明 |
| 圆角 | `SetCornerRadius(radius float64)` | 实时改圆角 |
| 回调 | `OnChanged func(string)` / `OnSubmitted func(string)` | 内容变化/回车 |
| 焦点 | `FocusGained()` / `Focused() bool` | 程序控制焦点 |
| 隔离容器 | `WrapWithIsolationContainer() fyne.CanvasObject` | 同上 |

示范：
```
import (
    "image/color"
    "2025-12-18-ggAndPng/tools"
    "fyne.io/fyne/v2"
)

// ...

style := tools.MaterialEntryStyle{
    Width:           320,
    Height:          48,
    FontSize:        20,
    LabelColor:      color.RGBA{82, 100, 174, 255},   // 标签颜色
    TextColor:       color.RGBA{33, 33, 33, 255},     // 输入文字颜色
    BorderColor:     color.RGBA{117, 117, 117, 255},  // 边框颜色
    BgColor:         color.RGBA{255, 255, 255, 255},  // 背景色
    Radius:          8,                               // 圆角
    UnderlineColor:  color.RGBA{82, 100, 174, 255},   // 下划线颜色
    UnderlineHeight: 4,                               // 下划线粗细
}

entry := tools.NewMaterialEntry("请输入内容")
entry.SetStyle(style)
entry.SetFontPath("ttf/chinese.ttf") // 如需自定义字体

// 推荐用隔离容器包裹
content := entry.WrapWithIsolationContainer()

myWindow.SetContent(content)
```

---

## 4. Material 风格复选框 MaterialCheckbox
| API | 签名 | 说明 |
|---|---|---|
| 构造 | `NewMaterialCheckbox(text string, checked bool, tileWidth, tileHeight float64) *MaterialCheckbox` | 默认 112×112 |
| 设置样式 | `SetStyle(style MaterialCheckboxStyle)` | 颜色、圆角、图标、动画色 |
| 选中状态 | `SetChecked(bool)` / `IsChecked() bool` | 读写 |
| 文字 | `SetText(string)` | 运行时改文字 |
| 图标 | `SetIconPath("svg/xxx.svg")` | 自动着色 |
| 字体 | `SetFontPath("ttf/xxx.ttf")` | 支持中英文 |
| 回调 | `OnChanged func(bool)` | 状态变化 |
| 隔离容器 | `WrapWithIsolationContainer() fyne.CanvasObject` | 同上 |

示范：
```
import (
    "image/color"
    "2025-12-18-ggAndPng/tools"
    "fyne.io/fyne/v2"
)

// ...

style := tools.MaterialCheckboxStyle{
    TileWidth:     120,
    TileHeight:    120,
    IconColor:     color.RGBA{34, 96, 255, 255},   // 图标主色
    LabelColor:    color.RGBA{33, 33, 33, 255},    // 文字颜色
    BorderColor:   color.RGBA{181, 191, 217, 255}, // 边框色
    BgColor:       color.White,                    // 背景色
    ShadowColor:   color.RGBA{0, 0, 0, 25},        // 阴影色
    CornerRadius:  10,                             // 圆角
    IconSize:      48,                             // 图标尺寸
    FontSize:      18,                             // 字体大小
    IconPath:      "svg/1.svg",                    // SVG图标路径（可选）
    HoverColor:    color.RGBA{34, 96, 255, 255},   // 悬停色
    SelectedColor: color.RGBA{34, 96, 255, 255},   // 选中色
    CircleColor:   color.RGBA{34, 96, 255, 255},   // 圆圈色
    CheckmarkColor: color.White,                   // 对勾色
}

checkbox := tools.NewMaterialCheckbox("自定义复选框", false)
checkbox.SetStyle(style)
checkbox.SetFontPath("ttf/chinese.ttf") // 如需自定义字体
checkbox.OnChanged = func(checked bool) {
    println("复选框状态:", checked)
}

// 推荐用隔离容器包裹
content := checkbox.WrapWithIsolationContainer()

myWindow.SetContent(content)
```

---

## 5. 动画开关 ToggleSwitch
| API | 签名 | 说明 |
|---|---|---|
| 构造 | `NewToggleSwitch(checked bool) *ToggleSwitch` | 默认 74×36 像素 |
| 设置效果 | `SetEffect(EffectSlide / EffectTwoBallSwap / ...)` | 6 种内置动画 |
| 设置尺寸 | `SetSize(width, height float64)` | 任意尺寸 |
| 设置配置 | `SetConfig(config SwitchConfig)` | 一次性配置颜色、文字、字体 |
| 链式调用 | `SetYesLabel("开").SetNoLabel("关").SetYesColor(green).SetNoColor(red)...` | 全字段支持 |
| 布尔映射 | `SetYesValue(true).SetNoValue(false)` | 自定义左右布尔值 |
| 读写 | `SetChecked(bool) / GetChecked() bool` | 状态 |
| 切换 | `Toggle()` | 反向状态 |
| 回调 | `OnChanged func(value bool)` | value 为映射后的布尔值 |
| 隔离容器 | `WrapWithIsolationContainer() fyne.CanvasObject` | 同上 |

---

示范：
```
import (
    "image/color"
    "2025-12-18-ggAndPng/tools"
    "fyne.io/fyne/v2"
)

// ...

// 自定义配置
config := tools.SwitchConfig{
    YesLabel:      "开",
    NoLabel:       "关",
    FontPath:      "ttf/toggle_switch.ttf",
    YesColor:      color.RGBA{3, 169, 244, 255},    // 开启色
    NoColor:       color.RGBA{244, 67, 54, 255},    // 关闭色
    YesBgColor:    color.RGBA{235, 247, 252, 255},  // 开启背景
    NoBgColor:     color.RGBA{252, 235, 235, 255},  // 关闭背景
    TextColor:     color.White,
    TextDarkColor: color.RGBA{78, 78, 78, 255},
    YesValue:      true,
    NoValue:       false,
}

sw := tools.NewToggleSwitch(false).
    SetConfig(config).
    SetEffect(tools.EffectTwoBallSwap). // 切换不同动画效果
    SetSize(120, 48).
    SetOnChanged(func(val bool) {
        println("开关状态:", val)
    })

// 推荐用隔离容器包裹
content := sw.WrapWithIsolationContainer()

myWindow.SetContent(content)
```

## 6. 步骤标签页 StepTabs
| API | 签名 | 说明 |
|---|---|---|
| 构造 | `NewStepTabs(items []*TabItem) (*StepTabs, error)` | 至少 2 项 |
| 选择 | `Select(index int) error` | 切换步骤 |
| 当前索引 | `GetCurrentIndex() int` | 读取 |
| 样式 | `SetStyle(style StyleConfig)` | 尺寸、圆角、间距 |
| 颜色 | `SetColors(colors ColorConfig)` | 背景、线条、激活色 |
| 更新图标 | `UpdateItemIcon(index int, svgPath string) error` | 运行时换图标 |
| 回调 | `OnChanged func(index int, id string)` | 切换事件 |
| 隔离容器 | `WrapWithIsolationContainer() fyne.CanvasObject` | 同上 |

示范：
```
import (
    "2025-12-18-ggAndPng/tools"
    "fyne.io/fyne/v2"
    "image/color"
)

// ...

// 定义每个步骤的Tab项
tabs := []*tools.TabItem{
    {
        ID:       "step1",
        Title:    "上传",
        IconPath: "svg/1.svg",
        Content:  nil, // 可选：每个Tab对应的内容
        Enabled:  true,
    },
    {
        ID:       "step2",
        Title:    "处理",
        IconPath: "svg/2.svg",
        Content:  nil,
        Enabled:  true,
    },
    {
        ID:       "step3",
        Title:    "完成",
        IconPath: "svg/3.svg",
        Content:  nil,
        Enabled:  true,
    },
}

// 创建StepTabs控件
stepTabs, err := tools.NewStepTabs(tabs)
if err != nil {
    panic(err)
}

// 可选：自定义样式
style := tools.StyleConfig{
    Width:         600,
    Height:        120,
    CircleSize:    32,
    IconSize:      24,
    Spacing:       100,
    TextOffsetY:   40,
    LineHeight:    2,
    LineWidth:     2,
    IndicatorSize: 8,
    IndicatorOffsetY: 60,
    ButtonAlpha:   0.12,
}
stepTabs.SetStyle(style)

// 可选：自定义颜色
colors := tools.ColorConfig{
    Background: color.White,
    Normal:     color.RGBA{224, 224, 224, 255},
    Active:     color.RGBA{91, 192, 222, 255},
    Disabled:   color.RGBA{200, 200, 200, 255},
    Line:       color.RGBA{224, 224, 224, 255},
    Text:       color.Black,
    Icon:       color.RGBA{85, 85, 85, 255},
    Button:     color.White,
}
stepTabs.SetColors(colors)

// 可选：监听切换事件
stepTabs.OnChanged = func(idx int, id string) {
    println("当前步骤:", idx, "ID:", id)
}

// 推荐用隔离容器包裹
content := stepTabs.WrapWithIsolationContainer()

myWindow.SetContent(content)
```

---

## 7. 小猪进度条 PigProgressBar （非gg绘制，图片叠加）
| API | 签名 | 说明 |
|---|---|---|
| 构造 | `NewPigProgressBar() *PigProgressBar` | 自动加载帧动画 |
| 设置进度 | `SetProgress(float64)` | 0.0~1.0，平滑移动 |
| 手动位置 | `SetPosition(x, y float32)` | 禁用自动居中 |
| 自动居中 | `SetAutoCenter(true)` | 恢复自动居中 |

---

## 8. 性能基准工具（benchmark 子包）
| API | 签名 | 说明 |
|---|---|---|
| 监控器 | `monitor := benchmark.NewMonitor("测试名")` | 创建实例 |
| 开始 | `monitor.Start()` | 启动后台采样 |
| 记录场景 | `monitor.StartRecording("组件名", "custom|native", "场景")` | 标记测试区间 |
| 结束场景 | `monitor.StopRecording()` | 结束标记 |
| 加帧 | `monitor.AddFrame()` | 在渲染循环调用 |
| 取结果 | `metrics := monitor.GetComponentMetrics("组件名", "custom")` | 拿到切片 |
| 对比 | `comparison := benchmark.CompareComponents(customMetrics, nativeMetrics)` | 科学对比 |
| 导出 CSV | `benchmark.NewCSVExporter("./result").ExportMetrics(metrics, summary)` | 一键生成报告 |

---

## 9. 主题透明辅助
| API | 签名 | 说明 |
|---|---|---|
| 透明主题 | `tools.NewInputTransparentTheme()` | 仅输入框背景/边框透明，其余保持默认主题 |

---

## 10. 物理隔离约定
所有自定义控件均提供  
`WrapWithIsolationContainer() fyne.CanvasObject`  
返回高度+2 px 的 `container.Stack`，解决 `canvas.Image` 刷新时污染兄弟节点的问题；  
**在main中调用的时候先加一层stack再放置**。

---
