package server

import (
	"net/http"

	"github.com/YL-Tan/GoHomeAi/internal/controllers"
	"github.com/YL-Tan/GoHomeAi/internal/db"
	"github.com/YL-Tan/GoHomeAi/internal/workers"
)

func InitRouter(q *db.Queries, pool *workers.WorkerPool, wsServer *WebSocketServer) http.Handler {
	mux := http.NewServeMux()

	deviceController := controllers.NewDeviceController(q)
	jobController := controllers.NewJobController(pool)
	monitoringController := controllers.NewMonitoringController(q)

	mux.HandleFunc("/devices", deviceController.GetDevices)
	mux.HandleFunc("/api/jobs/status", jobController.GetJobStatus)
	mux.HandleFunc("/api/metrics", monitoringController.GetHistSysMetrics)

	mux.HandleFunc("/ws", wsServer.HandleWebSocket)

	return LoggingMiddleware(SecurityHeadersMiddleware(CORSMiddleware(mux)))
}
