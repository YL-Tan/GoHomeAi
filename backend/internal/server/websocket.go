package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/YL-Tan/GoHomeAi/internal/controllers"
	"github.com/YL-Tan/GoHomeAi/internal/logger"
	"github.com/YL-Tan/GoHomeAi/internal/workers"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	Duration = 2 * time.Second
)

// Upgrade HTTP requests to WebSockets
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a WebSocket connection
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// WebSocketServer manages active connections
type WebSocketServer struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// NewWebSocketServer creates a WebSocket server
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Start runs the WebSocket server
func (s *WebSocketServer) Start(pool *workers.WorkerPool) {
	go func() {
		for {
			select {
			case client := <-s.register:
				s.mu.Lock()
				s.clients[client] = true
				s.mu.Unlock()
			case client := <-s.unregister:
				s.mu.Lock()
				if _, ok := s.clients[client]; ok {
					delete(s.clients, client)
					close(client.send)
				}
				s.mu.Unlock()
			case message := <-s.broadcast:
				s.mu.Lock()
				for client := range s.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(s.clients, client)
					}
				}
				s.mu.Unlock()
			}
		}
	}()

	// Periodically send job status updates
	go func() {
		for {
			time.Sleep(Duration)
			status := map[string]interface{}{
				"type":        "job_status",
				"active_jobs": pool.GetActiveJobs(),
			}
			msg, _ := json.Marshal(status)
			s.broadcast <- msg
		}
	}()

	// Periodically send system metrics
	go func() {
		for {
			time.Sleep(Duration)
			metrics, err := controllers.GetSystemMetrics()
			if err != nil {
				logger.Log.Error("Failed to fetch system metrics", zap.Error(err))
				continue
			}

			alert := ""
			if metrics.CpuUsage > 80 {
				alert = "⚠️ High CPU Usage: " + fmt.Sprintf("%.2f%%", metrics.CpuUsage)
				logger.Log.Warn(alert)
			}
			if float64(metrics.MemoryUsed)/float64(metrics.MemoryTotal) > 0.85 {
				alert = "⚠️ High Memory Usage!"
				logger.Log.Warn(alert)
			}

			status := map[string]interface{}{
				"type":         "system_metrics",
				"timestamp":    time.Now().Format(time.RFC3339),
				"cpu_usage":    metrics.CpuUsage,
				"memory_used":  metrics.MemoryUsed,
				"memory_total": metrics.MemoryTotal,
				"load_avg":     metrics.LoadAvg,
				"disk_used":    metrics.DiskUsed,
				"disk_total":   metrics.DiskTotal,
				"alert":        alert,
			}

			msg, _ := json.Marshal(status)
			s.broadcast <- msg
		}
	}()

}

// HandleWebSocket handles incoming WebSocket connections
func (s *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	s.register <- client

	// Start sending messages to client
	go client.writeMessages(s)
}

// writeMessages sends messages to the WebSocket client
func (c *Client) writeMessages(s *WebSocketServer) {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			s.unregister <- c
			break
		}
	}
}
