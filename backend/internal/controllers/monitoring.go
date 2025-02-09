package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YL-Tan/GoHomeAi/internal/db"
	"github.com/YL-Tan/GoHomeAi/internal/httpresponse"
	"github.com/YL-Tan/GoHomeAi/internal/logger"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
)

const (
	LastRecordNum = 50
)

type SystemMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsed  uint64  `json:"memory_used"`
	MemoryTotal uint64  `json:"memory_total"`
	LoadAvg     float64 `json:"load_avg"`
	DiskUsage   uint64  `json:"disk_used"`
	DiskTotal   uint64  `json:"disk_total"`
}

type MonitoringController struct {
	Queries *db.Queries
}

func NewMonitoringController(q *db.Queries) *MonitoringController {
	return &MonitoringController{Queries: q}

}

func (c *MonitoringController) GetSysMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := GetSystemMetrics()
	if err != nil {
		httpresponse.SendJSON(w, http.StatusInternalServerError, false, "Failed to fetch system metrics", nil, err)
		return
	}
	httpresponse.SendJSON(w, http.StatusOK, true, "System metrics fetched successfully", metrics, err)
}

func (c *MonitoringController) GetHistSysMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := c.Queries.GetRecentMetrics(r.Context(), LastRecordNum)
	if err != nil {
		httpresponse.SendJSON(w, http.StatusInternalServerError, false, "Failed to fetch recent system metrics", nil, err)
		return
	}
	httpresponse.SendJSON(w, http.StatusOK, true, "Historical system metrics fetched successfully", metrics, err)
}

func GetSystemMetrics() (*db.InsertSystemMetricsParams, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		logger.Log.Error("Failed to get CPU usage", zap.Error(err))
		return nil, fmt.Errorf("failed to get CPU usage: %w", err)
	}
	memStats, err := mem.VirtualMemory()
	if err != nil {
		logger.Log.Error("Failed to get memory usage", zap.Error(err))
		return nil, fmt.Errorf("failed to get memory usage: %w", err)
	}
	loadStats, err := load.Avg()
	if err != nil {
		logger.Log.Error("Failed to get load average", zap.Error(err))
		return nil, fmt.Errorf("failed to get load average: %w", err)
	}
	diskStats, err := disk.Usage("/")
	if err != nil {
		logger.Log.Error("Failed to get disk usage", zap.Error(err))
		return nil, fmt.Errorf("failed to get disk usage: %w", err)
	}

	metrics := &SystemMetrics{
		CPUUsage:    cpuPercent[0],
		MemoryUsed:  memStats.Used,
		MemoryTotal: memStats.Total,
		LoadAvg:     loadStats.Load1,
		DiskUsage:   diskStats.Used,
		DiskTotal:   diskStats.Total,
	}

	logger.Log.Info("System Metrics Collected",
		zap.Float64("CPU Usage", metrics.CPUUsage),
		zap.Uint64("Memory Used", metrics.MemoryUsed),
		zap.Uint64("Memory Total", metrics.MemoryTotal),
		zap.Float64("Load Avg", metrics.LoadAvg),
		zap.Uint64("Disk Used", metrics.DiskUsage),
		zap.Uint64("Disk Total", metrics.DiskTotal),
	)

	return &db.InsertSystemMetricsParams{
		CpuUsage:    metrics.CPUUsage,
		MemoryUsed:  int64(metrics.MemoryUsed),
		MemoryTotal: int64(metrics.MemoryTotal),
		LoadAvg:     metrics.LoadAvg,
		DiskUsed:    int64(metrics.DiskUsage),
		DiskTotal:   int64(metrics.DiskTotal),
	}, nil
}

func PrintMetricsEvery(interval time.Duration) {
	for {
		metrics, err := GetSystemMetrics()
		if err != nil {
			fmt.Println("Error getting system metrics:", err)
			continue
		}
		metricsJSON, _ := json.MarshalIndent(metrics, "", "  ")
		fmt.Println("System Metrics:", string(metricsJSON))

		time.Sleep(interval)
	}
}

func StoreMetrics(ctx context.Context, queries *db.Queries) {
	for {
		time.Sleep(5 * time.Second)
		metrics, err := GetSystemMetrics()
		if err != nil {
			logger.Log.Error("Failed to collect system metrics", zap.Error(err))
			continue
		}
		err = queries.InsertSystemMetrics(ctx, *metrics)
		if err != nil {
			logger.Log.Error("Failed to store system metrics", zap.Error(err))
		}
	}
}
