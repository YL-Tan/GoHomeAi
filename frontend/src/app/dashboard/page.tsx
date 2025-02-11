"use client";

import HistoryGraph from "@/app/components/HistoryGraph";
import SystemMetrics from "@/app/components/SystemMetrics";

export default function Dashboard() {
  return (
    <div className="p-6 grid grid-cols-1 md:grid-cols-2 gap-6">
      <SystemMetrics />
      <HistoryGraph />
    </div>
  );
}