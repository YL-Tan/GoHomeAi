"use client";

import { useWebSocketData } from "@/context/WebSocketProvider";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function JobStatusClientSide() {
  const { jobStatus } = useWebSocketData();

  return (
    <Card className="w-full max-w-sm bg-gray-900 text-white">
      <CardHeader>
        <CardTitle>Job Status</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-lg">Active Jobs: {jobStatus?.active_jobs ?? "N/A"}</p>
      </CardContent>
    </Card>
  );
}