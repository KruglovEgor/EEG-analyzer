package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// MemoryStats holds memory statistics
type MemoryStats struct {
	AllocBefore      uint64
	AllocAfter       uint64
	AllocDiff        uint64
	TotalAllocBefore uint64
	TotalAllocAfter  uint64
	TotalAllocDiff   uint64
	SysBefore        uint64
	SysAfter         uint64
	NumGCBefore      uint32
	NumGCAfter       uint32
}

// MemoryMonitorMiddleware logs memory usage for each request
func MemoryMonitorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Force GC to get more accurate measurements (optional, can impact performance)
		// runtime.GC()

		var memBefore runtime.MemStats
		runtime.ReadMemStats(&memBefore)

		start := time.Now()

		// Process request
		c.Next()

		latency := time.Since(start)

		var memAfter runtime.MemStats
		runtime.ReadMemStats(&memAfter)

		stats := MemoryStats{
			AllocBefore:      memBefore.Alloc,
			AllocAfter:       memAfter.Alloc,
			AllocDiff:        memAfter.Alloc - memBefore.Alloc,
			TotalAllocBefore: memBefore.TotalAlloc,
			TotalAllocAfter:  memAfter.TotalAlloc,
			TotalAllocDiff:   memAfter.TotalAlloc - memBefore.TotalAlloc,
			SysBefore:        memBefore.Sys,
			SysAfter:         memAfter.Sys,
			NumGCBefore:      memBefore.NumGC,
			NumGCAfter:       memAfter.NumGC,
		}

		// Format memory values
		allocDiff := formatBytes(stats.AllocDiff)
		totalAllocDiff := formatBytes(stats.TotalAllocDiff)
		heapAlloc := formatBytes(memAfter.Alloc)
		heapSys := formatBytes(memAfter.Sys)

		// Log with color based on memory usage
		var color string
		if stats.AllocDiff > 100*1024*1024 { // > 100 MB
			color = "\033[31m" // Red
		} else if stats.AllocDiff > 10*1024*1024 { // > 10 MB
			color = "\033[33m" // Yellow
		} else {
			color = "\033[32m" // Green
		}
		reset := "\033[0m"

		gcInfo := ""
		if stats.NumGCAfter > stats.NumGCBefore {
			gcInfo = fmt.Sprintf(" [GC: %d times]", stats.NumGCAfter-stats.NumGCBefore)
		}

		fmt.Printf("%s[MEMORY]%s %s | %13v | Heap: %s | Sys: %s | Alloc Δ: %s%s | TotalAlloc Δ: %s%s | %s %s\n",
			"\033[36m", reset, // Cyan for [MEMORY]
			c.Request.Method,
			latency,
			heapAlloc,
			heapSys,
			color, allocDiff,
			reset, totalAllocDiff,
			c.Request.URL.Path,
			gcInfo,
		)
	}
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
