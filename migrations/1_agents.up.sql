create table if not exists users (
    id SERIAL primary key,
    name varchar not null,
    email varchar unique not null,
    phone_number varchar not null,
    hashed_password bytea not null,
    user_role varchar not null,
    create_time timestamptz
);

create table if not exists jobs(
    id varchar primary key unique not null,
    short_id varchar not null,
    create_time timestamptz not null,
    last_time_modified timestamptz not null,
    tracking_url varchar not null,
    creator BIGINT references users(id) not null,
    worker BIGINT references users(id),
    status int not null,
    order_status int not null,
    origin_company_name varchar,
    origin_first_name varchar,
    origin_second_name varchar,
    origin_phone_number varchar not null,
    origin_email_address varchar,
    origin_first_line_address varchar,
    origin_second_line_address varchar,
    origin_third_line_address varchar,
    origin_town varchar,
    origin_city varchar,
    origin_postcode varchar not null,
    origin_latitude float8,
    origin_longitude float8,
    origin_notes varchar,
    destination_company_name varchar,
    destination_first_name varchar,
    destination_second_name varchar,
    destination_phone_number varchar not null,
    destination_email_address varchar,
    destination_first_line_address varchar,
    destination_second_line_address varchar,
    destination_third_line_address varchar,
    destination_town varchar,
    destination_city varchar,
    destination_postcode varchar not null,
    destination_latitude float8,
    destination_longitude float8,
    destination_notes varchar,
    worker_notes varchar
);

create table if not exists orders(
    id varchar primary key,
    requester_user_id BIGINT references users(id) not null,
    status varchar not null
);

create table if not exists driver_states(
    courier_id BIGINT primary key references users(id) not null,
    active bool not null
);

create table if not exists driver_locations(
    courier_id BIGINT primary key references users(id) not null,
    coordinates varchar not null
);

create table if not exists api_keys(
    api_key_owner BIGINT primary key references users(id) not null,
    api_key bytea,
    last_modified_time timestamptz
);