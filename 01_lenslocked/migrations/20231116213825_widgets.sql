-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS widgets (
	id SERIAL PRIMARY KEY,
	color TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS widgets;
-- +goose StatementEnd
