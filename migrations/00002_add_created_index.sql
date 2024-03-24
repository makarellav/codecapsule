-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_snippets_created ON snippets(created);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_snippets_created ON snippets;
-- +goose StatementEnd
