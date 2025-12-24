// tools/benchmark/metrics.go
package benchmark

import (
	"runtime"
	"time"
)

// PerformanceMetric 性能指标数据结构
type PerformanceMetric struct {
	Timestamp     time.Time // 时间戳
	ComponentName string    // 组件名称（如：ParticleButton, MaterialEntry等）
	ComponentType string    // 组件类型（custom/native）
	TestScenario  string    // 测试场景（如：static_render, click_animation等）

	// 性能指标
	FPS           float64 // 帧率（Frames Per Second）
	MemoryUsageMB float64 // 内存使用量（MB）
	MemoryAllocMB float64 // 内存分配量（MB）
	NumGC         uint32  // GC次数
	GCTimeMS      float64 // GC耗时（ms）
	CPUPercent    float64 // CPU使用率（%）
	RenderTimeMS  float64 // 渲染耗时（ms）
	UpdateTimeMS  float64 // 更新耗时（ms）

	// 系统信息
	Goroutines int // Goroutine数量
	NumCores   int // CPU核心数
}

// SystemInfo 系统信息
type SystemInfo struct {
	GoVersion string
	NumCPU    int
	GOOS      string
	GOARCH    string
}

// TestConfig 测试配置
type TestConfig struct {
	SampleInterval  time.Duration // 采样间隔
	RecordingLength time.Duration // 记录时长
	TestName        string        // 测试名称
	OutputPath      string        // 输出路径
}

// NewPerformanceMetric 创建新的性能指标实例
func NewPerformanceMetric(componentName, componentType, scenario string) *PerformanceMetric {
	return &PerformanceMetric{
		Timestamp:     time.Now(),
		ComponentName: componentName,
		ComponentType: componentType,
		TestScenario:  scenario,
	}
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() SystemInfo {
	return SystemInfo{
		GoVersion: runtime.Version(),
		NumCPU:    runtime.NumCPU(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
	}
}
