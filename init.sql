create table if not exists vocabulary
(
    id          INTEGER
        primary key autoincrement,
    count       INTEGER default 1,
    translation text not null,
    lang        varchar,
    source      varchar
);