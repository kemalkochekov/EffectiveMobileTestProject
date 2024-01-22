-- +goose Up
-- +goose StatementBegin
CREATE TABLE individuals (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255) DEFAULT '',
    age INTEGER NOT NULL ,
    gender VARCHAR(10) DEFAULT '',
    country_id VARCHAR(2) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE individuals;
-- +goose StatementEnd
