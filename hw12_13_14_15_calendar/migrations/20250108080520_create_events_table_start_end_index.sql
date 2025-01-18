-- +goose Up
-- +goose StatementBegin
create index owner_idx on events (owner_id);
create index event_start_end_idx on events using btree (event_start, event_end);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index owner_idx;
drop index event_start_end_idx;
-- +goose StatementEnd
