-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders_plain (
  id serial PRIMARY KEY,
  order_json jsonb
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders_plain;
-- +goose StatementEnd
