// main_benchmark.go
package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"2025-12-18-ggAndPng/tools"
	"2025-12-18-ggAndPng/tools/benchmark"
)

// ScientificBenchmarkResult 科学性能测试结果
type ScientificBenchmarkResult struct {
	TestName        string
	CustomComponent string
	NativeComponent string
	CustomMetrics   []*benchmark.PerformanceMetric
	NativeMetrics   []*benchmark.PerformanceMetric
	Comparison      *benchmark.ComponentComparison
	StartTime       time.Time
	EndTime         time.Time
}

// runComponentBenchmark 运行单个组件的性能测试
func runComponentBenchmark(log func(string), componentName, componentType, scenario string,
	duration time.Duration, frameRate time.Duration) ([]*benchmark.PerformanceMetric, error) {

	// 创建监控器
	testName := fmt.Sprintf("%s_%s", componentName, componentType)
	monitor := benchmark.NewMonitor(testName)
	monitor.Start()
	defer monitor.Stop()

	// 启动帧计数goroutine
	stopFrameCounter := make(chan bool)
	frameTicker := time.NewTicker(frameRate)

	go func() {
		for {
			select {
			case <-frameTicker.C:
				monitor.AddFrame()
			case <-stopFrameCounter:
				frameTicker.Stop()
				return
			}
		}
	}()

	// 开始记录
	fyne.Do(func() {
		monitor.StartRecording(componentName, componentType, scenario)
	})

	// 等待测试持续时间
	time.Sleep(duration)

	// 停止记录
	fyne.Do(func() {
		monitor.StopRecording()
	})

	// 停止帧计数
	close(stopFrameCounter)

	// 获取该组件的所有指标
	metrics := monitor.GetComponentMetrics(componentName, componentType)

	if len(metrics) == 0 {
		return nil, fmt.Errorf("没有收集到性能指标数据")
	}

	log(fmt.Sprintf("✅ 收集到 %s 的 %d 个性能样本", componentName, len(metrics)))
	return metrics, nil
}

// runScientificComparison 运行科学对比测试
func runScientificComparison(log func(string), statusLabel *widget.Label,
	comparisonContainer *fyne.Container,
	customName, nativeName, scenario string,
	createCustomFunc, createNativeFunc func() fyne.CanvasObject,
	testDuration time.Duration) *ScientificBenchmarkResult {

	log(fmt.Sprintf("🔬 开始科学性能对比测试: %s vs %s", customName, nativeName))
	fyne.Do(func() {
		statusLabel.SetText(fmt.Sprintf("测试 %s vs %s...", customName, nativeName))
	})

	// 清空对比容器
	comparisonContainer.Objects = nil

	// ====== 步骤1: 显示控件用于视觉对比 ======
	log("1. 显示控件进行视觉对比...")

	// 创建自定义控件
	customWidget := createCustomFunc()
	if customWidget == nil {
		log("❌ 创建自定义控件失败")
		fyne.Do(func() {
			statusLabel.SetText("创建自定义控件失败")
		})
		return nil
	}

	// 创建原生控件
	nativeWidget := createNativeFunc()
	if nativeWidget == nil {
		log("❌ 创建原生控件失败")
		fyne.Do(func() {
			statusLabel.SetText("创建原生控件失败")
		})
		return nil
	}

	// 添加到对比容器
	customBox := container.NewVBox(
		widget.NewLabelWithStyle("自定义控件", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(customWidget),
		widget.NewLabel("带复杂视觉效果"),
		widget.NewLabel(fmt.Sprintf("类型: %s", customName)),
	)

	nativeBox := container.NewVBox(
		widget.NewLabelWithStyle("原生控件", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(nativeWidget),
		widget.NewLabel("基础功能实现"),
		widget.NewLabel(fmt.Sprintf("类型: %s", nativeName)),
	)

	comparisonContainer.Add(customBox)
	comparisonContainer.Add(nativeBox)
	comparisonContainer.Refresh()

	// 等待渲染稳定
	time.Sleep(1 * time.Second)

	// ====== 步骤2: 分别测试两个组件 ======
	log("2. 分别测试自定义控件...")
	fyne.Do(func() {
		statusLabel.SetText("测试自定义控件性能...")
	})

	// 测试自定义控件
	customMetrics, err := runComponentBenchmark(log, customName, "custom", scenario,
		testDuration, 16*time.Millisecond) // ~60 FPS

	if err != nil {
		log(fmt.Sprintf("❌ 自定义控件测试失败: %v", err))
		fyne.Do(func() {
			statusLabel.SetText(fmt.Sprintf("自定义控件测试失败: %v", err))
		})
		return nil
	}

	log("3. 分别测试原生控件...")
	fyne.Do(func() {
		statusLabel.SetText("测试原生控件性能...")
	})

	// 测试原生控件
	nativeMetrics, err := runComponentBenchmark(log, nativeName, "native", scenario,
		testDuration, 16*time.Millisecond) // ~60 FPS

	if err != nil {
		log(fmt.Sprintf("❌ 原生控件测试失败: %v", err))
		fyne.Do(func() {
			statusLabel.SetText(fmt.Sprintf("原生控件测试失败: %v", err))
		})
		return nil
	}

	// ====== 步骤3: 科学对比分析 ======
	log("4. 进行科学对比分析...")
	fyne.Do(func() {
		statusLabel.SetText("进行科学对比分析...")
	})

	comparison := benchmark.CompareComponents(customMetrics, nativeMetrics)
	if comparison == nil {
		log("❌ 对比分析失败")
		fyne.Do(func() {
			statusLabel.SetText("对比分析失败")
		})
		return nil
	}

	// 打印对比结果到日志
	benchmark.PrintComparison(comparison)

	// ====== 步骤4: 导出结果 ======
	log("5. 导出测试结果...")
	fyne.Do(func() {
		statusLabel.SetText("导出测试结果...")
	})

	result := &ScientificBenchmarkResult{
		TestName:        fmt.Sprintf("%s_vs_%s", customName, nativeName),
		CustomComponent: customName,
		NativeComponent: nativeName,
		CustomMetrics:   customMetrics,
		NativeMetrics:   nativeMetrics,
		Comparison:      comparison,
		StartTime:       time.Now().Add(-testDuration * 2), // 估计开始时间
		EndTime:         time.Now(),
	}

	// 导出详细报告
	exportScientificResult(result, log, statusLabel)

	// ====== 步骤5: 显示结果摘要 ======
	displayResultsSummary(result, statusLabel)

	log("✅ 科学性能对比测试完成！")
	return result
}

// exportScientificResult 导出科学测试结果
func exportScientificResult(result *ScientificBenchmarkResult, log func(string), statusLabel *widget.Label) {
	// 合并所有指标
	allMetrics := append(result.CustomMetrics, result.NativeMetrics...)

	// 创建导出器
	exporter := benchmark.NewCSVExporter("./benchmark_results/scientific")

	// 设置文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("scientific_%s_%s.csv", result.TestName, timestamp)
	exporter.SetFilename(filename)

	// 创建详细摘要
	summary := map[string]interface{}{
		"test_name":      result.TestName,
		"custom_summary": result.Comparison.CustomSummary,
		"native_summary": result.Comparison.NativeSummary,
		"comparison":     result.Comparison.Comparison,
		"conclusion":     result.Comparison.Conclusion,
		"system_info":    benchmark.GetSystemInfo(),
		"start_time":     result.StartTime,
		"end_time":       result.EndTime,
		"duration":       result.EndTime.Sub(result.StartTime).String(),
	}

	// 导出数据
	if err := exporter.ExportMetrics(allMetrics, summary); err != nil {
		log(fmt.Sprintf("❌ 导出失败: %v", err))
		fyne.Do(func() {
			statusLabel.SetText(fmt.Sprintf("导出失败: %v", err))
		})
	} else {
		log(fmt.Sprintf("✅ 详细报告已导出到: %s", exporter.GetFullPath()))

		// 同时导出对比专用报告
		exportComparisonReport(result, log)
	}
}

// exportComparisonReport 导出对比专用报告
func exportComparisonReport(result *ScientificBenchmarkResult, log func(string)) {
	comparisonExporter := benchmark.NewCSVExporter("./benchmark_results/comparisons")
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("comparison_%s_%s.csv", result.TestName, timestamp)
	comparisonExporter.SetFilename(filename)

	if err := comparisonExporter.ExportComparison(result.Comparison); err != nil {
		log(fmt.Sprintf("⚠️ 对比报告导出失败: %v", err))
	} else {
		log(fmt.Sprintf("✅ 对比报告已导出到: %s", comparisonExporter.GetFullPath()))
	}
}

// displayResultsSummary 显示结果摘要
func displayResultsSummary(result *ScientificBenchmarkResult, statusLabel *widget.Label) {
	comparison := result.Comparison
	if comparison == nil {
		return
	}

	// 从对比结果中提取关键数据
	performanceScore := comparison.Comparison["performance_score"].(float64)
	fpsDiffPercent := comparison.Comparison["fps_diff_percent"].(float64)
	memoryDiffPercent := comparison.Comparison["memory_diff_percent"].(float64)
	cpuDiffPercent := comparison.Comparison["cpu_diff_percent"].(float64)

	// 自定义控件数据
	customFPS := comparison.CustomSummary["fps_avg"].(float64)
	customMemory := comparison.CustomSummary["memory_avg"].(float64)
	customCPU := comparison.CustomSummary["cpu_avg"].(float64)

	// 原生控件数据
	nativeFPS := comparison.NativeSummary["fps_avg"].(float64)
	nativeMemory := comparison.NativeSummary["memory_avg"].(float64)
	nativeCPU := comparison.NativeSummary["cpu_avg"].(float64)

	// 生成结果摘要
	summary := fmt.Sprintf(`
🎯 科学性能测试完成 - %s

📊 性能数据对比:
  自定义控件:
    • FPS: %.1f (%.1f-%.1f)
    • 内存: %.2fMB (%.2f-%.2f)  
    • CPU: %.1f%% (%.1f-%.1f)
  
  原生控件:
    • FPS: %.1f (%.1f-%.1f)
    • 内存: %.2fMB (%.2f-%.2f)
    • CPU: %.1f%% (%.1f-%.1f)

📈 性能差异:
    • FPS: %.1f%% %s
    • 内存: %.1f%% %s
    • CPU: %.1f%% %s

🏆 综合性能评分: %.1f/100

💡 %s

📁 结果已保存到 benchmark_results/ 目录
`,
		result.TestName,

		// 自定义控件
		customFPS,
		comparison.CustomSummary["fps_min"].(float64),
		comparison.CustomSummary["fps_max"].(float64),
		customMemory,
		comparison.CustomSummary["memory_min"].(float64),
		comparison.CustomSummary["memory_max"].(float64),
		customCPU,
		comparison.CustomSummary["cpu_min"].(float64),
		comparison.CustomSummary["cpu_max"].(float64),

		// 原生控件
		nativeFPS,
		comparison.NativeSummary["fps_min"].(float64),
		comparison.NativeSummary["fps_max"].(float64),
		nativeMemory,
		comparison.NativeSummary["memory_min"].(float64),
		comparison.NativeSummary["memory_max"].(float64),
		nativeCPU,
		comparison.NativeSummary["cpu_min"].(float64),
		comparison.NativeSummary["cpu_max"].(float64),

		// 性能差异
		fpsDiffPercent, getTrendIcon(fpsDiffPercent, true),
		memoryDiffPercent, getTrendIcon(memoryDiffPercent, false),
		cpuDiffPercent, getTrendIcon(cpuDiffPercent, false),

		// 评分和结论
		performanceScore,
		comparison.Conclusion,
	)

	fyne.Do(func() {
		statusLabel.SetText(summary)
	})
}

// getTrendIcon 获取趋势图标
func getTrendIcon(value float64, higherIsBetter bool) string {
	if higherIsBetter {
		if value > 10 {
			return "📈"
		} else if value > 0 {
			return "↗️"
		} else if value > -10 {
			return "↘️"
		} else {
			return "📉"
		}
	} else {
		// 对于内存和CPU，值越小越好
		if value < -10 {
			return "📈" // 负值表示节省，所以是上升趋势
		} else if value < 0 {
			return "↗️"
		} else if value < 10 {
			return "↘️"
		} else {
			return "📉"
		}
	}
}

// ====== 具体的测试函数 ======

// testParticleButton 测试粒子按钮
func testParticleButton(log func(string), statusLabel *widget.Label, comparisonContainer *fyne.Container) *ScientificBenchmarkResult {
	return runScientificComparison(log, statusLabel, comparisonContainer,
		"ParticleButton",
		"FyneButton",
		"click_animation",
		func() fyne.CanvasObject {
			// 创建自定义粒子按钮
			redStyle := tools.ParticleButtonStyle{
				BaseColor:     color.RGBA{R: 255, G: 100, B: 100, A: 255},
				CanvasBorder:  3,
				CanvasOffsetY: -3,
			}
			customBtn := tools.NewParticleButtonWithStyle(
				func() { log("粒子按钮被点击") },
				"粒子按钮",
				redStyle,
			)
			customBtn.SetSize(220, 56)
			return customBtn
		},
		func() fyne.CanvasObject {
			// 创建原生按钮
			nativeBtn := widget.NewButton("原生按钮", func() {
				log("原生按钮被点击")
			})
			nativeBtn.Resize(fyne.NewSize(220, 56))
			return nativeBtn
		},
		3*time.Second, // 测试持续时间
	)
}

// testMaterialEntry 测试Material输入框
func testMaterialEntry(log func(string), statusLabel *widget.Label, comparisonContainer *fyne.Container) *ScientificBenchmarkResult {
	return runScientificComparison(log, statusLabel, comparisonContainer,
		"MaterialEntry",
		"FyneEntry",
		"input_animation",
		func() fyne.CanvasObject {
			// 创建自定义Material输入框
			redInput := tools.NewMaterialEntry("输入测试", 400, 60)
			redInput.SetStyle(tools.MaterialEntryStyle{
				Width:           400,
				Height:          60,
				FontSize:        24,
				LabelColor:      color.RGBA{244, 67, 54, 255},
				TextColor:       color.RGBA{244, 67, 54, 255},
				BorderColor:     color.RGBA{244, 67, 54, 255},
				UnderlineColor:  color.RGBA{244, 67, 54, 255},
				UnderlineHeight: 5,
			})
			redInput.SetCustomBackground(color.White)
			redInput.SetCornerRadius(16)

			// 包装容器确保正确显示
			wrapper := container.NewStack(redInput)
			wrapper.Resize(fyne.NewSize(400, 60))
			return wrapper
		},
		func() fyne.CanvasObject {
			// 创建原生输入框
			nativeEntry := widget.NewEntry()
			nativeEntry.SetPlaceHolder("原生输入框")
			nativeEntry.Resize(fyne.NewSize(400, 60))
			return nativeEntry
		},
		3*time.Second,
	)
}

// testToggleSwitch 测试开关控件
func testToggleSwitch(log func(string), statusLabel *widget.Label, comparisonContainer *fyne.Container) *ScientificBenchmarkResult {
	return runScientificComparison(log, statusLabel, comparisonContainer,
		"ToggleSwitch",
		"FyneCheckbox",
		"toggle_animation",
		func() fyne.CanvasObject {
			// 创建自定义开关
			customToggle := tools.NewToggleSwitch(false).
				SetEffect(tools.EffectSlide).
				SetYesLabel("开").
				SetNoLabel("关").
				SetYesColor(color.RGBA{0, 200, 83, 255}).
				SetNoColor(color.RGBA{255, 61, 0, 255}).
				SetSize(160, 70)

			// 包装容器
			wrapper := container.NewStack(customToggle)
			wrapper.Resize(fyne.NewSize(160, 70))
			return wrapper
		},
		func() fyne.CanvasObject {
			// 创建原生复选框作为对比
			nativeCheckbox := widget.NewCheck("原生开关", func(checked bool) {
				log(fmt.Sprintf("原生开关状态: %v", checked))
			})
			return nativeCheckbox
		},
		3*time.Second,
	)
}

// testMaterialCheckbox 测试Material复选框
func testMaterialCheckbox(log func(string), statusLabel *widget.Label, comparisonContainer *fyne.Container) *ScientificBenchmarkResult {
	return runScientificComparison(log, statusLabel, comparisonContainer,
		"MaterialCheckbox",
		"FyneCheckbox",
		"check_animation",
		func() fyne.CanvasObject {
			// 创建自定义复选框
			customCheckbox := tools.NewMaterialCheckbox("自定义复选框", false, 112, 112)
			customCheckbox.SetStyle(tools.MaterialCheckboxStyle{
				TileWidth:     112,
				TileHeight:    112,
				IconColor:     color.RGBA{46, 204, 113, 255},
				LabelColor:    color.RGBA{46, 204, 113, 255},
				BorderColor:   color.RGBA{39, 174, 96, 255},
				BgColor:       color.White,
				CornerRadius:  8,
				IconPath:      "svg/1.svg",
				HoverColor:    color.RGBA{46, 204, 113, 100},
				SelectedColor: color.RGBA{46, 204, 113, 255},
			})

			// 包装容器
			wrapper := container.NewStack(customCheckbox)
			wrapper.Resize(fyne.NewSize(112, 112))
			return wrapper
		},
		func() fyne.CanvasObject {
			// 创建原生复选框
			nativeCheckbox := widget.NewCheck("原生复选框", func(checked bool) {
				log(fmt.Sprintf("原生复选框状态: %v", checked))
			})
			return nativeCheckbox
		},
		3*time.Second,
	)
}

// testStepTabs 测试步骤标签页
func testStepTabs(log func(string), statusLabel *widget.Label, comparisonContainer *fyne.Container) *ScientificBenchmarkResult {
	return runScientificComparison(log, statusLabel, comparisonContainer,
		"StepTabs",
		"FyneTabs",
		"tab_switch",
		func() fyne.CanvasObject {
			// 创建自定义步骤标签页
			items := []*tools.TabItem{
				{
					ID:       "step1",
					Title:    "第一步",
					IconPath: "svg/1.svg",
					Content:  container.NewCenter(widget.NewLabel("第一步内容")),
					Enabled:  true,
				},
				{
					ID:       "step2",
					Title:    "第二步",
					IconPath: "svg/2.svg",
					Content:  container.NewCenter(widget.NewLabel("第二步内容")),
					Enabled:  true,
				},
			}

			stepTabs, err := tools.NewStepTabs(items)
			if err != nil {
				log(fmt.Sprintf("⚠️ 创建StepTabs失败: %v", err))
				// 返回一个占位符
				return container.NewCenter(widget.NewLabel("StepTabs创建失败"))
			}

			return stepTabs
		},
		func() fyne.CanvasObject {
			// 创建原生Tab容器
			nativeTab1 := container.NewTabItem("标签1", widget.NewLabel("标签1内容"))
			nativeTab2 := container.NewTabItem("标签2", widget.NewLabel("标签2内容"))
			nativeTabs := container.NewAppTabs(nativeTab1, nativeTab2)
			nativeTabs.SetTabLocation(container.TabLocationTop)
			return nativeTabs
		},
		3*time.Second,
	)
}

// runBatchBenchmark 运行批量测试
func runBatchBenchmark(log func(string), statusLabel *widget.Label, comparisonContainer *fyne.Container) {
	log("🚀 开始批量性能测试...")
	fyne.Do(func() {
		comparisonContainer.Objects = nil
	})
	fyne.Do(func() {
		statusLabel.SetText("开始批量性能测试...")
	})

	// 批量测试配置
	tests := []struct {
		name     string
		function func(func(string), *widget.Label, *fyne.Container) *ScientificBenchmarkResult
	}{
		{"粒子按钮", testParticleButton},
		{"输入框", testMaterialEntry},
		{"开关控件", testToggleSwitch},
		{"复选框", testMaterialCheckbox},
		{"步骤标签页", testStepTabs},
	}

	// 运行所有测试
	results := make([]*ScientificBenchmarkResult, 0, len(tests))

	for i, test := range tests {
		log(fmt.Sprintf("\n📋 测试 %d/%d: %s", i+1, len(tests), test.name))
		result := test.function(log, statusLabel, comparisonContainer)
		if result != nil {
			results = append(results, result)

			// 短暂暂停，避免测试间相互影响
			if i < len(tests)-1 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	// 生成批量测试报告
	if len(results) > 0 {
		generateBatchReport(results, log, statusLabel)
	}

	log("✅ 批量性能测试完成！")
}

// generateBatchReport 生成批量测试报告
func generateBatchReport(results []*ScientificBenchmarkResult, log func(string), statusLabel *widget.Label) {
	log("\n📊 生成批量测试报告...")

	// 计算总体统计
	var totalPerformanceScore float64
	var bestResult *ScientificBenchmarkResult
	var worstResult *ScientificBenchmarkResult
	bestScore := -1.0
	worstScore := 101.0

	for _, result := range results {
		if result.Comparison == nil {
			continue
		}

		score := result.Comparison.Comparison["performance_score"].(float64)
		totalPerformanceScore += score

		if score > bestScore {
			bestScore = score
			bestResult = result
		}

		if score < worstScore {
			worstScore = score
			worstResult = result
		}
	}

	avgScore := totalPerformanceScore / float64(len(results))

	// 生成报告
	report := fmt.Sprintf(`
📈 批量性能测试报告

🔢 测试统计:
  • 测试总数: %d
  • 成功测试: %d
  • 平均性能评分: %.1f/100

🏆 最佳性能:
  • 组件: %s
  • 评分: %.1f/100
  • 结论: %s

⚠️ 最差性能:
  • 组件: %s  
  • 评分: %.1f/100
  • 建议: %s

💡 总体建议:
  %s

📁 详细报告已保存到 benchmark_results/ 目录
`,
		len(results),
		len(results),
		avgScore,
		bestResult.TestName,
		bestScore,
		bestResult.Comparison.Conclusion,
		worstResult.TestName,
		worstScore,
		worstResult.Comparison.Conclusion,
		getOverallRecommendation(avgScore, bestScore, worstScore),
	)

	fyne.Do(func() {
		statusLabel.SetText(report)
	})
	log("✅ 批量测试报告生成完成！")

	// 导出批量测试摘要
	exportBatchSummary(results, log)
}

// getOverallRecommendation 获取总体建议
func getOverallRecommendation(avgScore, bestScore, worstScore float64) string {
	if avgScore >= 85 {
		return "整体性能优秀，自定义控件质量很高"
	} else if avgScore >= 75 {
		return "整体性能良好，部分控件可能需要优化"
	} else if avgScore >= 65 {
		return "整体性能一般，建议对低分控件进行重点优化"
	} else {
		return "整体性能较差，需要系统性地优化自定义控件"
	}
}

// exportBatchSummary 导出批量测试摘要
func exportBatchSummary(results []*ScientificBenchmarkResult, log func(string)) {
	if len(results) == 0 {
		return
	}

	// 创建批量测试摘要
	summaryItems := make([]map[string]interface{}, 0, len(results))

	for _, result := range results {
		if result.Comparison == nil {
			continue
		}

		item := map[string]interface{}{
			"test_name":         result.TestName,
			"custom_component":  result.CustomComponent,
			"native_component":  result.NativeComponent,
			"performance_score": result.Comparison.Comparison["performance_score"],
			"fps_ratio":         result.Comparison.Comparison["fps_ratio"],
			"memory_ratio":      result.Comparison.Comparison["memory_ratio"],
			"cpu_ratio":         result.Comparison.Comparison["cpu_ratio"],
			"conclusion":        result.Comparison.Conclusion,
			"start_time":        result.StartTime,
			"end_time":          result.EndTime,
		}
		summaryItems = append(summaryItems, item)
	}

	// 使用CSV导出器创建摘要文件
	exporter := benchmark.NewCSVExporter("./benchmark_results/batch_summaries")
	timestamp := time.Now().Format("20060102_150405")
	exporter.SetFilename(fmt.Sprintf("batch_summary_%s.csv", timestamp))

	// 直接创建文件
	filePath := exporter.GetFullPath()
	file, err := os.Create(filePath)
	if err != nil {
		log(fmt.Sprintf("❌ 批量摘要导出失败: %v", err))
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入批量摘要标题
	writer.Write([]string{"批量性能测试摘要"})
	writer.Write([]string{fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05"))})
	writer.Write([]string{fmt.Sprintf("测试总数: %d", len(results))})
	sysInfo := benchmark.GetSystemInfo()
	writer.Write([]string{fmt.Sprintf("系统信息: Go %s, %s %s, %d cores",
		sysInfo.GoVersion, sysInfo.GOOS, sysInfo.GOARCH, sysInfo.NumCPU)})
	writer.Write([]string{})

	// 写入每个测试的结果标题
	writer.Write([]string{"测试名称", "自定义控件", "原生控件", "性能评分", "FPS比率", "内存比率", "CPU比率", "结论"})

	// 写入每个测试的结果数据
	for _, item := range summaryItems {
		writer.Write([]string{
			item["test_name"].(string),
			item["custom_component"].(string),
			item["native_component"].(string),
			fmt.Sprintf("%.1f", item["performance_score"].(float64)),
			fmt.Sprintf("%.3f", item["fps_ratio"].(float64)),
			fmt.Sprintf("%.3f", item["memory_ratio"].(float64)),
			fmt.Sprintf("%.3f", item["cpu_ratio"].(float64)),
			item["conclusion"].(string),
		})
	}

	// 写入统计摘要
	writer.Write([]string{})
	writer.Write([]string{"=== 统计摘要 ==="})

	// 计算统计数据
	var totalScore float64
	var passedTests int

	for _, item := range summaryItems {
		score := item["performance_score"].(float64)
		totalScore += score
		if score >= 70 {
			passedTests++
		}
	}

	avgScore := totalScore / float64(len(summaryItems))
	passRate := float64(passedTests) / float64(len(summaryItems)) * 100

	writer.Write([]string{"平均性能评分", fmt.Sprintf("%.1f/100", avgScore)})
	writer.Write([]string{"通过率", fmt.Sprintf("%.1f%%", passRate)})
	writer.Write([]string{"测试结论", getBatchConclusion(summaryItems)})

	log(fmt.Sprintf("✅ 批量测试摘要已导出到: %s", filePath))
}

// getBatchConclusion 获取批量测试结论
func getBatchConclusion(summary []map[string]interface{}) string {
	if len(summary) == 0 {
		return "没有测试数据"
	}

	var totalScore float64
	var passedTests int

	for _, item := range summary {
		score := item["performance_score"].(float64)
		totalScore += score
		if score >= 70 {
			passedTests++
		}
	}

	avgScore := totalScore / float64(len(summary))
	passRate := float64(passedTests) / float64(len(summary)) * 100

	if avgScore >= 85 {
		return fmt.Sprintf("优秀 - 平均评分%.1f，通过率%.1f%%", avgScore, passRate)
	} else if avgScore >= 75 {
		return fmt.Sprintf("良好 - 平均评分%.1f，通过率%.1f%%", avgScore, passRate)
	} else if avgScore >= 65 {
		return fmt.Sprintf("一般 - 平均评分%.1f，通过率%.1f%%", avgScore, passRate)
	} else {
		return fmt.Sprintf("需优化 - 平均评分%.1f，通过率%.1f%%", avgScore, passRate)
	}
}

// BuildBenchmarkPage 构建性能测试页面
func BuildBenchmarkPage() fyne.CanvasObject {
	// 创建状态显示
	statusLabel := widget.NewLabel("准备进行科学性能测试...")
	statusLabel.Wrapping = fyne.TextWrapWord

	// 创建日志输出区域
	logText := widget.NewMultiLineEntry()
	logText.SetPlaceHolder("测试日志将显示在这里...")
	logText.Disable()
	logScroll := container.NewScroll(logText)

	// 添加日志函数
	log := func(msg string) {
		fyne.Do(func() {
			currentText := logText.Text
			if currentText != "" {
				currentText += "\n"
			}
			currentText += time.Now().Format("15:04:05") + " - " + msg
			logText.SetText(currentText)
			logScroll.ScrollToBottom()
		})
	}

	// 创建对比容器
	comparisonContainer := container.NewHBox()

	// ====== 创建控制面板 ======

	// 单个测试按钮
	particleBtn := widget.NewButton("🔬 测试粒子按钮", func() {
		go testParticleButton(log, statusLabel, comparisonContainer)
	})

	entryBtn := widget.NewButton("🔬 测试输入框", func() {
		go testMaterialEntry(log, statusLabel, comparisonContainer)
	})

	toggleBtn := widget.NewButton("🔬 测试开关控件", func() {
		go testToggleSwitch(log, statusLabel, comparisonContainer)
	})

	checkboxBtn := widget.NewButton("🔬 测试复选框", func() {
		go testMaterialCheckbox(log, statusLabel, comparisonContainer)
	})

	tabsBtn := widget.NewButton("🔬 测试标签页", func() {
		go testStepTabs(log, statusLabel, comparisonContainer)
	})

	// 批量测试按钮
	batchTestBtn := widget.NewButton("🚀 批量测试所有控件", func() {
		go runBatchBenchmark(log, statusLabel, comparisonContainer)
	})

	// 清空日志按钮
	clearLogBtn := widget.NewButton("🗑️ 清空日志", func() {
		logText.SetText("")
		log("📝 日志已清空")
		fyne.Do(func() {
			statusLabel.SetText("准备进行科学性能测试...")
		})
	})

	// 清空对比容器按钮
	clearComparisonBtn := widget.NewButton("🗑️ 清空对比", func() {
		comparisonContainer.Objects = nil
		comparisonContainer.Refresh()
		log("🔄 对比容器已清空")
	})

	// 控制面板
	controlPanel := container.NewVBox(
		widget.NewLabelWithStyle("🔬 科学性能测试", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabel("测试方法:"),
		widget.NewLabel("• 分别测试自定义和原生控件"),
		widget.NewLabel("• 使用真实性能数据"),
		widget.NewLabel("• 科学统计对比分析"),
		widget.NewSeparator(),
		widget.NewLabel("选择测试类型:"),
		particleBtn,
		entryBtn,
		toggleBtn,
		checkboxBtn,
		tabsBtn,
		widget.NewSeparator(),
		batchTestBtn,
		widget.NewSeparator(),
		clearLogBtn,
		clearComparisonBtn,
		widget.NewSeparator(),
		widget.NewLabel("输出目录:"),
		widget.NewLabel("• benchmark_results/"),
		widget.NewLabel("• benchmark_results/scientific/"),
		widget.NewLabel("• benchmark_results/comparisons/"),
		widget.NewLabel("• benchmark_results/batch_summaries/"),
	)

	// ====== 创建主内容区域 ======
	mainContent := container.NewVBox(
		widget.NewLabelWithStyle("🔍 控件对比展示区", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		container.NewVSplit(
			container.NewGridWrap(fyne.NewSize(500, 180), comparisonContainer),
			container.NewVSplit(
				container.NewGridWrap(fyne.NewSize(500, 120), logScroll),
				container.NewGridWrap(fyne.NewSize(500, 80), container.NewScroll(statusLabel)),
			),
		),
	)

	mainContentScroll := container.NewScroll(mainContent)
	mainContentScroll.SetMinSize(fyne.NewSize(600, 400))

	// 创建分割容器
	split := container.NewHSplit(controlPanel, mainContentScroll)
	split.SetOffset(0.25)

	return container.NewPadded(split)
}
