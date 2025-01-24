-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD COLUMN notify_start TIMESTAMP;
ALTER TABLE events ADD COLUMN notified BOOLEAN;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'No down SQL query for this migration';
-- +goose StatementEnd
