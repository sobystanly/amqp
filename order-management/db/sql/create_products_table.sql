CREATE TABLE IF NOT EXISTS products (
    product_id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    description TEXT,
    quantity_available INT NOT NULL,
    product_type TEXT NOT NULL
);