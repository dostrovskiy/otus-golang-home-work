-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    event_start TIMESTAMP NOT NULL,
    event_end TIMESTAMP NOT NULL,
    description VARCHAR(255),
    owner_id VARCHAR(255) NOT NULL,
    notify_before BIGINT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'No down SQL query for this migration';
-- +goose StatementEnd
