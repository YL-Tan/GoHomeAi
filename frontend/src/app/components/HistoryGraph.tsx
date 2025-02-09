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
        const metrics = await response.json();
        setData(metrics);
      } catch (error) {
        console.error("Failed to fetch metrics:", error);
      }
    }
    fetchMetrics();
  }, []);

  return (
    <Card className="w-full bg-gray-900 text-white shadow-md p-4">
      <CardHeader>
        <CardTitle>CPU Usage Over Time</CardTitle>
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
