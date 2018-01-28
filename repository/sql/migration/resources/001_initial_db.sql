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
  product_id INT(11)        NULL,
  value      DECIMAL(15, 2) NOT NULL,
  tag        VARCHAR(255)                         DEFAULT NULL,
  added_at   TIMESTAMP      NOT NULL              DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP      NULL                  DEFAULT NULL,
  PRIMARY KEY (payment_id),
  INDEX (user_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS products (
  product_id    INT(11)        NOT NULL              AUTO_INCREMENT,
  name          VARCHAR(255)   NOT NULL,
  GTIN          CHAR(14)                             DEFAULT NULL,
  price         DECIMAL(15, 2) NOT NULL,
  added_at      DATETIME       NOT NULL              DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME                             DEFAULT NULL,
  deleted_at    DATETIME                             DEFAULT NULL,
  is_visible    TINYINT(1)     NOT NULL              DEFAULT '1',
  quantity      DECIMAL(15, 2)                       DEFAULT NULL,
  quantity_unit VARCHAR(255)                         DEFAULT NULL,
  PRIMARY KEY (product_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS product_category_map (
  product_id  INT(11) NOT NULL  AUTO_INCREMENT,
  category_id INT(11) NOT NULL,
  PRIMARY KEY (product_id, category_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS stock (
  stock_id   INT(11)   NOT NULL  AUTO_INCREMENT,
  product_id INT(11)   NOT NULL,
  user_id    INT(11)             DEFAULT NULL,
  quantity   INT(11)   NOT NULL,
  addded_at  TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (stock_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE users
(
  user_id    INT AUTO_INCREMENT
    PRIMARY KEY,
  name       VARCHAR(191)                       NOT NULL,
  email      VARCHAR(191)                       NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME                           NULL,
  deleted_at DATETIME                           NULL,
  blocked_at TIMESTAMP                          NULL,
  is_admin   TINYINT(1) DEFAULT '0'             NOT NULL,
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
  ON UPDATE NO ACTION,
  ADD CONSTRAINT fk_transactions_products FOREIGN KEY (product_id) REFERENCES products (product_id)
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
  ADD CONSTRAINT stock_ibfk_2 FOREIGN KEY (product_id) REFERENCES products (product_id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE user_auths
  ADD CONSTRAINT fk_user_auths_user FOREIGN KEY (user_id) REFERENCES users (user_id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

# INSERT INTO payment (paymentID, userID, paymentValue, paymentTag, paymentTimeAdded, paymentTimeChanged) VALUES
#   (1, 20, '1.23', '', '2017-12-11 13:53:38', NULL),
#   (3, 20, '9999999999999.99', '', '2017-12-11 13:57:03', NULL),
#   (4, 20, '9999999999999.99', '', '2017-12-11 13:57:04', NULL),
#   (5, 20, '9999999999999.99', '', '2017-12-11 13:57:04', NULL),
#   (6, 20, '9999999999999.99', '', '2017-12-11 13:57:05', NULL),
#   (7, 20, '21344.00', '', '2017-12-11 13:57:20', NULL),
#   (8, 20, '21344.00', '', '2017-12-11 15:17:17', NULL),
#   (9, 20, '21344.00', '', '2017-12-11 15:21:01', NULL),
#   (10, 20, '21344.00', '', '2017-12-11 15:22:54', NULL),
#   (11, 20, '21344.00', '', '2017-12-11 15:23:30', NULL),
#   (12, 20, '21344.00', '', '2017-12-11 15:23:31', NULL),
#   (13, 20, '21344.00', '', '2017-12-11 15:23:33', NULL),
#   (14, 20, '21344.00', '', '2017-12-11 15:25:37', NULL),
#   (15, 20, '21344.00', '', '2017-12-11 15:25:39', NULL),
#   (16, 20, '21344.00', '', '2017-12-11 15:25:53', NULL),
#   (17, 20, '123.00', '', '2017-12-11 15:26:01', NULL),
#   (18, 20, '123.32', '', '2017-12-11 15:26:22', NULL),
#   (19, 20, '-1123.32', '', '2017-12-11 15:26:50', NULL),
#   (20, 20, '-1123.32', '', '2017-12-15 12:19:46', NULL);
#
# INSERT INTO product (productID, productName, productGTIN, productPrice, productTimeAdded, productTimeChanged, productTimeDeleted, productVisible, ProductQuantity, ProductQuantityUnit)
# VALUES
#   (3, 'Club-Mate 0,5l Flasche', '4029764001807', '1.50', '2017-12-15 14:10:18', NULL, NULL, 1, '0.50', 'liter'),
#   (4, 'Club-Mate 0,5l Flasche', '40297236400180', '1.50', '2017-12-15 14:11:06', NULL, NULL, 1, '0.50', 'liter'),
#   (5, 'Club-Mate 0,5l Flasche', '40297236402301', '1.50', '2017-12-15 14:11:22', NULL, NULL, 1, '0.50', 'liter'),
#   (6, 'Club-Mate 0,5l Flasche', '40297223640230', '1.50', '2017-12-15 14:11:44', NULL, NULL, 1, '0.50', 'liter'),
#   (8, 'Club-Mate 0,5l Flasche', '40297223640301', '1.50', '2017-12-15 14:12:06', NULL, NULL, 1, '0.50', 'liter'),
#   (9, 'Club-Mate 0,5l Flasche', '40297223640310', '1.50', '2017-12-15 14:12:14', NULL, NULL, 1, '0.50', 'liter'),
#   (10, 'Club-Mate 0,5l Flasche', '40292236403107', '1.50', '2017-12-15 14:17:34', NULL, NULL, 1, '0.50', ''),
#   (11, 'Club-Mate 0,5l Flasche', '40297223403107', '1.50', '2017-12-15 14:28:50', NULL, NULL, 1, '0.50', 'liter');
#
# INSERT INTO USER (userID, userName, userEmail, userTimeCreated, userTimeChanged, userTimeDeleted, userTimeBlocked, userIsAdmin)
# VALUES
#   (15, 'Marei', 'marei@binary-kitchen.de', '2017-11-21 16:00:29', NULL, NULL, NULL, 0),
#   (17, 'Timo', 'timo@binary-kitchen.de', '2017-11-21 16:00:43', NULL, NULL, NULL, 0),
#   (18, 'marove2000', 'timo@awefawef.de', '2017-11-27 15:40:08', NULL, NULL, NULL, 1),
#   (20, 'marove', '', '2017-12-08 13:11:11', NULL, NULL, NULL, 1);
#
# INSERT INTO userAuth (userAuthID, userID, userAuthMethod, userAuthValue) VALUES
#   (9, 15, 'password',
#    0x2432612431302478754b7975716847632e71466f3546614d7253535a2e5862632e74756645455734654e67506c394a5a4c3879303250654a37443553),
#   (10, 17, 'password',
#    0x2432612431302453734d757135445448566977436f396d7a38754f792e4d72743852635131546630454f67525648786b6f43753435656a4972796c36),
#   (11, 18, 'password',
#    0x24326124313024486d77514d2e74784c642f5246482e594d447166362e487368314d396933457469454767374c33474b7a693348616c3854756d3357),
#   (13, 20, 'ldap', NULL);