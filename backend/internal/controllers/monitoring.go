package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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

type MetricResponse struct {
	Timestamp string  `json:"timestamp"`
	CPUUsage  float64 `json:"cpu_usage"`
}

type MonitoringController struct {
	Queries *db.Queries
}

func NewMonitoringController(q *db.Queries) *MonitoringController {
	return &MonitoringController{Queries: q}

}

func (c *MonitoringController) GetSysMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := c.Queries.GetRecentMetrics(r.Context(), 10)
	if err != nil {
		httpresponse.SendJSON(w, http.StatusInternalServerError, false, "Failed to fetch system metrics", nil, err)
		return
	}

	var formattedMetrics []MetricResponse
	for _, metric := range metrics {
		if metric.Timestamp.Valid {
			formattedMetrics = append(formattedMetrics, MetricResponse{
				Timestamp: metric.Timestamp.Time.Format(time.RFC3339),
				CPUUsage:  metric.CpuUsage,
			})
		}
	}
	sort.Slice(formattedMetrics, func(i, j int) bool {
		return formattedMetrics[i].Timestamp < formattedMetrics[j].Timestamp
	})
	httpresponse.SendJSON(w, http.StatusOK, true, "Historical system metrics fetched successfully", formattedMetrics, nil)
}

func (c *MonitoringController) GetHistSysMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := c.Queries.GetRecentMetrics(r.Context(), LastRecordNum)
	if err != nil {
		httpresponse.SendJSON(w, http.StatusInternalServerError, false, "Failed to fetch recent system metrics", nil, err)
		return
	}

	// Convert to proper JSON format
	var formattedMetrics []MetricResponse
	for _, metric := range metrics {
		if metric.Timestamp.Valid {
			formattedMetrics = append(formattedMetrics, MetricResponse{
				Timestamp: metric.Timestamp.Time.Format(time.RFC3339),
				CPUUsage:  metric.CpuUsage,
			})
		}
	}

	// Sort in ascending order
	sort.Slice(formattedMetrics, func(i, j int) bool {
		return formattedMetrics[i].Timestamp < formattedMetrics[j].Timestamp
	})

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

func StartMetricsCollection(ctx context.Context, queries *db.Queries) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics, err := GetSystemMetrics()
			if err != nil {
				logger.Log.Error("Failed to collect system metrics", zap.Error(err))
				continue
			}
			err = queries.InsertSystemMetrics(ctx, *metrics)
			if err != nil {
				logger.Log.Error("Failed to store system metrics", zap.Error(err))
			}
		case <-ctx.Done():
			logger.Log.Info("Stopping system metrics collection")
			return
		}
	}
}
