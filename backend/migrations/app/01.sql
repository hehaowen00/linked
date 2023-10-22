create table collections (
    id text unique not null,
    user_id text not null,
    name text not null,
    created_at int not null,
    deleted_at int not null default 0,
    primary key (id),
    unique (user_id, name)
);

create table items (
    id text not null,
    collection_id text not null,
    user_id text not null,
    url text not null,
    title text not null,
    description text not null,
    created_at int not null,
    deleted_at int not null default 0,
    primary key (id),
    unique (user_id, url)
);
