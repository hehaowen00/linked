alter table items rename to items_backup;

create table items (
    id text not null,
    user_id text not null,
    url text not null,
    title text not null,
    description text not null,
    created_at int not null,
    deleted_at int not null default 0,
    primary key (id),
    unique (user_id, url)
);

insert into items (id, user_id, url, title, description, created_at, deleted_at)
select id, user_id, url, title, description, created_at, deleted_at
from items_backup;

create table tags (
    id text not null,
    user_id text not null,
    name text not null,
    created_at int not null,
    deleted_at int not null default 0,
    primary key (id),
    unique (user_id, name)
);

create table item_collection_map (
    user_id text not null,
    collection_id text not null,
    item_id text not null,
    primary key (user_id, collection_id, item_id)
);

insert into item_collection_map (user_id, collection_id, item_id)
select user_id, collection_id, id
from items_backup;

drop table profiles;

drop table items_backup;
