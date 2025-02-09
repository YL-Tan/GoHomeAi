-- name: InsertSystemMetrics :exec
INSERT INTO system_metrics (cpu_usage, memory_used, memory_total, load_avg, disk_used, disk_total)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetRecentMetrics :many
SELECT timestamp, cpu_usage, memory_used, memory_total, load_avg, disk_used, disk_total
FROM system_metrics
ORDER BY timestamp DESC
LIMIT $1;
