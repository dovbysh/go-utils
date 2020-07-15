-- +migrate Up

CREATE TABLE gender
(
    id              bigint NOT NULL,
    name            text,
    primary key (id)
);

-- +migrate Down
drop table gender;
