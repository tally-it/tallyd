-- +migrate Up
CREATE TABLE categories
(
  category_id INT AUTO_INCREMENT
    PRIMARY KEY,
  name        INT                    NOT NULL,
  visible     TINYINT(1) DEFAULT '1' NOT NULL,
  active      TINYINT(1) DEFAULT '0' NOT NULL,
  is_root     TINYINT(1) DEFAULT '1' NOT NULL
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

CREATE TABLE category_parent_map
(
  cateogry_id        INT AUTO_INCREMENT,
  parent_category_id INT NOT NULL,
  PRIMARY KEY (cateogry_id, parent_category_id)
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

CREATE TABLE product_category_map
(
  product_id  INT AUTO_INCREMENT,
  category_id INT NOT NULL,
  PRIMARY KEY (product_id, category_id),
  CONSTRAINT product_category_map_ibfk_1
  FOREIGN KEY (category_id) REFERENCES categories (category_id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
)
  ENGINE = InnoDB
  CHARSET = utf8mb4;

CREATE INDEX product_category_map_ibfk_1
  ON product_category_map (category_id);

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

ALTER TABLE product_category_map
  ADD CONSTRAINT product_category_map_ibfk_2
FOREIGN KEY (product_id) REFERENCES products (product_id)
  ON UPDATE CASCADE
  ON DELETE CASCADE;

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
