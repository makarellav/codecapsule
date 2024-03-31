-- +goose Up
-- +goose StatementBegin
CREATE INDEX sessions_expiry_idx ON sessions(expiry);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX sessions_expiry_idx ON sessions;
-- +goose StatementEnd
