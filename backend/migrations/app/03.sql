alter table collections add column archived boolean not null default false;
alter table collections add column archived_at int not null default 0;

alter table collections drop column deleted_at;
alter table items drop column deleted_at;
