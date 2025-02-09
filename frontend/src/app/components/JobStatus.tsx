"use client";

import { useState } from "react";
import useWebSocket from "react-use-websocket";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface JobStatusData {
  active_jobs: number;
}

interface SystemMetrics {
  cpu_usage: number;
  memory_used: number;
  memory_total: number;
  load_avg: number;
  disk_used: number;
  disk_total: number;
}

export default function JobStatusClientSide() {
  const [status, setStatus] = useState<JobStatusData | null>(null);

  useWebSocket("ws://localhost:8080/ws", {
    onMessage: (event) => {
      const data = JSON.parse(event.data);
      
      if (data.type === "job_status") {
        setStatus(data);
      }
    },
    shouldReconnect: () => true,
  });

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <Card className="w-full max-w-sm bg-gray-900 text-white">
        <CardHeader>
          <CardTitle>Job Status</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-lg">Active Jobs: {status?.active_jobs ?? "N/A"}</p>
        </CardContent>
      </Card>
    </div>
  );
}
