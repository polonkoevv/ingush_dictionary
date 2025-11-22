-- +goose Up
-- +goose StatementBegin
CREATE TABLE dictionary (
    dict_id SERIAL PRIMARY KEY,
    abbr VARCHAR(10) UNIQUE,
    name VARCHAR NOT NULL,
    author VARCHAR
);

-- Индекс для поиска по названию
CREATE INDEX idx_dictionary_name ON dictionary(name);

-- Индекс для поиска по автору
CREATE INDEX idx_dictionary_author ON dictionary(author);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dictionary;
-- +goose StatementEnd
