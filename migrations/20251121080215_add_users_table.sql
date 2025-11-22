-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    user_id INTEGER NOT NULL,
    tg_user_id BIGINT NOT NULL,
    first_name VARCHAR,
    language_code VARCHAR,
    signup_date TIMESTAMP WITH TIME ZONE NOT NULL,
    language VARCHAR(3)
);

CREATE SEQUENCE user_user_id_seq
    AS INTEGER
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE user_user_id_seq OWNED BY users.user_id;

ALTER TABLE ONLY users ALTER COLUMN user_id SET DEFAULT nextval('user_user_id_seq');

ALTER TABLE ONLY users
    ADD CONSTRAINT user_tg_user_id_key UNIQUE (tg_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP SEQUENCE IF EXISTS user_user_id_seq;
-- +goose StatementEnd
