-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_metrics (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    cpu_usage FLOAT NOT NULL,
    memory_used BIGINT NOT NULL,
    memory_total BIGINT NOT NULL,
    load_avg FLOAT NOT NULL,
    disk_used BIGINT NOT NULL,
    disk_total BIGINT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE system_metrics;
-- +goose StatementEnd
