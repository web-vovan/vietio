-- +goose Up
-- +goose StatementBegin
INSERT INTO cities ("name_vn", "name_rus")
VALUES ('Nha Trang', 'Нячанг')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE cities RESTART IDENTITY CASCADE;
-- +goose StatementEnd
