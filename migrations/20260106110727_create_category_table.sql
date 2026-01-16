-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
  id bigserial NOT NULL,
  "name" varchar(255) NOT NULL,
  "order" int4 NOT NULL,
  CONSTRAINT categories_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
