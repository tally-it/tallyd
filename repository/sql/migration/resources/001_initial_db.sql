-- +migrate Up
CREATE TABLE users
(
  user_id    INT AUTO_INCREMENT
    PRIMARY KEY,
  name       VARCHAR(191)                       NOT NULL,
  email      VARCHAR(191)                       NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME                           NULL,
  is_blocked BIT DEFAULT b'0'                   NOT NULL,
  is_admin   BIT DEFAULT b'0'                   NOT NULL,
  CONSTRAINT users_name_uindex
  UNIQUE (name)
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

create table user_auths
(
  user_auth_id int auto_increment
    primary key,
  user_id      int          not null,
  method       varchar(255) not null,
  value        tinyblob     null,
  constraint fk_user_auths_user
  foreign key (user_id) references users (user_id)
    on update cascade
    on delete cascade
)
  engine = InnoDB
  charset = utf8mb4;

create index fk_user_auths_user
  on user_auths (user_id);

create table categories
(
  category_id int auto_increment
    primary key,
  name        varchar(255)     not null,
  is_visible  bit default b'1' not null,
  is_active   bit default b'0' not null,
  parent_id   int              null,
  constraint categories_categories_category_id_fk
  foreign key (parent_id) references categories (category_id)
    on update cascade
    on delete cascade
)
  engine = InnoDB
  charset = utf8mb4;

create index categories_categories_category_id_fk
  on categories (parent_id);

create table products
(
  product_id int auto_increment
    primary key
)
  engine = InnoDB;


create table product_versions
(
  product_version_id int auto_increment
    primary key,
  product_id         int                                not null,
  name               varchar(255)                       not null,
  GTIN               char(14)                           null,
  price              decimal(15, 2)                     not null,
  added_at           datetime default CURRENT_TIMESTAMP not null,
  deleted_at         datetime                           null,
  is_visible         bit default b'1'                   not null,
  quantity           decimal(15, 4)                     null,
  quantity_unit      varchar(255)                       null,
  constraint product_versions_products_product_id_fk
  foreign key (product_id) references products (product_id)
    on update cascade
)
  engine = InnoDB
  charset = utf8mb4;

create index products_SKU_index
  on product_versions (product_id);

create table product_category_map
(
  product_id  int not null,
  category_id int not null,
  primary key (product_id, category_id),
  constraint product_category_map_ibfk_2
  foreign key (product_id) references products (product_id)
    on update cascade
    on delete cascade,
  constraint product_category_map_ibfk_1
  foreign key (category_id) references categories (category_id)
    on update cascade
    on delete cascade
)
  engine = InnoDB
  charset = utf8mb4;

create index product_category_map_ibfk_1
  on product_category_map (category_id);

create table stock
(
  stock_id   int auto_increment
    primary key,
  product_id int                                 not null,
  user_id    int                                 null,
  quantity   int                                 not null,
  addded_at  timestamp default CURRENT_TIMESTAMP not null,
  constraint stock_products_product_id_fk
  foreign key (product_id) references products (product_id)
    on update cascade,
  constraint stock_ibfk_1
  foreign key (user_id) references users (user_id)
    on update set null
    on delete set null
)
  engine = InnoDB
  charset = utf8mb4;

create index stock_ibfk_1
  on stock (user_id);

create index stock_products_SKU_fk
  on stock (product_id);

create table transactions
(
  transaction_id int auto_increment
    primary key,
  user_id        int                                 null,
  product_id     int                                 null,
  value          decimal(15, 2)                      not null,
  tag            varchar(255)                        null,
  added_at       timestamp default CURRENT_TIMESTAMP not null,
  updated_at     timestamp                           null,
  constraint fk_transactions_users
  foreign key (user_id) references users (user_id)
    on update cascade
    on delete set null,
  constraint transactions_products_product_id_fk
  foreign key (product_id) references products (product_id)
    on update cascade
    ON DELETE SET NULL
)
  engine = InnoDB
  charset = utf8mb4;

create index `transactions__product.SKU_id_fk`
  on transactions (product_id);

create index user_id
  on transactions (user_id);
