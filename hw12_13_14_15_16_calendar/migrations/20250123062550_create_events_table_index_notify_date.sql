-- +goose Up
-- +goose StatementBegin
create index event_notify_date_idx on events using btree (notify_start, event_start, notified);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index event_notify_date_idx;
-- +goose StatementEnd
