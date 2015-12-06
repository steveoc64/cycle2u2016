drop table if exists users;
create table users (
	id serial not null primary key,
	username text not null default '',
	password text not null default '',
	email text not null default '',
	role text not null default 'public'
);
create unique index users_idx on users (username);

drop table if exists booking;
create table booking (
	id serial not null primary key,
	name text not null default '',
	address text not null default '',
	bike text not null default '',
	email text not null default '',
	telephone text not null default '',
	enquiry text not null default '',
	message text not null default ''
);
create index booking_idx on booking (email);

drop table if exists customer;
create table customer (
	id serial not null primary key,
	email text not null default ''
);
create unique index customer_idx on customer (email);