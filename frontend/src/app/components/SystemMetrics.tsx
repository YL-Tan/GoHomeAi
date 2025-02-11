"use client";

import { useWebSocketData } from "@/context/WebSocketProvider";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function SystemMetricsComponent() {
  const { metrics } = useWebSocketData();

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