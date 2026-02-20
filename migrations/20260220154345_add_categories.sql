-- +goose Up
-- +goose StatementBegin
INSERT INTO categories ("name", "order")
VALUES
    ('Барахолка', 1),
    ('Байки', 2),
    ('Жильё', 3),
    ('Услуги', 4),
    ('Разное', 5)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE categories;
-- +goose StatementEnd
