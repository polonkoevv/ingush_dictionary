-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_to_dict (
    tg_user_id BIGINT NOT NULL,
    dict_id INTEGER NOT NULL
);

ALTER TABLE ONLY users_to_dict
    ADD CONSTRAINT users_to_dict_pkey PRIMARY KEY (tg_user_id, dict_id);

ALTER TABLE ONLY users_to_dict
    ADD CONSTRAINT dict_id_to_dict_id FOREIGN KEY (dict_id) REFERENCES dictionary(dict_id);

ALTER TABLE ONLY users_to_dict
    ADD CONSTRAINT tg_user_id_to_tg_user_id FOREIGN KEY (tg_user_id) REFERENCES users(tg_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_to_dict;
-- +goose StatementEnd
