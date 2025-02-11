"use client";

import { createContext, useContext, useState, useEffect } from "react";
import useWebSocket from "react-use-websocket";

interface SystemMetrics {
  timestamp: string;
  cpu_usage: number;
  memory_used: number;
  memory_total: number;
  load_avg: number;
  disk_used: number;
  disk_total: number;
  alert?: string; // Optional
}

interface JobStatus {
  active_jobs: number;
}

interface WebSocketContextType {
  metrics: SystemMetrics | null;
  historicalData: SystemMetrics[];
  jobStatus: JobStatus | null;
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
  const [historicalData, setHistoricalData] = useState<SystemMetrics[]>([]);
  const [jobStatus, setJobStatus] = useState<JobStatus | null>(null);

  useWebSocket("ws://localhost:8080/ws", {
    onMessage: (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.type === "system_metrics") {
          const formattedMetric: SystemMetrics = {
            timestamp: new Date(data.timestamp).toLocaleTimeString([], { 
              hour: "2-digit", 
              minute: "2-digit" 
            }),
            cpu_usage: data.cpu_usage,
            memory_used: data.memory_used,
            memory_total: data.memory_total,
            load_avg: data.load_avg,
            disk_used: data.disk_used,
            disk_total: data.disk_total,
          };
          setMetrics(formattedMetric);
          setHistoricalData((prev) => [...prev, formattedMetric].slice(-50));
        } else if (data.type === "job_status") {
          setJobStatus((prev) => {
            if (prev?.active_jobs !== data.active_jobs) {
              return { active_jobs: data.active_jobs };
            }
            return prev;
          });
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    },
    shouldReconnect: () => true,
  });

  return (
    <WebSocketContext.Provider value={{ metrics, historicalData, jobStatus }}>
      {children}
    </WebSocketContext.Provider>
  );
}

export function useWebSocketData() {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocketData must be used within a WebSocketProvider");
  }
  return context;
}
