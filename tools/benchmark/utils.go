// tools/benchmark/utils.go
package benchmark

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// RunGarbageCollection 运行垃圾回收并测量耗时
func RunGarbageCollection() (duration time.Duration) {
	start := time.Now()
	runtime.GC()
	return time.Since(start)
}

// GetMemoryStats 获取当前内存统计
func GetMemoryStats() map[string]float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]float64{
		"alloc_mb":       float64(memStats.Alloc) / 1024 / 1024,
		"total_alloc_mb": float64(memStats.TotalAlloc) / 1024 / 1024,
		"sys_mb":         float64(memStats.Sys) / 1024 / 1024,
		"heap_alloc_mb":  float64(memStats.HeapAlloc) / 1024 / 1024,
		"heap_sys_mb":    float64(memStats.HeapSys) / 1024 / 1024,
		"num_gc":         float64(memStats.NumGC),
		"gc_pause_ms":    float64(memStats.PauseTotalNs) / 1_000_000,
	}
}

// PrintMemoryStats 打印内存统计信息
func PrintMemoryStats(label string) {
	stats := GetMemoryStats()

	fmt.Printf("\n📊 %s 内存统计:\n", label)
	fmt.Printf("  当前分配: %.2f MB\n", stats["alloc_mb"])
	fmt.Printf("  累计分配: %.2f MB\n", stats["total_alloc_mb"])
	fmt.Printf("  系统内存: %.2f MB\n", stats["sys_mb"])
	fmt.Printf("  堆分配:   %.2f MB\n", stats["heap_alloc_mb"])
	fmt.Printf("  堆系统:   %.2f MB\n", stats["heap_sys_mb"])
	fmt.Printf("  GC次数:   %.0f\n", stats["num_gc"])
	fmt.Printf("  GC暂停:   %.2f ms\n", stats["gc_pause_ms"])
}

// MeasureExecutionTime 测量函数执行时间
func MeasureExecutionTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

// BenchmarkComponent 基准测试组件性能的通用函数
func BenchmarkComponent(name, componentType string, setupFunc, testFunc, cleanupFunc func()) (*PerformanceMetric, error) {
	fmt.Printf("🔧 开始测试组件: %s (%s)\n", name, componentType)

	// 运行GC确保测试环境干净
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	// 记录初始内存
	initialStats := GetMemoryStats()

	// 执行设置函数
	if setupFunc != nil {
		fmt.Println("  设置测试环境...")
		setupFunc()
	}

	// 等待稳定
	time.Sleep(200 * time.Millisecond)

	// 记录测试前内存
	preTestStats := GetMemoryStats()

	// 执行测试函数并测量时间
	fmt.Println("  执行测试...")
	executionTime := MeasureExecutionTime(testFunc)

	// 记录测试后内存
	postTestStats := GetMemoryStats()

	// 执行清理函数
	if cleanupFunc != nil {
		fmt.Println("  清理测试环境...")
		cleanupFunc()
	}

	// 再次运行GC查看最终内存
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	finalStats := GetMemoryStats()

	// 创建性能指标
	metric := NewPerformanceMetric(name, componentType, "benchmark")
	metric.RenderTimeMS = executionTime.Seconds() * 1000

	// 计算内存增量
	metric.MemoryUsageMB = postTestStats["alloc_mb"] - preTestStats["alloc_mb"]
	metric.MemoryAllocMB = finalStats["alloc_mb"] - initialStats["alloc_mb"]

	// 获取其他系统信息
	metric.Goroutines = runtime.NumGoroutine()
	metric.NumCores = runtime.NumCPU()

	fmt.Printf("✅ 测试完成 - 执行时间: %.2fms, 内存增量: %.2fMB\n",
		metric.RenderTimeMS, metric.MemoryUsageMB)

	return metric, nil
}

// CreateTestReport 创建测试报告
func CreateTestReport(metrics []*PerformanceMetric, outputPath string) error {
	if len(metrics) == 0 {
		return fmt.Errorf("no metrics to report")
	}

	exporter := NewCSVExporter(filepath.Dir(outputPath))
	exporter.SetFilename(filepath.Base(outputPath))

	return exporter.ExportMetrics(metrics, nil)
}

// 组件对比结果
type ComponentComparison struct {
	CustomMetrics []*PerformanceMetric
	NativeMetrics []*PerformanceMetric
	CustomSummary map[string]interface{}
	NativeSummary map[string]interface{}
	Comparison    map[string]interface{}
	Conclusion    string
}

// CompareComponents 对比两个组件性能（科学方法）
func CompareComponents(customMetrics, nativeMetrics []*PerformanceMetric) *ComponentComparison {
	if len(customMetrics) == 0 || len(nativeMetrics) == 0 {
		return nil
	}

	// 1. 分别计算统计摘要
	customSummary := calculateComponentSummary(customMetrics, "custom")
	nativeSummary := calculateComponentSummary(nativeMetrics, "native")

	// 2. 科学对比分析
	comparison := scientificComparison(customSummary, nativeSummary)

	// 3. 生成结论
	conclusion := generateConclusion(comparison)

	return &ComponentComparison{
		CustomMetrics: customMetrics,
		NativeMetrics: nativeMetrics,
		CustomSummary: customSummary,
		NativeSummary: nativeSummary,
		Comparison:    comparison,
		Conclusion:    conclusion,
	}
}

// calculateComponentSummary 计算组件性能统计摘要
func calculateComponentSummary(metrics []*PerformanceMetric, componentType string) map[string]interface{} {
	if len(metrics) == 0 {
		return nil
	}

	// 提取关键指标
	var fpsValues, memoryValues, cpuValues []float64
	var totalRenderTime float64

	for _, metric := range metrics {
		fpsValues = append(fpsValues, metric.FPS)
		memoryValues = append(memoryValues, metric.MemoryUsageMB)
		cpuValues = append(cpuValues, metric.CPUPercent)
		totalRenderTime += metric.RenderTimeMS
	}

	// 计算统计指标
	fpsStats := calculateStatistics(fpsValues)
	memoryStats := calculateStatistics(memoryValues)
	cpuStats := calculateStatistics(cpuValues)

	avgRenderTime := totalRenderTime / float64(len(metrics))

	return map[string]interface{}{
		"component_type": componentType,
		"sample_count":   len(metrics),

		"fps_avg": fpsStats["avg"],
		"fps_min": fpsStats["min"],
		"fps_max": fpsStats["max"],
		"fps_std": fpsStats["std"],
		"fps_cv":  fpsStats["cv"], // 变异系数

		"memory_avg": memoryStats["avg"],
		"memory_min": memoryStats["min"],
		"memory_max": memoryStats["max"],
		"memory_std": memoryStats["std"],
		"memory_cv":  memoryStats["cv"],

		"cpu_avg": cpuStats["avg"],
		"cpu_min": cpuStats["min"],
		"cpu_max": cpuStats["max"],
		"cpu_std": cpuStats["std"],
		"cpu_cv":  cpuStats["cv"],

		"render_time_avg": avgRenderTime,
	}
}

// calculateStatistics 计算统计指标
func calculateStatistics(values []float64) map[string]float64 {
	if len(values) == 0 {
		return nil
	}

	// 计算基本统计
	var sum, min, max float64
	min = math.MaxFloat64

	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	avg := sum / float64(len(values))

	// 计算标准差
	var variance float64
	for _, v := range values {
		diff := v - avg
		variance += diff * diff
	}
	variance /= float64(len(values))
	std := math.Sqrt(variance)

	// 计算变异系数（标准差/均值）
	cv := 0.0
	if avg != 0 {
		cv = (std / avg) * 100
	}

	return map[string]float64{
		"avg": math.Round(avg*100) / 100,
		"min": math.Round(min*100) / 100,
		"max": math.Round(max*100) / 100,
		"std": math.Round(std*100) / 100,
		"cv":  math.Round(cv*100) / 100,
	}
}

// scientificComparison 科学对比两个组件的性能
func scientificComparison(customSummary, nativeSummary map[string]interface{}) map[string]interface{} {
	if customSummary == nil || nativeSummary == nil {
		return nil
	}

	// 提取关键指标
	customFPSAvg := customSummary["fps_avg"].(float64)
	customMemoryAvg := customSummary["memory_avg"].(float64)
	customCPUAvg := customSummary["cpu_avg"].(float64)

	nativeFPSAvg := nativeSummary["fps_avg"].(float64)
	nativeMemoryAvg := nativeSummary["memory_avg"].(float64)
	nativeCPUAvg := nativeSummary["cpu_avg"].(float64)

	// 计算相对性能
	fpsRatio := customFPSAvg / nativeFPSAvg
	memoryRatio := customMemoryAvg / nativeMemoryAvg
	cpuRatio := customCPUAvg / nativeCPUAvg

	// 计算性能差异百分比
	fpsDiffPercent := (customFPSAvg - nativeFPSAvg) / nativeFPSAvg * 100
	memoryDiffPercent := (customMemoryAvg - nativeMemoryAvg) / nativeMemoryAvg * 100
	cpuDiffPercent := (customCPUAvg - nativeCPUAvg) / nativeCPUAvg * 100

	// 计算综合性能评分
	performanceScore := calculatePerformanceScore(
		fpsRatio, memoryRatio, cpuRatio,
		customSummary["fps_cv"].(float64),
		customSummary["memory_cv"].(float64),
		customSummary["cpu_cv"].(float64),
	)

	// 显著性分析
	significance := analyzeSignificance(
		customSummary, nativeSummary,
		fpsDiffPercent, memoryDiffPercent, cpuDiffPercent,
	)

	return map[string]interface{}{
		// 性能比率
		"fps_ratio":    math.Round(fpsRatio*1000) / 1000,
		"memory_ratio": math.Round(memoryRatio*1000) / 1000,
		"cpu_ratio":    math.Round(cpuRatio*1000) / 1000,

		// 性能差异百分比
		"fps_diff_percent":    math.Round(fpsDiffPercent*100) / 100,
		"memory_diff_percent": math.Round(memoryDiffPercent*100) / 100,
		"cpu_diff_percent":    math.Round(cpuDiffPercent*100) / 100,

		// 性能评分
		"performance_score": performanceScore,

		// 显著性分析
		"significance": significance,

		// 性能分类
		"performance_category": classifyPerformance(
			fpsRatio, memoryRatio, cpuRatio,
			fpsDiffPercent, memoryDiffPercent, cpuDiffPercent,
		),
	}
}

// calculatePerformanceScore 计算综合性能评分
func calculatePerformanceScore(fpsRatio, memoryRatio, cpuRatio, fpsCV, memoryCV, cpuCV float64) float64 {
	// 权重分配：CPU 40%，FPS 30%，内存 30%
	// 变异系数惩罚：稳定性越差，扣分越多

	// FPS得分（越高越好）
	fpsScore := 0.0
	if fpsRatio >= 1.0 {
		fpsScore = 100 // 持平或更好
	} else {
		fpsScore = fpsRatio * 100 // 按比例得分
	}

	// 内存得分（越低越好）
	memoryScore := 0.0
	if memoryRatio <= 1.0 {
		memoryScore = 100 // 持平或更好
	} else {
		memoryScore = (1.0 / memoryRatio) * 100 // 内存使用越多，得分越低
	}

	// CPU得分（越低越好）
	cpuScore := 0.0
	if cpuRatio <= 1.0 {
		cpuScore = 100 // 持平或更好
	} else {
		cpuScore = (1.0 / cpuRatio) * 100 // CPU使用越多，得分越低
	}

	// 稳定性惩罚（变异系数越高，扣分越多）
	stabilityPenalty := (fpsCV * 0.5) + (memoryCV * 0.3) + (cpuCV * 0.2)

	// 综合得分
	totalScore := (fpsScore * 0.3) + (memoryScore * 0.3) + (cpuScore * 0.4) - stabilityPenalty

	// 确保在0-100范围内
	if totalScore < 0 {
		totalScore = 0
	}
	if totalScore > 100 {
		totalScore = 100
	}

	return math.Round(totalScore*10) / 10
}

// analyzeSignificance 分析性能差异的显著性
func analyzeSignificance(customSummary, nativeSummary map[string]interface{},
	fpsDiffPercent, memoryDiffPercent, cpuDiffPercent float64) map[string]interface{} {

	// 提取变异系数
	customFPSCV := customSummary["fps_cv"].(float64)
	customMemoryCV := customSummary["memory_cv"].(float64)
	customCPUCV := customSummary["cpu_cv"].(float64)

	// 判断显著性
	isFPSSignificant := math.Abs(fpsDiffPercent) > (customFPSCV * 2) // 差异大于2倍变异系数
	isMemorySignificant := math.Abs(memoryDiffPercent) > (customMemoryCV * 2)
	isCPUSignificant := math.Abs(cpuDiffPercent) > (customCPUCV * 2)

	return map[string]interface{}{
		"fps_significant":    isFPSSignificant,
		"memory_significant": isMemorySignificant,
		"cpu_significant":    isCPUSignificant,

		"fps_confidence":    calculateConfidenceLevel(math.Abs(fpsDiffPercent), customFPSCV),
		"memory_confidence": calculateConfidenceLevel(math.Abs(memoryDiffPercent), customMemoryCV),
		"cpu_confidence":    calculateConfidenceLevel(math.Abs(cpuDiffPercent), customCPUCV),
	}
}

// calculateConfidenceLevel 计算置信水平
func calculateConfidenceLevel(diffPercent, cv float64) string {
	if cv == 0 {
		return "high"
	}

	ratio := diffPercent / cv

	if ratio >= 3 {
		return "very high"
	} else if ratio >= 2 {
		return "high"
	} else if ratio >= 1 {
		return "medium"
	} else {
		return "low"
	}
}

// classifyPerformance 分类性能表现
func classifyPerformance(fpsRatio, memoryRatio, cpuRatio,
	fpsDiffPercent, memoryDiffPercent, cpuDiffPercent float64) map[string]interface{} {

	// FPS分类
	fpsCategory := ""
	if fpsDiffPercent >= 10 {
		fpsCategory = "excellent" // 显著优于
	} else if fpsDiffPercent >= 0 {
		fpsCategory = "good" // 略优或持平
	} else if fpsDiffPercent >= -10 {
		fpsCategory = "acceptable" // 略差但可接受
	} else if fpsDiffPercent >= -30 {
		fpsCategory = "poor" // 较差
	} else {
		fpsCategory = "bad" // 很差
	}

	// 内存分类
	memoryCategory := ""
	if memoryDiffPercent <= -10 {
		memoryCategory = "excellent" // 显著节省内存
	} else if memoryDiffPercent <= 0 {
		memoryCategory = "good" // 略省或持平
	} else if memoryDiffPercent <= 20 {
		memoryCategory = "acceptable" // 略多但可接受
	} else if memoryDiffPercent <= 50 {
		memoryCategory = "poor" // 较多
	} else {
		memoryCategory = "bad" // 很多
	}

	// CPU分类
	cpuCategory := ""
	if cpuDiffPercent <= -10 {
		cpuCategory = "excellent" // 显著节省CPU
	} else if cpuDiffPercent <= 0 {
		cpuCategory = "good" // 略省或持平
	} else if cpuDiffPercent <= 20 {
		cpuCategory = "acceptable" // 略多但可接受
	} else if cpuDiffPercent <= 50 {
		cpuCategory = "poor" // 较多
	} else {
		cpuCategory = "bad" // 很多
	}

	return map[string]interface{}{
		"fps":    fpsCategory,
		"memory": memoryCategory,
		"cpu":    cpuCategory,
	}
}

// generateConclusion 生成测试结论
func generateConclusion(comparison map[string]interface{}) string {
	if comparison == nil {
		return "数据不足，无法生成结论"
	}

	performanceScore := comparison["performance_score"].(float64)
	fpsDiffPercent := comparison["fps_diff_percent"].(float64)
	memoryDiffPercent := comparison["memory_diff_percent"].(float64)
	cpuDiffPercent := comparison["cpu_diff_percent"].(float64)

	performanceCategory := comparison["performance_category"].(map[string]interface{})
	fpsCategory := performanceCategory["fps"].(string)
	memoryCategory := performanceCategory["memory"].(string)
	cpuCategory := performanceCategory["cpu"].(string)

	// 根据综合评分和各项指标生成结论
	if performanceScore >= 90 {
		return fmt.Sprintf("✅ 性能优秀 (%.1f分)。自定义控件在保持视觉效果的同时，性能表现优于或接近原生控件。", performanceScore)
	} else if performanceScore >= 80 {
		return fmt.Sprintf("🟡 性能良好 (%.1f分)。自定义控件性能可接受，FPS: %s, 内存: %s, CPU: %s。",
			performanceScore, fpsCategory, memoryCategory, cpuCategory)
	} else if performanceScore >= 70 {
		return fmt.Sprintf("🟡 性能一般 (%.1f分)。存在一定的性能开销，建议优化。FPS差异: %.1f%%, 内存开销: %.1f%%, CPU开销: %.1f%%。",
			performanceScore, fpsDiffPercent, memoryDiffPercent, cpuDiffPercent)
	} else {
		return fmt.Sprintf("🔴 性能较差 (%.1f分)。性能开销较大，需要重点优化。FPS显著降低，内存和CPU使用明显增加。",
			performanceScore)
	}
}

// EnsureOutputDir 确保输出目录存在
func EnsureOutputDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// GetTimestamp 获取时间戳字符串
func GetTimestamp() string {
	return time.Now().Format("2006-01-02_15-04-05")
}

// PrintComparison 打印对比结果
func PrintComparison(comparison *ComponentComparison) {
	if comparison == nil {
		fmt.Println("对比结果为空")
		return
	}

	fmt.Println("\n📊 ========== 性能对比结果 ==========")

	// 打印自定义控件性能
	customSummary := comparison.CustomSummary
	fmt.Printf("\n🔧 自定义控件性能:\n")
	fmt.Printf("  样本数: %d\n", customSummary["sample_count"])
	fmt.Printf("  FPS: %.1f (%.1f-%.1f) σ=%.2f CV=%.1f%%\n",
		customSummary["fps_avg"], customSummary["fps_min"],
		customSummary["fps_max"], customSummary["fps_std"], customSummary["fps_cv"])
	fmt.Printf("  内存: %.2fMB (%.2f-%.2f) σ=%.2f CV=%.1f%%\n",
		customSummary["memory_avg"], customSummary["memory_min"],
		customSummary["memory_max"], customSummary["memory_std"], customSummary["memory_cv"])
	fmt.Printf("  CPU: %.1f%% (%.1f-%.1f) σ=%.2f CV=%.1f%%\n",
		customSummary["cpu_avg"], customSummary["cpu_min"],
		customSummary["cpu_max"], customSummary["cpu_std"], customSummary["cpu_cv"])

	// 打印原生控件性能
	nativeSummary := comparison.NativeSummary
	fmt.Printf("\n🔧 原生控件性能:\n")
	fmt.Printf("  样本数: %d\n", nativeSummary["sample_count"])
	fmt.Printf("  FPS: %.1f (%.1f-%.1f) σ=%.2f CV=%.1f%%\n",
		nativeSummary["fps_avg"], nativeSummary["fps_min"],
		nativeSummary["fps_max"], nativeSummary["fps_std"], nativeSummary["fps_cv"])
	fmt.Printf("  内存: %.2fMB (%.2f-%.2f) σ=%.2f CV=%.1f%%\n",
		nativeSummary["memory_avg"], nativeSummary["memory_min"],
		nativeSummary["memory_max"], nativeSummary["memory_std"], nativeSummary["memory_cv"])
	fmt.Printf("  CPU: %.1f%% (%.1f-%.1f) σ=%.2f CV=%.1f%%\n",
		nativeSummary["cpu_avg"], nativeSummary["cpu_min"],
		nativeSummary["cpu_max"], nativeSummary["cpu_std"], nativeSummary["cpu_cv"])

	// 打印对比结果
	comp := comparison.Comparison
	fmt.Printf("\n📈 性能对比分析:\n")
	fmt.Printf("  FPS比率: %.3f (差异: %.1f%%)\n",
		comp["fps_ratio"], comp["fps_diff_percent"])
	fmt.Printf("  内存比率: %.3f (差异: %.1f%%)\n",
		comp["memory_ratio"], comp["memory_diff_percent"])
	fmt.Printf("  CPU比率: %.3f (差异: %.1f%%)\n",
		comp["cpu_ratio"], comp["cpu_diff_percent"])

	fmt.Printf("\n🏆 综合性能评分: %.1f/100\n", comp["performance_score"])

	// 打印显著性分析
	sig := comp["significance"].(map[string]interface{})
	fmt.Printf("\n🔍 显著性分析:\n")
	fmt.Printf("  FPS显著性: %v (置信度: %s)\n",
		sig["fps_significant"], sig["fps_confidence"])
	fmt.Printf("  内存显著性: %v (置信度: %s)\n",
		sig["memory_significant"], sig["memory_confidence"])
	fmt.Printf("  CPU显著性: %v (置信度: %s)\n",
		sig["cpu_significant"], sig["cpu_confidence"])

	fmt.Printf("\n💡 结论: %s\n", comparison.Conclusion)
	fmt.Println("====================================\n")
}
