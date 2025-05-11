CREATE DATABASE IF NOT EXISTS checkout_service;

CREATE TABLE IF NOT EXISTS products (
  id INT NOT NULL AUTO_INCREMENT,
  sku VARCHAR(20) DEFAULT NULL,
  name VARCHAR(20) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  inventory_qty INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (id)
);

INSERT INTO products (sku, name, price, inventory_qty)
SELECT * FROM (
    SELECT '120P90', 'Google Home', 49.99, 10
    UNION ALL
    SELECT '43N23P', 'MacBook Pro', 5399.99, 5
    UNION ALL
    SELECT 'A304SD', 'Alexa Speaker', 109.50, 10
    UNION ALL
    SELECT '234234', 'Raspberry Pi B', 30.00, 2
) AS tmp
WHERE NOT EXISTS (SELECT 1 FROM products LIMIT 1);

CREATE TABLE IF NOT EXISTS discounts (
  id INT NOT NULL AUTO_INCREMENT,
  type INT NOT NULL,
  is_percentage INT NOT NULL DEFAULT 0,
  amount INT NOT NULL DEFAULT 0,
  required_qty INT NOT NULL DEFAULT 0,
  final_qty INT NOT NULL DEFAULT 0,
  free_id_product INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (id)
);

INSERT INTO discounts (type, is_percentage, amount, required_qty, final_qty, free_id_product, created_at, updated_at, deleted_at)
SELECT 1, 0, 0, 1, 0, 4, '2025-05-09 20:03:41', '2025-05-09 20:03:41', NULL
WHERE NOT EXISTS (SELECT 1 FROM discounts WHERE type = 1 AND amount = 0);

INSERT INTO discounts (type, is_percentage, amount, required_qty, final_qty, free_id_product, created_at, updated_at, deleted_at)
SELECT 2, 0, 0, 3, 2, 0, '2025-05-09 20:05:14', '2025-05-09 20:05:14', NULL
WHERE NOT EXISTS (SELECT 1 FROM discounts WHERE type = 2 AND amount = 0);

INSERT INTO discounts (type, is_percentage, amount, required_qty, final_qty, free_id_product, created_at, updated_at, deleted_at)
SELECT 3, 1, 10, 3, 0, 0, '2025-05-09 20:07:04', '2025-05-09 20:07:04', NULL
WHERE NOT EXISTS (SELECT 1 FROM discounts WHERE type = 3 AND amount = 10);

CREATE TABLE IF NOT EXISTS products_discounts (
  product_id INT NOT NULL,
  discounts_id INT NOT NULL,
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  FOREIGN KEY (discounts_id) REFERENCES discounts(id) ON DELETE CASCADE
);

INSERT INTO products_discounts (product_id, discounts_id)
SELECT * FROM (
    SELECT 2, 1
    UNION ALL
    SELECT 1, 2
    UNION ALL
    SELECT 3, 3
) AS tmp
WHERE NOT EXISTS (SELECT 1 FROM products_discounts LIMIT 1);

CREATE TABLE IF NOT EXISTS transaction (
  id INT NOT NULL AUTO_INCREMENT,
  customer_name VARCHAR(20) DEFAULT NULL,
  customer_phone VARCHAR(20) DEFAULT NULL,
  status INT NOT NULL DEFAULT 0,
  total_price DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS transaction_detail (
  id INT NOT NULL AUTO_INCREMENT,
  transaction_id INT NOT NULL,
  product_id INT NOT NULL,
  product_name VARCHAR(20) DEFAULT NULL,
  sku VARCHAR(20) DEFAULT NULL,
  qty INT NOT NULL,
  price DECIMAL(10, 2) NOT NULL DEFAULT 0,
  status INT NOT NULL DEFAULT 0,
  total_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);