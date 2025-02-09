"use client";

import { useState } from "react";
import useWebSocket from "react-use-websocket";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface SystemMetrics {
  cpu_usage: number;
  memory_used: number;
  memory_total: number;
  load_avg: number;
  disk_used: number;
  disk_total: number;
}

export default function SystemMetricsComponent() {
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);

  useWebSocket("ws://localhost:8080/ws", {
    onMessage: (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.type === "system_metrics" && data.metrics) {
          setMetrics(data.metrics);
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    },
    shouldReconnect: () => true,
  });

  return (
    <Card className="w-full max-w-md bg-gray-900 text-white shadow-md">
      <CardHeader>
        <CardTitle>System Metrics</CardTitle>
      </CardHeader>
      <CardContent>
        {metrics ? (
          <>
            <p>CPU Usage: {metrics.cpu_usage?.toFixed(2) ?? "N/A"}%</p>
            <p>Memory: {(metrics.memory_used / 1e9).toFixed(2) ?? "N/A"} GB / {(metrics.memory_total / 1e9).toFixed(2) ?? "N/A"} GB</p>
            <p>Load Average: {metrics.load_avg?.toFixed(2) ?? "N/A"}</p>
            <p>Disk: {(metrics.disk_used / 1e9).toFixed(2) ?? "N/A"} GB / {(metrics.disk_total / 1e9).toFixed(2) ?? "N/A"} GB</p>
          </>
        ) : (
          <p className="text-gray-500">Waiting for metrics...</p>
        )}
      </CardContent>
    </Card>
  );
}
