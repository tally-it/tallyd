-- +migrate Up
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

CREATE TABLE products
(
  product_id    INT AUTO_INCREMENT
    PRIMARY KEY,
  SKU_id        INT                                NOT NULL,
  name          VARCHAR(255)                       NOT NULL,
  GTIN          CHAR(14)                           NULL,
  price         DECIMAL(15, 2)                     NOT NULL,
  added_at      DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at    DATETIME                           NULL,
  is_visible    BIT DEFAULT b'1'                   NOT NULL,
  quantity      DECIMAL(15, 4)                     NULL,
  quantity_unit VARCHAR(255)                       NULL
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

CREATE INDEX products_SKU_index
  ON products (SKU_id);

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

CREATE TABLE stock
(
  stock_id  INT AUTO_INCREMENT
    PRIMARY KEY,
  SKU_id    INT                                 NOT NULL,
  user_id   INT                                 NULL,
  quantity  INT                                 NOT NULL,
  addded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  CONSTRAINT stock_products_SKU_fk
  FOREIGN KEY (SKU_id) REFERENCES products (SKU_id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

CREATE INDEX stock_products_SKU_fk
  ON stock (SKU_id);

CREATE INDEX stock_ibfk_1
  ON stock (user_id);

CREATE TABLE transactions
(
  transaction_id INT AUTO_INCREMENT
    PRIMARY KEY,
  user_id        INT                                 NULL,
  SKU_id         INT                                 NULL,
  value          DECIMAL(15, 2)                      NOT NULL,
  tag            VARCHAR(255)                        NULL,
  added_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at     TIMESTAMP                           NULL,
  CONSTRAINT `transactions__product.SKU_id_fk`
  FOREIGN KEY (SKU_id) REFERENCES products (SKU_id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

CREATE INDEX user_id
  ON transactions (user_id);

CREATE INDEX `transactions__product.SKU_id_fk`
  ON transactions (SKU_id);

CREATE TABLE user_auths
(
  user_auth_id INT AUTO_INCREMENT
    PRIMARY KEY,
  user_id      INT          NOT NULL,
  method       VARCHAR(255) NOT NULL,
  value        TINYBLOB     NULL
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

CREATE INDEX fk_user_auths_user
  ON user_auths (user_id);

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

ALTER TABLE stock
  ADD CONSTRAINT stock_ibfk_1
FOREIGN KEY (user_id) REFERENCES users (user_id)
  ON UPDATE SET NULL
  ON DELETE SET NULL;

ALTER TABLE transactions
  ADD CONSTRAINT fk_transactions_users
FOREIGN KEY (user_id) REFERENCES users (user_id)
  ON UPDATE CASCADE
  ON DELETE SET NULL;

ALTER TABLE user_auths
  ADD CONSTRAINT fk_user_auths_user
FOREIGN KEY (user_id) REFERENCES users (user_id)
  ON UPDATE CASCADE
  ON DELETE CASCADE;
