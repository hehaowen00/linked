alter table collections drop column deleted_at;

alter table collections add column archived boolean not null default false;

alter table items drop column deleted_at;
