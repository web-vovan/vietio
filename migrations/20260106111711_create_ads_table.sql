-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ads (
  id bigserial NOT NULL,
  user_id int8 NOT NULL,
  category_id int8 NOT NULL,
  city_id int8 NOT NULL,
  title varchar(255) NOT NULL,
  description text NOT NULL,
  price int4 NOT NULL,
  currency varchar(255) NOT NULL,
  district varchar(255) NULL,
  status int2 NOT NULL,
  expires_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT ads_pkey PRIMARY KEY (id),
  CONSTRAINT ads_category_id_foreign FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
  CONSTRAINT ads_city_id_foreign FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE SET NULL,
  CONSTRAINT ads_user_id_foreign FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ads;
-- +goose StatementEnd
