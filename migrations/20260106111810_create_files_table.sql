-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists files (
  id bigserial NOT NULL,
  ad_id int8 NULL,
  "path" varchar(255) NOT NULL,
  "order" int2 NULL,
  "size" int NULL,
  "mime" varchar(255) NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT files_pkey PRIMARY KEY (id)
);

ALTER TABLE files ADD CONSTRAINT files_ad_id_foreign FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files;
-- +goose StatementEnd
