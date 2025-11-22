-- +goose Up
-- +goose StatementBegin
CREATE TABLE word (
    word_id INTEGER NOT NULL,
    word VARCHAR NOT NULL,
    speech_part VARCHAR,
    translation VARCHAR NOT NULL,
    topic VARCHAR,
    dict_id INTEGER NOT NULL
);

CREATE SEQUENCE word_word_id_seq
    AS INTEGER
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE word_word_id_seq OWNED BY word.word_id;

ALTER TABLE ONLY word ALTER COLUMN word_id SET DEFAULT nextval('word_word_id_seq');

ALTER TABLE ONLY word
    ADD CONSTRAINT word_dict_id FOREIGN KEY (dict_id) REFERENCES dictionary(dict_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS word;
DROP SEQUENCE IF EXISTS word_word_id_seq;
-- +goose StatementEnd
