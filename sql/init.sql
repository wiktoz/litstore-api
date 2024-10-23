drop database if exists litstore;

DO $$
DECLARE
r RECORD;
BEGIN
FOR r IN
SELECT table_schema, table_name
FROM information_schema.tables
WHERE table_type = 'BASE TABLE' AND table_schema NOT IN ('pg_catalog', 'information_schema')
    LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || r.table_schema || '.' || r.table_name || ' CASCADE;';
END LOOP;
END $$;

DO $$
DECLARE
r RECORD;
BEGIN
FOR r IN
SELECT sequence_schema, sequence_name
FROM information_schema.sequences
         LOOP
    EXECUTE 'DROP SEQUENCE IF EXISTS ' || r.sequence_schema || '.' || r.sequence_name || ' CASCADE;';
END LOOP;
END $$;

DO $$
DECLARE
enum_record RECORD;
BEGIN
FOR enum_record IN
SELECT enumtypid::regtype AS enum_name
FROM pg_enum
GROUP BY enumtypid
    LOOP
        EXECUTE 'DROP TYPE IF EXISTS ' || enum_record.enum_name || ' CASCADE;';
END LOOP;
END $$;

create database litstore;

\connect litstore;

--- Enums
CREATE TYPE languages AS ENUM ('pl', 'en', 'fr', 'de');
CREATE TYPE currencies AS ENUM ('PLN', 'USD', 'EUR', 'CHF');
CREATE TYPE selection_types as ENUM ('select', 'button');
CREATE TYPE discount_unit AS ENUM ('percent', 'cash');
CREATE TYPE item_unit AS ENUM ('pcs.', 'set', 'kg', 'l');
CREATE TYPE order_status AS ENUM ('pending', 'processing', 'shipped', 'delivered', 'cancelled');
CREATE TYPE payment_status AS ENUM ('PENDING', 'SUCCESS', 'CANCELLED');

--- Sequences for auto ID generation
create sequence user_id_seq start 1 increment 1 cache 1;
create sequence role_id_seq start 1 increment 1 cache 1;
create sequence permission_id_seq start 1 increment 1 cache 1;
create sequence address_id_seq start 1 increment 1 cache 1;
create sequence category_id_seq start 1 increment 1 cache 1;
create sequence subcategory_id_seq start 1 increment 1 cache 1;
create sequence promo_code_id_seq start 1 increment 1 cache 1;
create sequence product_id_seq start 1 increment 1 cache 1;
create sequence product_photo_id_seq start 1 increment 1 cache 1;
create sequence product_description_id_seq start 1 increment 1 cache 1;
create sequence variant_type_id_seq start 1 increment 1 cache 1;
create sequence variant_option_id_seq start 1 increment 1 cache 1;
create sequence item_id_seq start 1 increment 1 cache 1;
create sequence delivery_id_seq start 1 increment 1 cache 1;
create sequence order_id_seq start 1000 increment 1 cache 1;
create sequence token_id_seq start 1 increment 1 cache 1;


create table Roles(
                      role_id int primary key unique not null default nextval('role_id_seq'),
                      name varchar(60) not null
);

create table Addresses(
                          address_id int primary key unique not null default nextval('address_id_seq'),
                          name varchar(60) not null,
                          surname varchar(60) not null,
                          street varchar(60) not null,
                          house varchar(20) not null,
                          flat varchar(20),
                          post_code varchar(10) not null,
                          city varchar(60) not null,
                          phone varchar(20) not null,
                          country varchar(60) not null
);

create table Users(
                      user_id int primary key unique not null default nextval('user_id_seq'),
                      email varchar(60) not null,
                      password char(64) not null,
                      Mainaddress_id int,
                      role_id int,
                      active boolean default true,
                      created_at timestamp default current_timestamp not null,
                      updated_at timestamp default current_timestamp not null,
                      FOREIGN KEY (role_id) REFERENCES Roles(role_id),
                      FOREIGN KEY (Mainaddress_id) REFERENCES Addresses(address_id)
);

insert into Roles(name) values ('admin');
insert into Users(email, password, role_id) values ('wiktoz05@icloud.com', 'aaa', 1);

create table Permissions(
                            permission_id int primary key unique not null default nextval('permission_id_seq'),
                            name varchar(60) not null
);

create table PermissionsRoles(
                                 role_id int not null,
                                 permission_id int not null,
                                 PRIMARY KEY (role_id, permission_id),
                                 FOREIGN KEY (role_id) REFERENCES Roles(role_id),
                                 FOREIGN KEY (permission_id) REFERENCES Permissions(permission_id)
);

create table PermissionsUsers(
                                 user_id int not null,
                                 permission_id int not null,
                                 PRIMARY KEY (user_id, permission_id),
                                 FOREIGN KEY (user_id) REFERENCES Users(user_id),
                                 FOREIGN KEY (permission_id) REFERENCES Permissions(permission_id)
);


create table Categories(
                           category_id int primary key unique not null default nextval('category_id_seq'),
                           name varchar(50) not null,
                           description text,
                           seo_description text,
                           img varchar(100),
                           bg_img varchar(100),
                           display_navbar boolean default true,
                           display_footer boolean default true,
                           active boolean default true,
                           slug varchar(60) not null
);

create table Subcategories(
                              subcategory_id int primary key unique not null default nextval('subcategory_id_seq'),
                              category_id int not null,
                              name varchar(50) not null,
                              description text,
                              seo_description text,
                              img varchar(100),
                              bg_img varchar(100),
                              display_navbar boolean default true,
                              display_footer boolean default true,
                              active boolean default true,
                              slug varchar(60) not null,
                              FOREIGN KEY (category_id) REFERENCES Categories(category_id)
);

create table Products(
                         product_id int primary key unique not null default nextval('product_id_seq'),
                         category_id int,
                         subcategory_id int,
                         name varchar(60),
                         manufacturer varchar(60),
                         new boolean not null default true,
                         active boolean not null default true,
                         slug varchar(70) not null,
                         created_at timestamp default current_timestamp not null,
                         updated_at timestamp default current_timestamp not null
);


insert into Categories(name, slug) values ('Test category', 'test-category');
insert into Subcategories(category_id, name, slug) values (1, 'Test subcategory', 'test-subcategory');
insert into Products(category_id, subcategory_id, name, manufacturer, new, active, slug) values (1, 1, 'Test product', 'Test manufacturer', true, true, 'test-product');

create table VariantTypes(
                             variant_type_id int primary key unique not null default nextval('variant_type_id_seq'),
                             name varchar(60) not null,
                             display_name varchar(40) not null,
                             selection_type selection_types not null,
                             slug varchar(70) not null
);

create table VariantOptions(
                               variant_option_id int primary key unique not null default nextval('variant_option_id_seq'),
                               variant_type_id int not null,
                               name varchar(40) not null,
                               order_index int not null,
                               FOREIGN KEY (variant_type_id) REFERENCES VariantTypes(variant_type_id)
);

create table PromoCodes(
                           promo_code_id int primary key unique not null default nextval('promo_code_id_seq'),
                           code varchar(50) not null,
                           discount numeric(6, 2) not null,
                           unit discount_unit not null,
                           name varchar(60) not null,
                           description text,
                           current_use int not null default 0,
                           max_use int,
                           start_date timestamp not null default current_timestamp,
                           end_date timestamp
);

create table PromoCodeCategories(
                                    promo_code_id int not null,
                                    category_id int not null,
                                    PRIMARY KEY (promo_code_id, category_id),
                                    FOREIGN KEY (promo_code_id) REFERENCES PromoCodes(promo_code_id),
                                    FOREIGN KEY (category_id) REFERENCES Categories(category_id)
);

create table PromoCodeSubcategories(
                                       promo_code_id int not null,
                                       subcategory_id int not null,
                                       PRIMARY KEY (promo_code_id, subcategory_id),
                                       FOREIGN KEY (promo_code_id) REFERENCES PromoCodes(promo_code_id),
                                       FOREIGN KEY (subcategory_id) REFERENCES Subcategories(subcategory_id)
);

create table PromoCodeProducts(
                                  promo_code_id int not null,
                                  product_id int not null,
                                  PRIMARY KEY (promo_code_id, product_id),
                                  FOREIGN KEY (promo_code_id) REFERENCES PromoCodes(promo_code_id),
                                  FOREIGN KEY (product_id) REFERENCES Products(product_id)
);

create table PromoCodeVariantOptions(
                                        promo_code_id int not null,
                                        variant_option_id int not null,
                                        PRIMARY KEY (promo_code_id, variant_option_id),
                                        FOREIGN KEY (promo_code_id) REFERENCES PromoCodes(promo_code_id),
                                        FOREIGN KEY (variant_option_id) REFERENCES VariantOptions(variant_option_id)
);

create table ProductPhotos(
                              product_photo_id int primary key unique not null default nextval('product_photo_id_seq'),
                              product_id int not null,
                              url varchar(150) not null,
                              order_index int not null,
                              FOREIGN KEY (product_id) REFERENCES Products(product_id)
);

create table ProductDescriptions(
                                    product_description_id int primary key unique not null default nextval('product_description_id_seq'),
                                    product_id int not null,
                                    lang_code languages not null,
                                    description text,
                                    FOREIGN KEY (product_id) REFERENCES Products(product_id)
);

create table Items(
                      item_id int primary key unique not null default nextval('item_id_seq'),
                      product_id int not null,
                      variant_option_id int not null,
                      price numeric(6,2) not null,
                      stock int not null,
                      unit item_unit not null default 'pcs.',
                      sku varchar(50) not null unique,
                      FOREIGN KEY (product_id) REFERENCES Products(product_id),
                      FOREIGN KEY (variant_option_id) REFERENCES VariantOptions(variant_option_id)
);

create table Deliveries(
                           delivery_id int primary key unique not null default nextval('delivery_id_seq'),
                           name varchar(60) not null,
                           description text,
                           img varchar(150) not null,
                           price numeric(6,2) not null,
                           free_from numeric(6,2) not null,
                           personal_collect boolean not null default false,
                           cash_on_delivery boolean not null default false,
                           active boolean not null default true,
                           slug varchar(70) not null
);

create table DeliveriesItems(
                                delivery_id int not null,
                                item_id int not null,
                                PRIMARY KEY (delivery_id, item_id),
                                FOREIGN KEY (delivery_id) REFERENCES Deliveries(delivery_id),
                                FOREIGN KEY (item_id) REFERENCES Items(item_id)
);

create table Orders(
                       order_id int primary key unique not null default nextval('order_id_seq'),
                       address_id int not null,
                       delivery_id int not null,
                       Deliveryaddress_id int not null,
                       promo_code_id int,
                       user_id int,
                       guest_email varchar(60),
                       discount_amount numeric(6,2) not null,
                       payment_amount numeric(6,2) not null,
                       payment_currency currencies not null,
                       payment_method varchar(50),
                       payment_hash varchar(64),
                       payment_url varchar(300),
                       payment_status payment_status not null default 'PENDING',
                       status order_status not null default 'pending',
                       created_at timestamp default current_timestamp not null,
                       updated_at timestamp default current_timestamp not null
);

create table OrderItems(
                           order_id int not null,
                           item_id int not null,
                           quantity int not null,
                           price numeric(6,2) not null,
                           PRIMARY KEY (order_id, item_id),
                           FOREIGN KEY (order_id) REFERENCES Orders(order_id),
                           FOREIGN KEY (item_id) REFERENCES Items(item_id)
);

create table Tokens(
                       token_id int primary key unique not null default nextval('token_id_seq'),
                       token char(64) not null,
                       created_at timestamp not null default current_timestamp,
                       expiring_at timestamp not null
);

--- Functions
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_order_index_photos() RETURNS TRIGGER AS $$
BEGIN
    NEW.order_index := (
        SELECT COUNT(*) + 1
        FROM ProductPhotos
        WHERE product_id = NEW.product_id
    );
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_order_index_variant_options() RETURNS TRIGGER AS $$
BEGIN
    NEW.order_index := (
        SELECT COUNT(*) + 1
        FROM VariantOptions
        WHERE variant_type_id = NEW.variant_type_id
    );
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--- Triggers for updating updated_at columns
CREATE TRIGGER update_updated_at_users
    BEFORE UPDATE ON Users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_updated_at_products
    BEFORE UPDATE ON Products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_updated_at_orders
    BEFORE UPDATE ON Orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_order_index_product_photos
    BEFORE UPDATE ON ProductPhotos
    FOR EACH ROW
    EXECUTE FUNCTION set_order_index_photos();

CREATE TRIGGER set_order_index_variant_options
    BEFORE UPDATE ON VariantOptions
    FOR EACH ROW
    EXECUTE FUNCTION set_order_index_variant_options();