create table todos (
    id text not null,
    user_id text not null
    title text not null,
    desc text not null,
    status text not null,
    created_at int not null,
    completed_at int not null,
    deleted_at int not null,
);

create table statuses (
    id text not null
    user_id text not null,
    name text not null,
    primary key (id, user_id),
    unique (user_id, name)
);
