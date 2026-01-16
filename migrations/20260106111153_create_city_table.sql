-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cities (
  id bigserial NOT NULL,
  name_vn varchar(255) NOT NULL,
  name_rus varchar(255) NOT NULL,
  CONSTRAINT city_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cities;
-- +goose StatementEnd
