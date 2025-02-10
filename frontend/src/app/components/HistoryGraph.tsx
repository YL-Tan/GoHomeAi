"use client";

import { useEffect, useState } from "react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Metric {
  timestamp: string;
  cpu_usage: number;
}

export default function HistoryGraph() {
  const [data, setData] = useState<Metric[]>([]);

  useEffect(() => {
    async function fetchMetrics() {
      try {
        const response = await fetch("http://localhost:8080/api/metrics");
        const jsonResponse = await response.json();
        const metrics = jsonResponse.data || [];
        console.log("Fetched Initial Metrics:", metrics);

        const formattedMetrics = metrics.map((item: Metric) => ({
          ...item,
          timestamp: new Date(item.timestamp).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
        }));

        setData(formattedMetrics);
      } catch (error) {
        console.error("Failed to fetch metrics:", error);
      }
    }

    fetchMetrics();

    // Establish WebSocket connection
    const ws = new WebSocket("ws://localhost:8080/ws"); // Adjust WebSocket URL if needed
    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.type === "system_metrics") {
        console.log("Received WebSocket Update:", message);

        // Format timestamp before updating the graph
        const newMetric = {
          timestamp: new Date(message.timestamp).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
          cpu_usage: message.cpu_usage,
        };

        setData((prevData) => {
          // Keep the last 50 records only (FIFO queue)
          const updatedData = [...prevData, newMetric].slice(-50);
          return updatedData;
        });
      }
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <Card className="w-full bg-gray-900 text-white shadow-md p-4">
      <CardHeader>
        <CardTitle>Live CPU Usage Over Time</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="timestamp" />
            <YAxis />
            <Tooltip />
            <Line type="monotone" dataKey="cpu_usage" stroke="#8884d8" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
