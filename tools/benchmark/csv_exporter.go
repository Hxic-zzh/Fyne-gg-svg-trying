// tools/benchmark/csv_exporter.go
package benchmark

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// CSVExporter CSV导出器
type CSVExporter struct {
	outputDir     string
	filename      string
	includeHeader bool
}

// NewCSVExporter 创建CSV导出器
func NewCSVExporter(outputDir string) *CSVExporter {
	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		// 如果创建失败，使用当前目录
		outputDir = "."
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("benchmark_%s.csv", timestamp)

	return &CSVExporter{
		outputDir:     outputDir,
		filename:      filename,
		includeHeader: true,
	}
}

// ExportMetrics 导出性能指标到CSV
func (e *CSVExporter) ExportMetrics(metrics []*PerformanceMetric, summary map[string]interface{}) error {
	if len(metrics) == 0 {
		return fmt.Errorf("no metrics to export")
	}

	filePath := filepath.Join(e.outputDir, e.filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if e.includeHeader {
		headers := []string{
			"timestamp",
			"component_name",
			"component_type",
			"test_scenario",
			"fps",
			"memory_usage_mb",
			"memory_alloc_mb",
			"num_gc",
			"gc_time_ms",
			"cpu_percent",
			"render_time_ms",
			"update_time_ms",
			"goroutines",
			"num_cores",
		}
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write header: %v", err)
		}
	}

	// 写入数据行
	for _, metric := range metrics {
		record := []string{
			metric.Timestamp.Format("2006-01-02 15:04:05.000"),
			metric.ComponentName,
			metric.ComponentType,
			metric.TestScenario,
			strconv.FormatFloat(metric.FPS, 'f', 2, 64),
			strconv.FormatFloat(metric.MemoryUsageMB, 'f', 2, 64),
			strconv.FormatFloat(metric.MemoryAllocMB, 'f', 2, 64),
			strconv.FormatUint(uint64(metric.NumGC), 10),
			strconv.FormatFloat(metric.GCTimeMS, 'f', 2, 64),
			strconv.FormatFloat(metric.CPUPercent, 'f', 2, 64),
			strconv.FormatFloat(metric.RenderTimeMS, 'f', 2, 64),
			strconv.FormatFloat(metric.UpdateTimeMS, 'f', 2, 64),
			strconv.Itoa(metric.Goroutines),
			strconv.Itoa(metric.NumCores),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
	}

	// 写入摘要信息
	if summary != nil {
		if err := e.writeSummary(writer, summary); err != nil {
			return err
		}
	}

	fmt.Printf("✅ 性能数据已导出到: %s\n", filePath)
	return nil
}

// writeSummary 写入摘要信息
// tools/benchmark/csv_exporter.go
// 修复第124行的语法错误

func (e *CSVExporter) writeSummary(writer *csv.Writer, summary map[string]interface{}) error {
	// 写入空行分隔
	if err := writer.Write([]string{}); err != nil {
		return err
	}

	// 写入摘要标题
	if err := writer.Write([]string{"=== 性能测试摘要 ==="}); err != nil {
		return err
	}

	// 判断摘要类型并写入相应内容
	_, hasCustom := summary["custom_summary"]
	_, hasNative := summary["native_summary"]
	comparisonResult, hasComparison := summary["comparison_result"]

	if hasCustom && hasNative {
		// 这是对比测试的摘要
		return e.writeComparisonSummary(writer, summary)
	} else if hasComparison {
		// 包含对比结果的摘要
		return e.writeComparisonResult(writer, summary, comparisonResult.(map[string]interface{}))
	} else {
		// 常规测试摘要
		return e.writeRegularSummary(writer, summary)
	}
}

// writeRegularSummary 写入常规测试摘要
func (e *CSVExporter) writeRegularSummary(writer *csv.Writer, summary map[string]interface{}) error {
	// 直接使用 summary 中的值，不存储到临时变量

	summaryRows := [][]string{
		{"测试名称", fmt.Sprintf("%v", summary["test_name"])},
		{"组件名称", fmt.Sprintf("%v", summary["component_name"])},
		{"组件类型", fmt.Sprintf("%v", summary["component_type"])},
		{"测试场景", fmt.Sprintf("%v", summary["scenario"])},
		{"样本数量", fmt.Sprintf("%v", summary["total_samples"])},
		{"测试时长", fmt.Sprintf("%.1f 秒", summary["duration_seconds"])},
		{"总帧数", fmt.Sprintf("%v", summary["total_frames"])},
		{},
		{"=== 性能统计 ==="},
		{"平均帧率 (FPS)", fmt.Sprintf("%.2f", summary["avg_fps"])},
		{"最低帧率 (FPS)", fmt.Sprintf("%.2f", summary["min_fps"])},
		{"最高帧率 (FPS)", fmt.Sprintf("%.2f", summary["max_fps"])},
		{"帧率标准差", fmt.Sprintf("%.2f", summary["fps_std_dev"])},
		{},
		{"平均内存使用 (MB)", fmt.Sprintf("%.2f", summary["avg_memory_mb"])},
		{"最低内存 (MB)", fmt.Sprintf("%.2f", summary["min_memory_mb"])},
		{"最高内存 (MB)", fmt.Sprintf("%.2f", summary["max_memory_mb"])},
		{"内存标准差", fmt.Sprintf("%.2f", summary["memory_std_dev"])},
		{},
		{"平均CPU使用率 (%)", fmt.Sprintf("%.2f", summary["avg_cpu_percent"])},
		{"CPU标准差", fmt.Sprintf("%.2f", summary["cpu_std_dev"])},
		{},
		{"=== 系统信息 ==="},
	}

	// 写入基础行
	for _, row := range summaryRows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	// 添加系统信息
	if sysInfo, ok := summary["system_info"]; ok {
		info := sysInfo.(SystemInfo)
		sysRows := [][]string{
			{"Go版本", info.GoVersion},
			{"CPU核心数", fmt.Sprintf("%d", info.NumCPU)},
			{"操作系统", info.GOOS},
			{"系统架构", info.GOARCH},
		}

		for _, row := range sysRows {
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	// 添加性能评估
	return e.writePerformanceEvaluation(writer, summary)
}

// writeComparisonResult 写入对比结果
func (e *CSVExporter) writeComparisonResult(writer *csv.Writer, summary, comparison map[string]interface{}) error {
	// 写入基础信息
	baseRows := [][]string{
		{"测试名称", fmt.Sprintf("%v", summary["test_name"])},
		{"测试类型", "性能对比测试"},
		{},
		{"=== 对比结果 ==="},
		{"性能评分", fmt.Sprintf("%.1f/100", comparison["performance_score"])},
		{"FPS比率", fmt.Sprintf("%.3f", comparison["fps_ratio"])},
		{"内存比率", fmt.Sprintf("%.3f", comparison["memory_ratio"])},
		{"CPU比率", fmt.Sprintf("%.3f", comparison["cpu_ratio"])},
		{},
		{"FPS差异百分比", fmt.Sprintf("%.1f%%", comparison["fps_diff_percent"])},
		{"内存差异百分比", fmt.Sprintf("%.1f%%", comparison["memory_diff_percent"])},
		{"CPU差异百分比", fmt.Sprintf("%.1f%%", comparison["cpu_diff_percent"])},
	}

	for _, row := range baseRows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	// 写入显著性分析
	if significance, ok := comparison["significance"]; ok {
		sig := significance.(map[string]interface{})
		sigRows := [][]string{
			{},
			{"=== 显著性分析 ==="},
			{"FPS显著性", fmt.Sprintf("%v", sig["fps_significant"])},
			{"FPS置信度", fmt.Sprintf("%v", sig["fps_confidence"])},
			{"内存显著性", fmt.Sprintf("%v", sig["memory_significant"])},
			{"内存置信度", fmt.Sprintf("%v", sig["memory_confidence"])},
			{"CPU显著性", fmt.Sprintf("%v", sig["cpu_significant"])},
			{"CPU置信度", fmt.Sprintf("%v", sig["cpu_confidence"])},
		}

		for _, row := range sigRows {
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	// 写入性能分类
	if perfCategory, ok := comparison["performance_category"]; ok {
		category := perfCategory.(map[string]interface{})
		categoryRows := [][]string{
			{},
			{"=== 性能分类 ==="},
			{"FPS分类", fmt.Sprintf("%v", category["fps"])},
			{"内存分类", fmt.Sprintf("%v", category["memory"])},
			{"CPU分类", fmt.Sprintf("%v", category["cpu"])},
		}

		for _, row := range categoryRows {
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	// 写入结论
	conclusionRows := [][]string{
		{},
		{"=== 测试结论 ==="},
		{"结论", fmt.Sprintf("%v", comparison["conclusion"])},
	}

	for _, row := range conclusionRows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	// 写入系统信息
	if sysInfo, ok := summary["system_info"]; ok {
		info := sysInfo.(SystemInfo)
		sysRows := [][]string{
			{},
			{"=== 系统信息 ==="},
			{"Go版本", info.GoVersion},
			{"CPU核心数", fmt.Sprintf("%d", info.NumCPU)},
			{"操作系统", info.GOOS},
			{"系统架构", info.GOARCH},
		}

		for _, row := range sysRows {
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	return nil
}

// writeComparisonSummary 写入对比测试的详细摘要
func (e *CSVExporter) writeComparisonSummary(writer *csv.Writer, summary map[string]interface{}) error {
	// 写入标题
	titleRows := [][]string{
		{"测试名称", fmt.Sprintf("%v", summary["test_name"])},
		{"测试类型", "科学对比测试"},
		{},
		{"=== 自定义组件性能 ==="},
	}

	customSummary := summary["custom_summary"].(map[string]interface{})
	nativeSummary := summary["native_summary"].(map[string]interface{})
	comparison := summary["comparison"].(map[string]interface{})
	conclusion := summary["conclusion"].(string)

	// 写入自定义组件性能
	customRows := [][]string{
		{"组件类型", "custom"},
		{"样本数量", fmt.Sprintf("%v", customSummary["sample_count"])},
		{"平均FPS", fmt.Sprintf("%.2f", customSummary["fps_avg"])},
		{"FPS范围", fmt.Sprintf("%.2f-%.2f", customSummary["fps_min"], customSummary["fps_max"])},
		{"FPS标准差", fmt.Sprintf("%.2f", customSummary["fps_std"])},
		{"FPS变异系数", fmt.Sprintf("%.1f%%", customSummary["fps_cv"])},
		{},
		{"平均内存(MB)", fmt.Sprintf("%.2f", customSummary["memory_avg"])},
		{"内存范围", fmt.Sprintf("%.2f-%.2f", customSummary["memory_min"], customSummary["memory_max"])},
		{"内存标准差", fmt.Sprintf("%.2f", customSummary["memory_std"])},
		{"内存变异系数", fmt.Sprintf("%.1f%%", customSummary["memory_cv"])},
		{},
		{"平均CPU(%)", fmt.Sprintf("%.2f", customSummary["cpu_avg"])},
		{"CPU范围", fmt.Sprintf("%.2f-%.2f", customSummary["cpu_min"], customSummary["cpu_max"])},
		{"CPU标准差", fmt.Sprintf("%.2f", customSummary["cpu_std"])},
		{"CPU变异系数", fmt.Sprintf("%.1f%%", customSummary["cpu_cv"])},
	}

	// 写入原生组件性能
	nativeRows := [][]string{
		{},
		{"=== 原生组件性能 ==="},
		{"组件类型", "native"},
		{"样本数量", fmt.Sprintf("%v", nativeSummary["sample_count"])},
		{"平均FPS", fmt.Sprintf("%.2f", nativeSummary["fps_avg"])},
		{"FPS范围", fmt.Sprintf("%.2f-%.2f", nativeSummary["fps_min"], nativeSummary["fps_max"])},
		{"FPS标准差", fmt.Sprintf("%.2f", nativeSummary["fps_std"])},
		{"FPS变异系数", fmt.Sprintf("%.1f%%", nativeSummary["fps_cv"])},
		{},
		{"平均内存(MB)", fmt.Sprintf("%.2f", nativeSummary["memory_avg"])},
		{"内存范围", fmt.Sprintf("%.2f-%.2f", nativeSummary["memory_min"], nativeSummary["memory_max"])},
		{"内存标准差", fmt.Sprintf("%.2f", nativeSummary["memory_std"])},
		{"内存变异系数", fmt.Sprintf("%.1f%%", nativeSummary["memory_cv"])},
		{},
		{"平均CPU(%)", fmt.Sprintf("%.2f", nativeSummary["cpu_avg"])},
		{"CPU范围", fmt.Sprintf("%.2f-%.2f", nativeSummary["cpu_min"], nativeSummary["cpu_max"])},
		{"CPU标准差", fmt.Sprintf("%.2f", nativeSummary["cpu_std"])},
		{"CPU变异系数", fmt.Sprintf("%.1f%%", nativeSummary["cpu_cv"])},
	}

	// 写入对比分析
	comparisonRows := [][]string{
		{},
		{"=== 性能对比分析 ==="},
		{"FPS比率", fmt.Sprintf("%.3f", comparison["fps_ratio"])},
		{"FPS差异", fmt.Sprintf("%.1f%%", comparison["fps_diff_percent"])},
		{},
		{"内存比率", fmt.Sprintf("%.3f", comparison["memory_ratio"])},
		{"内存差异", fmt.Sprintf("%.1f%%", comparison["memory_diff_percent"])},
		{},
		{"CPU比率", fmt.Sprintf("%.3f", comparison["cpu_ratio"])},
		{"CPU差异", fmt.Sprintf("%.1f%%", comparison["cpu_diff_percent"])},
		{},
		{"综合性能评分", fmt.Sprintf("%.1f/100", comparison["performance_score"])},
	}

	// 写入显著性分析
	if significance, ok := comparison["significance"]; ok {
		sig := significance.(map[string]interface{})
		sigRows := [][]string{
			{},
			{"=== 显著性分析 ==="},
			{"FPS显著性", fmt.Sprintf("%v", sig["fps_significant"])},
			{"FPS置信度", fmt.Sprintf("%v", sig["fps_confidence"])},
			{"内存显著性", fmt.Sprintf("%v", sig["memory_significant"])},
			{"内存置信度", fmt.Sprintf("%v", sig["memory_confidence"])},
			{"CPU显著性", fmt.Sprintf("%v", sig["cpu_significant"])},
			{"CPU置信度", fmt.Sprintf("%v", sig["cpu_confidence"])},
		}
		comparisonRows = append(comparisonRows, sigRows...)
	}

	// 写入结论
	conclusionRows := [][]string{
		{},
		{"=== 测试结论 ==="},
		{"结论", conclusion},
	}

	// 合并所有行并写入
	allRows := append(titleRows, customRows...)
	allRows = append(allRows, nativeRows...)
	allRows = append(allRows, comparisonRows...)
	allRows = append(allRows, conclusionRows...)

	// 写入系统信息
	if sysInfo, ok := summary["system_info"]; ok {
		info := sysInfo.(SystemInfo)
		sysRows := [][]string{
			{},
			{"=== 系统信息 ==="},
			{"Go版本", info.GoVersion},
			{"CPU核心数", fmt.Sprintf("%d", info.NumCPU)},
			{"操作系统", info.GOOS},
			{"系统架构", info.GOARCH},
		}
		allRows = append(allRows, sysRows...)
	}

	for _, row := range allRows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// writePerformanceEvaluation 写入性能评估
func (e *CSVExporter) writePerformanceEvaluation(writer *csv.Writer, summary map[string]interface{}) error {
	// 写入空行
	if err := writer.Write([]string{}); err != nil {
		return err
	}

	// 写入评估标题
	if err := writer.Write([]string{"=== 性能评估 ==="}); err != nil {
		return err
	}

	// 提取关键指标
	avgFPS, hasFPS := summary["avg_fps"].(float64)
	avgMemory, hasMemory := summary["avg_memory_mb"].(float64)
	avgCPU, hasCPU := summary["avg_cpu_percent"].(float64)

	// 初始化评估结果
	var evaluation string
	var suggestions []string

	if hasFPS && hasMemory && hasCPU {
		// 基于真实数据的评估
		if avgFPS >= 55 && avgMemory < 10 && avgCPU < 20 {
			evaluation = "优秀 - 帧率高，内存和CPU占用低"
		} else if avgFPS >= 45 && avgMemory < 20 && avgCPU < 30 {
			evaluation = "良好 - 性能表现均衡"
		} else if avgFPS >= 30 && avgMemory < 30 && avgCPU < 40 {
			evaluation = "一般 - 可接受性能"
			if avgFPS < 45 {
				suggestions = append(suggestions, "考虑优化渲染性能")
			}
			if avgMemory >= 20 {
				suggestions = append(suggestions, "注意内存使用")
			}
			if avgCPU >= 30 {
				suggestions = append(suggestions, "优化CPU使用")
			}
		} else {
			evaluation = "需优化 - 性能开销较大"
			if avgFPS < 30 {
				suggestions = append(suggestions, "FPS过低，需重点优化渲染")
			}
			if avgMemory >= 30 {
				suggestions = append(suggestions, "内存占用过高")
			}
			if avgCPU >= 40 {
				suggestions = append(suggestions, "CPU使用率过高")
			}
		}
	} else {
		evaluation = "数据不足，无法评估"
	}

	// 写入评估结果
	if err := writer.Write([]string{"性能评估", evaluation}); err != nil {
		return err
	}

	// 写入具体指标
	if hasFPS {
		if err := writer.Write([]string{"FPS评估", getFPSEvaluation(avgFPS)}); err != nil {
			return err
		}
	}
	if hasMemory {
		if err := writer.Write([]string{"内存评估", getMemoryEvaluation(avgMemory)}); err != nil {
			return err
		}
	}
	if hasCPU {
		if err := writer.Write([]string{"CPU评估", getCPUEvaluation(avgCPU)}); err != nil {
			return err
		}
	}

	// 写入优化建议
	if len(suggestions) > 0 {
		if err := writer.Write([]string{}); err != nil {
			return err
		}
		if err := writer.Write([]string{"优化建议"}); err != nil {
			return err
		}
		for _, suggestion := range suggestions {
			if err := writer.Write([]string{"•", suggestion}); err != nil {
				return err
			}
		}
	}

	return nil
}

// 辅助评估函数
func getFPSEvaluation(fps float64) string {
	if fps >= 55 {
		return fmt.Sprintf("优秀 (%.1f FPS)", fps)
	} else if fps >= 45 {
		return fmt.Sprintf("良好 (%.1f FPS)", fps)
	} else if fps >= 30 {
		return fmt.Sprintf("一般 (%.1f FPS)", fps)
	} else if fps >= 15 {
		return fmt.Sprintf("较差 (%.1f FPS)", fps)
	} else {
		return fmt.Sprintf("很差 (%.1f FPS)", fps)
	}
}

func getMemoryEvaluation(memory float64) string {
	if memory < 10 {
		return fmt.Sprintf("优秀 (%.2f MB)", memory)
	} else if memory < 20 {
		return fmt.Sprintf("良好 (%.2f MB)", memory)
	} else if memory < 30 {
		return fmt.Sprintf("一般 (%.2f MB)", memory)
	} else if memory < 50 {
		return fmt.Sprintf("较高 (%.2f MB)", memory)
	} else {
		return fmt.Sprintf("过高 (%.2f MB)", memory)
	}
}

func getCPUEvaluation(cpu float64) string {
	if cpu < 10 {
		return fmt.Sprintf("优秀 (%.1f%%)", cpu)
	} else if cpu < 20 {
		return fmt.Sprintf("良好 (%.1f%%)", cpu)
	} else if cpu < 30 {
		return fmt.Sprintf("一般 (%.1f%%)", cpu)
	} else if cpu < 40 {
		return fmt.Sprintf("较高 (%.1f%%)", cpu)
	} else {
		return fmt.Sprintf("过高 (%.1f%%)", cpu)
	}
}

// ExportComparison 专门导出对比结果
func (e *CSVExporter) ExportComparison(comparison *ComponentComparison) error {
	if comparison == nil {
		return fmt.Errorf("comparison data is nil")
	}

	// 设置对比专用的文件名
	e.filename = fmt.Sprintf("comparison_%s.csv", time.Now().Format("20060102_150405"))

	// 合并所有指标
	allMetrics := append(comparison.CustomMetrics, comparison.NativeMetrics...)

	// 创建摘要
	summary := map[string]interface{}{
		"test_name":      "ComponentComparison",
		"custom_summary": comparison.CustomSummary,
		"native_summary": comparison.NativeSummary,
		"comparison":     comparison.Comparison,
		"conclusion":     comparison.Conclusion,
		"system_info":    GetSystemInfo(),
	}

	return e.ExportMetrics(allMetrics, summary)
}

// ExportSummaryOnly 仅导出摘要信息
func (e *CSVExporter) ExportSummaryOnly(summary map[string]interface{}) error {
	// 修改文件名为summary
	e.filename = fmt.Sprintf("summary_%s.csv", time.Now().Format("20060102_150405"))

	filePath := filepath.Join(e.outputDir, e.filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create summary CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return e.writeSummary(writer, summary)
}

// SetFilename 设置输出文件名
func (e *CSVExporter) SetFilename(filename string) {
	e.filename = filename
	if filepath.Ext(e.filename) != ".csv" {
		e.filename += ".csv"
	}
}

// GetFullPath 获取完整文件路径
func (e *CSVExporter) GetFullPath() string {
	return filepath.Join(e.outputDir, e.filename)
}
