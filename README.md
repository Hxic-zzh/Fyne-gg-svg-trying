# 自定义 UI 控件库 API 速查

> 基于 Fyne 的纯 Go 实现，零 CGo 依赖，可直接 `go run` / `go build`。
> 暂时展示页，过年去完善
---

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

---

## 7. 小猪进度条 PigProgressBar
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
**推荐将控件放入此容器后再加入布局**。

---
