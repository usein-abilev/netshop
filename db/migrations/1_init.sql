-- migrate:up
create table categories (
    id serial primary key,
    name varchar(255) not null,
    description text,
    unique(name)
);
create index categories_name_idx on categories(name);
create table sizes (
    id serial primary key,
    name varchar(32) not null,
    unique(name)
);
create index sizes_name_idx on sizes(name);
create table products (
    id serial primary key,
    name varchar(255) not null,
    description text,
    price decimal(10, 2) not null,
    size_id integer references sizes(id) on delete
    set null,
        category_id integer references categories(id) on delete
    set null
);
create index products_name_idx on products(name);
create index products_price_idx on products(price);
create table customers (
    id serial primary key,
    phone varchar(15) not null,
    email varchar(255) not null,
    metadata jsonb,
    unique(phone),
    unique(email)
);
create index customers_phone_idx on customers(phone);
create index customers_email_idx on customers(email);
create table orders (
    id serial primary key,
    order_date timestamp not null,
    customer_id integer references customers(id) on delete cascade,
    total decimal(10, 2) not null
);
create index orders_order_date_idx on orders(order_date);
create index orders_total_idx on orders(total);
create index orders_customer_id_idx on orders(customer_id);
create table roles (
    id serial primary key,
    name varchar(255) not null,
    permissions jsonb,
    unique(name)
);
create index roles_name_idx on roles(name);
create table employees (
    id serial primary key,
    username varchar(255) not null,
    password varchar(255) not null,
    role_id integer references roles(id) on delete
    set null,
        unique(username)
);
-- migrate:down
drop table if exists categories;
drop table if exists products;
drop table if exists customers;
drop table if exists orders;
drop table if exists roles;
drop table if exists employees;