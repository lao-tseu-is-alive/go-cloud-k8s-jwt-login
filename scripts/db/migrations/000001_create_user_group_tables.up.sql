CREATE TYPE  orgunit_type AS ENUM ('Entreprise', 'Direction', 'Service', 'Office', 'Bureau', 'Unité', 'Division');
CREATE TABLE IF NOT EXISTS go_user
(
    id        	serial    CONSTRAINT go_user_pk   primary key,
    name			text	not null	constraint go_user_unique_name	unique
        constraint name_min_length check (length(btrim(name)) > 2),
    email			text	not null	constraint go_user_unique_email unique
        constraint email_min_length	check (length(btrim(email)) > 3),
    username		text	not null	constraint go_user_unique_username unique
        constraint username_min_length check (length(btrim(username)) > 2),
    password_hash	text	not null 	constraint password_hash_min_length check (length(btrim(password_hash)) > 30),
    external_id		int,
    orgunit_id		int,
    groups_id		int [],
    phone			text,
    is_locked		boolean   default false not null,
    is_admin		boolean   default false not null,
    create_time		timestamp default now() not null,
    creator			integer	not null,
    last_modification_time	timestamp,
    last_modification_user	integer,
    is_active				boolean default true not null,
    inactivation_time		timestamp,
    inactivation_reason    	text,
    comment                	text,
    bad_password_count     	integer default 0 not null
);
comment on table go_user is 'go_user is the main table of the GO_USER microservice';

CREATE TABLE IF NOT EXISTS go_group
(
    id        	serial    CONSTRAINT go_group_pk   primary key,
    name			text	not null	constraint go_group_unique_name	unique
        constraint name_min_length check (length(btrim(name)) > 2),
    create_time		timestamp default now() not null,
    creator			integer	not null,
    last_modification_time	timestamp,
    last_modification_user	integer,
    is_active				boolean default true not null,
    inactivation_time		timestamp,
    inactivation_reason    	text,
    comment                	text
);
comment on table go_group is 'go_group contains the list of groups of the GO_USER microservice';

create table if not exists public.go_orgunit
(
    id                     integer              not null
        constraint pk_go_orgunit
            primary key,
    type                   orgunit_type         not null,
    name                   text                 not null
        constraint go_orgunit_unique_name unique
        constraint name_min_length check (length(btrim(name)) > 2),
    parent_id              integer,
    abbreviation           text                 not null
        constraint go_orgunit_unique_abbreviation unique
        constraint abbreviation_min_length check (length(btrim(abbreviation)) > 1),
    description            text,
    order_list             text                 not null
        constraint go_orgunit_unique_order_list unique
        constraint order_list_min_length check (length(btrim(order_list)) > 1),
    phone                  text,
    email                  text,
    create_time            timestamp            not null,
    creator                integer              not null,
    last_modification_time timestamp,
    last_modification_user integer,
    is_active              boolean default true not null,
    inactivation_time      timestamp,
    inactivation_reason    text,
    comment                text,
    guid                   uuid,
    full_name_de_norm      text
);

alter table public.go_orgunit   owner to go_cloud_k8s_user_group;
comment on table go_orgunit is 'go_orgunit contains the list of organisational units of the GO_USER microservice';

INSERT INTO public.go_group (name, creator, comment)
VALUES ('global_admin', 1, 'global administrators');

