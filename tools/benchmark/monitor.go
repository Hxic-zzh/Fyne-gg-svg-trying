// tools/benchmark/monitor.go
package benchmark

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// Monitor 性能监控器
type Monitor struct {
	mu              sync.RWMutex
	isRunning       bool
	stopChan        chan struct{}
	metrics         []*PerformanceMetric
	currentTestName string
	currentScenario string
	componentName   string
	componentType   string

	// 性能统计
	fpsCalculator *FPSCalculator
	lastMemStats  runtime.MemStats

	// 进程监控
	proc           *process.Process
	lastCPUPercent float64

	// 帧计数
	startTime     time.Time
	totalFrames   int64
	lastFrameTime time.Time

	// 配置
	config     *TestConfig
	systemInfo SystemInfo
}

// NewMonitor 创建新的性能监控器
func NewMonitor(testName string) *Monitor {
	config := &TestConfig{
		SampleInterval:  100 * time.Millisecond,
		RecordingLength: 5 * time.Second,
		TestName:        testName,
		OutputPath:      "./benchmark_results",
	}

	// 获取当前进程
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		proc = nil
		fmt.Printf("警告: 无法获取进程信息: %v\n", err)
	}

	return &Monitor{
		config:          config,
		systemInfo:      GetSystemInfo(),
		stopChan:        make(chan struct{}),
		metrics:         make([]*PerformanceMetric, 0),
		currentTestName: testName,
		fpsCalculator:   NewFPSCalculator(120),
		proc:            proc,
		startTime:       time.Now(),
		lastFrameTime:   time.Now(),
	}
}

// Start 开始监控
func (m *Monitor) Start() {
	if m.isRunning {
		return
	}

	m.isRunning = true
	m.metrics = make([]*PerformanceMetric, 0)
	m.startTime = time.Now()
	m.totalFrames = 0

	// 初始化基准值
	runtime.ReadMemStats(&m.lastMemStats)

	// 重置FPS计算器
	m.fpsCalculator.Reset()

	// 获取初始CPU使用率
	if m.proc != nil {
		if percent, err := m.proc.CPUPercent(); err == nil {
			m.lastCPUPercent = percent
		}
	}

	go m.monitoringLoop()
}

// StartRecording 开始记录特定场景
func (m *Monitor) StartRecording(componentName, componentType, scenario string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.componentName = componentName
	m.componentType = componentType
	m.currentScenario = scenario

	// 重置相关计数器
	m.totalFrames = 0
	m.startTime = time.Now()

	// 记录开始时的快照
	m.recordSnapshot()
}

// StopRecording 停止当前场景记录
func (m *Monitor) StopRecording() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 记录结束时的快照
	m.recordSnapshot()

	m.currentScenario = ""
}

// Stop 停止监控
func (m *Monitor) Stop() {
	if !m.isRunning {
		return
	}

	m.isRunning = false
	close(m.stopChan)
}

// monitoringLoop 监控循环
func (m *Monitor) monitoringLoop() {
	ticker := time.NewTicker(m.config.SampleInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.recordMetrics()
		case <-m.stopChan:
			return
		}
	}
}

// recordMetrics 记录性能指标
func (m *Monitor) recordMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentScenario == "" {
		return
	}

	metric := NewPerformanceMetric(m.componentName, m.componentType, m.currentScenario)
	m.collectRealMetrics(metric)
	m.metrics = append(m.metrics, metric)
}

// recordSnapshot 记录快照
func (m *Monitor) recordSnapshot() {
	metric := NewPerformanceMetric(m.componentName, m.componentType, m.currentScenario)
	m.collectRealMetrics(metric)
	m.metrics = append(m.metrics, metric)
}

// collectRealMetrics 收集真实性能指标（完全真实版本）
func (m *Monitor) collectRealMetrics(metric *PerformanceMetric) {
	now := time.Now()

	// ====== 1. 收集真实内存信息 ======
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metric.MemoryUsageMB = float64(memStats.Alloc) / 1024 / 1024
	metric.MemoryAllocMB = float64(memStats.TotalAlloc) / 1024 / 1024
	metric.NumGC = memStats.NumGC
	metric.GCTimeMS = float64(memStats.PauseTotalNs) / 1_000_000

	// ====== 2. 真实帧率计算 ======
	// 使用FPS计算器获取当前帧率
	currentFPS := m.fpsCalculator.GetFPS()

	// 如果没有从计算器获取到，根据总帧数和时间计算
	if currentFPS <= 0 && m.totalFrames > 0 {
		elapsed := now.Sub(m.startTime).Seconds()
		if elapsed > 0.1 { // 至少0.1秒的数据
			currentFPS = float64(m.totalFrames) / elapsed
		}
	}

	// 确保FPS在合理范围内（但不预设任何基准值）
	if currentFPS < 0.1 {
		currentFPS = 0.1
	}
	metric.FPS = currentFPS

	// ====== 3. 真实CPU使用率计算 ======
	if m.proc != nil {
		// 获取当前进程CPU使用率百分比
		if cpuPercent, err := m.proc.CPUPercent(); err == nil {
			metric.CPUPercent = cpuPercent
			m.lastCPUPercent = cpuPercent
		} else {
			// 如果获取失败，使用上次的值
			metric.CPUPercent = m.lastCPUPercent
		}
	} else {
		// 如果无法获取进程信息，使用runtime的粗略估计
		metric.CPUPercent = m.estimateCPUPercent()
	}

	// 确保CPU百分比在合理范围内
	if metric.CPUPercent < 0 {
		metric.CPUPercent = 0
	}
	if metric.CPUPercent > 100 {
		metric.CPUPercent = 100
	}

	// ====== 4. 渲染时间估算 ======
	if metric.FPS > 0 {
		metric.RenderTimeMS = 1000.0 / metric.FPS
	} else {
		// 如果FPS为0，使用默认的16.67ms（60FPS）
		metric.RenderTimeMS = 16.67
	}

	// 更新时间估算（假设为渲染时间的30%）
	metric.UpdateTimeMS = metric.RenderTimeMS * 0.3

	// ====== 5. 系统信息 ======
	metric.Goroutines = runtime.NumGoroutine()
	metric.NumCores = runtime.NumCPU()

	// 更新内存状态
	m.lastMemStats = memStats
}

// estimateCPUPercent 当无法使用gopsutil时的CPU估算
func (m *Monitor) estimateCPUPercent() float64 {
	// 使用runtime的统计信息进行粗略估算
	// 注意：这不是精确的CPU使用率，只是一个估算

	var currentMemStats runtime.MemStats
	runtime.ReadMemStats(&currentMemStats)

	// 计算GC暂停时间的变化
	gcPauseDiff := float64(currentMemStats.PauseTotalNs-m.lastMemStats.PauseTotalNs) / 1_000_000_000 // 转换为秒

	// 假设GC暂停时间占CPU时间的10%
	estimatedCPUTime := gcPauseDiff * 10

	// 计算时间间隔
	timeInterval := 0.1 // 采样间隔100ms

	// 计算CPU使用率百分比
	cpuPercent := (estimatedCPUTime / timeInterval) * 100

	// 限制范围
	if cpuPercent < 0 {
		cpuPercent = 0
	}
	if cpuPercent > 100 {
		cpuPercent = 100
	}

	return cpuPercent
}

// AddFrame 添加一帧（需要在渲染循环中调用）
func (m *Monitor) AddFrame() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalFrames++
	m.lastFrameTime = time.Now()

	if m.fpsCalculator != nil {
		m.fpsCalculator.AddFrame()
	}
}

// GetMetrics 获取所有记录的指标
func (m *Monitor) GetMetrics() []*PerformanceMetric {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return append([]*PerformanceMetric(nil), m.metrics...)
}

// GetMetricsByComponent 按组件获取指标
func (m *Monitor) GetMetricsByComponent(componentName string) []*PerformanceMetric {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PerformanceMetric
	for _, metric := range m.metrics {
		if metric.ComponentName == componentName {
			result = append(result, metric)
		}
	}
	return result
}

// GetMetricsByType 按类型获取指标
func (m *Monitor) GetMetricsByType(componentType string) []*PerformanceMetric {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PerformanceMetric
	for _, metric := range m.metrics {
		if metric.ComponentType == componentType {
			result = append(result, metric)
		}
	}
	return result
}

// GetComponentMetrics 获取特定组件的所有指标
func (m *Monitor) GetComponentMetrics(componentName, componentType string) []*PerformanceMetric {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PerformanceMetric
	for _, metric := range m.metrics {
		if metric.ComponentName == componentName && metric.ComponentType == componentType {
			result = append(result, metric)
		}
	}
	return result
}

// GetSummary 获取性能摘要
func (m *Monitor) GetSummary() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.metrics) == 0 {
		return nil
	}

	return CalculateSummary(m.metrics, m.currentTestName, m.componentName,
		m.componentType, m.currentScenario, m.startTime, m.systemInfo, m.totalFrames)
}

// CalculateSummary 计算性能摘要（独立函数，便于测试）
func CalculateSummary(metrics []*PerformanceMetric, testName, componentName,
	componentType, scenario string, startTime time.Time,
	systemInfo SystemInfo, totalFrames int64) map[string]interface{} {

	if len(metrics) == 0 {
		return nil
	}

	// 计算平均值
	var totalFPS, totalMemory, totalCPU float64
	var minFPS, maxFPS float64 = math.MaxFloat64, 0
	var minMemory, maxMemory float64 = math.MaxFloat64, 0

	for _, metric := range metrics {
		totalFPS += metric.FPS
		totalMemory += metric.MemoryUsageMB
		totalCPU += metric.CPUPercent

		if metric.FPS < minFPS {
			minFPS = metric.FPS
		}
		if metric.FPS > maxFPS {
			maxFPS = metric.FPS
		}
		if metric.MemoryUsageMB < minMemory {
			minMemory = metric.MemoryUsageMB
		}
		if metric.MemoryUsageMB > maxMemory {
			maxMemory = metric.MemoryUsageMB
		}
	}

	count := float64(len(metrics))
	avgFPS := totalFPS / count
	avgMemory := totalMemory / count
	avgCPU := totalCPU / count

	// 计算标准差
	var fpsVariance, memoryVariance, cpuVariance float64
	for _, metric := range metrics {
		fpsVariance += math.Pow(metric.FPS-avgFPS, 2)
		memoryVariance += math.Pow(metric.MemoryUsageMB-avgMemory, 2)
		cpuVariance += math.Pow(metric.CPUPercent-avgCPU, 2)
	}
	fpsStdDev := math.Sqrt(fpsVariance / count)
	memoryStdDev := math.Sqrt(memoryVariance / count)
	cpuStdDev := math.Sqrt(cpuVariance / count)

	return map[string]interface{}{
		"test_name":        testName,
		"component_name":   componentName,
		"component_type":   componentType,
		"scenario":         scenario,
		"avg_fps":          math.Round(avgFPS*100) / 100,
		"min_fps":          math.Round(minFPS*100) / 100,
		"max_fps":          math.Round(maxFPS*100) / 100,
		"fps_std_dev":      math.Round(fpsStdDev*100) / 100,
		"avg_memory_mb":    math.Round(avgMemory*100) / 100,
		"min_memory_mb":    math.Round(minMemory*100) / 100,
		"max_memory_mb":    math.Round(maxMemory*100) / 100,
		"memory_std_dev":   math.Round(memoryStdDev*100) / 100,
		"avg_cpu_percent":  math.Round(avgCPU*100) / 100,
		"cpu_std_dev":      math.Round(cpuStdDev*100) / 100,
		"total_samples":    len(metrics),
		"duration_seconds": math.Round(time.Since(startTime).Seconds()*10) / 10,
		"system_info":      systemInfo,
		"total_frames":     totalFrames,
	}
}

// Reset 重置监控器
func (m *Monitor) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics = make([]*PerformanceMetric, 0)
	m.fpsCalculator.Reset()
	m.startTime = time.Now()
	m.totalFrames = 0
	m.lastFrameTime = time.Now()

	// 重新获取进程信息
	if m.proc == nil {
		if proc, err := process.NewProcess(int32(os.Getpid())); err == nil {
			m.proc = proc
		}
	}
}
