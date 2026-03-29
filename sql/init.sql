CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    item_name VARCHAR(255),
    sku VARCHAR(100) UNIQUE NOT NULL,
    stock_count INT NOT NULL,
    price DECIMAL(10,2),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO inventory (item_name, sku, stock_count, price) VALUES
('Wireless Mouse', 'MOUSE-001', 25, 499.99),
('Mechanical Keyboard', 'KEY-001', 15, 1999.99),
('USB-C Cable', 'CABLE-001', 50, 199.99),
('Laptop Stand', 'STAND-001', 8, 999.99),
('Webcam', 'CAM-001', 5, 1499.99);