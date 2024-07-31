-- migrate:up

create table person (
    id serial primary key,
    first_name varchar(255) not null default '',
    last_name varchar(255) not null default '',
    phone varchar(15) not null,
    email varchar(255) not null,
    email_verified boolean not null default false,
    metadata jsonb,
    unique(phone),
    unique(email)
);
create index person_first_name_idx on person(first_name);
create index person_last_name_idx on person(last_name);
create index person_phone_idx on person(phone);
create index person_email_idx on person(email);

create table customers (
    id serial primary key,
    username varchar(255) null,
    password varchar(255) null,
    is_verified boolean not null default false,
    person_id integer references person(id) on delete cascade,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
create index customers_username_idx on customers(username);
create index customers_is_verified_idx on customers(is_verified);

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
    role_id integer references roles(id) on delete set null,
    person_id integer references person(id) on delete set null,
    unique(username)
);

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

create table colors (
    id serial primary key,
    name varchar(32) not null,
    unique(name)
);
create index colors_name_idx on colors(name);

create table files (
    id serial primary key,
    filename varchar(255) not null,
    filetype varchar(255) not null,
    size_bytes integer not null,
    path varchar(255) not null,
    width integer not null,
    height integer not null,
    created_at timestamp not null default now()
);

create table products (
    id serial primary key,
    name varchar(255) not null,
    description text,
    base_price decimal(10, 2) not null,
    category_id integer references categories(id) on delete set null,
    employee_id integer references employees(id) on delete set null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
create index products_name_idx on products(name);
create index products_price_idx on products(base_price);

create table product_variants (
    id serial primary key, 
    product_id integer references products(id) on delete cascade,
    size_id integer references sizes(id) on delete set null,
    color_id integer references colors(id) on delete set null,
    price decimal(10, 2) not null,
    stock integer not null,
    unique(product_id, size_id)
);
create index product_variants_product_id_idx on product_variants(product_id);
create index product_variants_size_id_idx on product_variants(size_id);
create index product_variants_color_id_idx on product_variants(color_id);
create index product_variants_price_idx on product_variants(price);
create index product_variants_stock_idx on product_variants(stock);
alter table product_variants add constraint check_stock_nonnegative check (stock >= 0);


create table product_variant_images (
    id serial primary key,
    file_id integer references files(id) on delete cascade,
    product_variant_id integer references product_variants(id) on delete cascade
);

create type order_status as enum('pending', 'processing', 'shipped', 'delivered');
create table orders (
    id serial primary key,
    order_date timestamp not null,
    customer_id integer references customers(id) on delete cascade,
    
    -- delivery info
    delivery_address varchar(255) not null,
    delivery_zipcode varchar(10) not null,
    delivery_city varchar(255) not null,
    delivery_country varchar(255) not null,

    status order_status not null default 'pending',
    status_date timestamp not null default now(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
create index orders_order_date_idx on orders(order_date);
create index orders_status_idx on orders(status);
create index orders_status_date_idx on orders(status_date);

create table order_items (
    id serial primary key,
    order_id integer references orders(id) on delete cascade,
    product_variant_id integer references product_variants(id) on delete set null,
    quantity integer not null,
    price decimal(10, 2) not null
);
create index order_items_order_id_idx on order_items(order_id);
create index order_items_product_variant_id_idx on order_items(product_variant_id);
create index order_items_quantity_idx on order_items(quantity);

insert into categories (name, description) values
    ('Shirts', 'All kinds of shirts'),
    ('Pants', 'All kinds of pants'),
    ('Shoes', 'All kinds of shoes');

insert into sizes (name) values
    ('XS'), 
    ('S'), 
    ('M'), 
    ('L'), 
    ('XL');

insert into colors (name) values
    ('Red'),
    ('Blue'),
    ('Green'),
    ('Yellow'),
    ('White'),
    ('Black');

insert into roles (name, permissions) values 
    ('admin', '{}'),
    ('employee', '{}'),
    ('customer', '{}');

insert into employees (username, password, role_id) values
    ('admin', '$argon2id$v=19$m=65536,t=3,p=2$ZuBBUFIG98Vm+53eTBRHaQ$nMKgOpEQdXdCcnkjDO4mc+0gS4BWzC/M/ZH4SOjUCJ8', 1);

-- migrate:down
drop table if exists categories;
drop table if exists products;
drop table if exists sizes;
drop table if exists colors;
drop table if exists product_variants;
drop table if exists customers;
drop table if exists orders;
drop table if exists roles;
drop table if exists employees;