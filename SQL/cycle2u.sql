drop table if exists users;
create table users (
	id serial not null primary key,
	username text not null default '',
	password text not null default '',
	email text not null default '',
	role text not null default 'public'
);
create unique index users_idx on users (username)