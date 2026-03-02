-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wishlist (
  id bigserial NOT NULL,
  ad_uuid UUID NULL,
  user_id int8 NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT wishlist_pkey PRIMARY KEY (id),
  CONSTRAINT wishlist_user_ad_unique UNIQUE (user_id, ad_uuid)
);

ALTER TABLE wishlist ADD CONSTRAINT wishlist_user_id_foreign FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE wishlist ADD CONSTRAINT wishlist_ad_uuid_foreign FOREIGN KEY (ad_uuid) REFERENCES ads(uuid) ON DELETE CASCADE;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wishlist;
-- +goose StatementEnd
