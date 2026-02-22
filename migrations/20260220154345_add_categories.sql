-- +goose Up
-- +goose StatementBegin
INSERT INTO categories ("name", "order")
VALUES
    ('Маркет', 1),
    ('Байки', 2),
    ('Жильё', 3),
    ('Услуги', 4),
    ('Работа', 4),
    ('Разное', 6)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE categories RESTART IDENTITY CASCADE;
-- +goose StatementEnd
