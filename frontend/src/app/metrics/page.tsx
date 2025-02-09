"use client";

import HistoryGraph from "@/app/components/HistoryGraph";

export default function MetricsPage() {
  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">System Metrics</h1>
      <HistoryGraph />
    </div>
  );
}
