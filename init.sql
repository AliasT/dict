create table if not exists vocabulary
(
    id          INTEGER
        primary key autoincrement,
    count       INTEGER default 1,
    translation varchar,
    lang        varchar,
    source      varchar
);