-- +migrate Up

CREATE TABLE IF NOT EXISTS categories (
  category_id INT(11)    NOT NULL  AUTO_INCREMENT,
  name        INT(11)    NOT NULL,
  visible     TINYINT(1) NOT NULL  DEFAULT '1',
  active      TINYINT(1) NOT NULL  DEFAULT '0',
  is_root     TINYINT(1) NOT NULL  DEFAULT '1',
  PRIMARY KEY (category_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS category_parent_map (
  cateogry_id        INT(11) NOT NULL  AUTO_INCREMENT,
  parent_category_id INT(11) NOT NULL,
  PRIMARY KEY (cateogry_id, parent_category_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS transactions (
  payment_id INT(11)        NOT NULL              AUTO_INCREMENT,
  user_id    INT(11)        NOT NULL,
  SKU_id     INT(11)        NULL,
  value      DECIMAL(15, 2) NOT NULL,
  tag        VARCHAR(255)                         DEFAULT NULL,
  added_at   TIMESTAMP      NOT NULL              DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP      NULL                  DEFAULT NULL,
  PRIMARY KEY (payment_id),
  INDEX (user_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

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
  quantity      DECIMAL(15, 2)                     NULL,
  quantity_unit VARCHAR(255)                       NULL
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

CREATE INDEX products_SKU_index
  ON products (SKU_id);


CREATE TABLE IF NOT EXISTS product_category_map (
  product_id  INT(11) NOT NULL  AUTO_INCREMENT,
  category_id INT(11) NOT NULL,
  PRIMARY KEY (product_id, category_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

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
    ON DELETE CASCADE,
  CONSTRAINT stock_ibfk_1
  FOREIGN KEY (user_id) REFERENCES users (user_id)
    ON UPDATE SET NULL
    ON DELETE SET NULL
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

CREATE INDEX stock_products_SKU_fk
  ON stock (SKU_id);

CREATE INDEX stock_ibfk_1
  ON stock (user_id);

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

CREATE TABLE IF NOT EXISTS user_auths (
  user_auth_id INT(11)      NOT NULL  AUTO_INCREMENT,
  user_id      INT(11)      NOT NULL,
  method       VARCHAR(255) NOT NULL,
  value        TINYBLOB,
  PRIMARY KEY (user_auth_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

ALTER TABLE transactions
  ADD CONSTRAINT fk_transactions_users FOREIGN KEY (user_id) REFERENCES users (user_id)
  ON DELETE NO ACTION
  ON UPDATE NO ACTION;

ALTER TABLE product_category_map
  ADD CONSTRAINT product_category_map_ibfk_1 FOREIGN KEY (category_id) REFERENCES categories (category_id)
  ON DELETE CASCADE
  ON UPDATE CASCADE,
  ADD CONSTRAINT product_category_map_ibfk_2 FOREIGN KEY (product_id) REFERENCES products (product_id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE stock
  ADD CONSTRAINT stock_ibfk_1 FOREIGN KEY (user_id) REFERENCES users (user_id)
  ON DELETE SET NULL
  ON UPDATE SET NULL,
  ADD CONSTRAINT stock_ibfk_2 FOREIGN KEY (SKU_id) REFERENCES products (product_id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE user_auths
  ADD CONSTRAINT fk_user_auths_user FOREIGN KEY (user_id) REFERENCES users (user_id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;