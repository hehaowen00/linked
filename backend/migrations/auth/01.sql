create table users (
    id text not null,
    secret text not null,
    email text not null,
    first text not null,
    last text not null,
    primary key (id)
    unique (secret),
    unique (email)
);
