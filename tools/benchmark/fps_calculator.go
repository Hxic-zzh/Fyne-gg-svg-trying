// tools/benchmark/fps_calculator.go
package benchmark

import (
	"sync"
	"time"
)

// FPSCalculator 帧率计算器
type FPSCalculator struct {
	mu         sync.RWMutex
	frameTimes []time.Time
	maxSamples int
	lastFPS    float64
	lastUpdate time.Time
}

// NewFPSCalculator 创建帧率计算器
func NewFPSCalculator(maxSamples int) *FPSCalculator {
	return &FPSCalculator{
		frameTimes: make([]time.Time, 0, maxSamples),
		maxSamples: maxSamples,
		lastFPS:    60.0, // 默认值
		lastUpdate: time.Now(),
	}
}

// AddFrame 记录一帧
func (f *FPSCalculator) AddFrame() {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	f.frameTimes = append(f.frameTimes, now)

	// 保持样本数量
	if len(f.frameTimes) > f.maxSamples {
		f.frameTimes = f.frameTimes[1:]
	}

	// 每秒更新一次FPS
	if now.Sub(f.lastUpdate) >= time.Second {
		f.calculateFPS()
		f.lastUpdate = now
	}
}

// calculateFPS 计算帧率
func (f *FPSCalculator) calculateFPS() {
	if len(f.frameTimes) < 2 {
		return
	}

	// 计算最后一秒内的帧数
	oneSecondAgo := time.Now().Add(-time.Second)
	count := 0
	for i := len(f.frameTimes) - 1; i >= 0; i-- {
		if f.frameTimes[i].After(oneSecondAgo) {
			count++
		} else {
			break
		}
	}

	if count > 0 {
		f.lastFPS = float64(count)
	}
}

// GetFPS 获取当前帧率
func (f *FPSCalculator) GetFPS() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 如果太久没更新，重新计算
	if time.Since(f.lastUpdate) > 2*time.Second {
		f.calculateFPS()
		f.lastUpdate = time.Now()
	}

	return f.lastFPS
}

// Reset 重置计算器
func (f *FPSCalculator) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.frameTimes = make([]time.Time, 0, f.maxSamples)
	f.lastFPS = 60.0
	f.lastUpdate = time.Now()
}
