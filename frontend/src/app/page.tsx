"use client";

import JobStatus from "@/app/components/JobStatus";
import SystemMetrics from "@/app/components/SystemMetrics";
import HistoryGraph from "@/app/components/HistoryGraph";

export default function Dashboard() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-6">
      {/* Real-time Job Status */}
      <JobStatus />

      {/* Real-time System Metrics */}
      <SystemMetrics />

      {/* Historical Graphs */}
      <div className="col-span-1 md:col-span-2">
        <HistoryGraph />
      </div>
    </div>
  );
}
