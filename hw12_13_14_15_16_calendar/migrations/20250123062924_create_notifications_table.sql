-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications (
    id VARCHAR(255) PRIMARY KEY,
    event_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    event_start TIMESTAMP NOT NULL,
    event_end TIMESTAMP NOT NULL,
    owner_id VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'No down SQL query for this migration';
-- +goose StatementEnd
