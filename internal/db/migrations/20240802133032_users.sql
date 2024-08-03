-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id       bigserial primary key,
    login    text not null,
    email    text not null,
    password text not null,

    UNIQUE (login, email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
