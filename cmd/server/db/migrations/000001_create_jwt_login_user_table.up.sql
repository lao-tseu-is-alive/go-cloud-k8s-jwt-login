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
comment on table go_user is 'go_user is the main table of the GO_JWT_LOGIN microservice';



