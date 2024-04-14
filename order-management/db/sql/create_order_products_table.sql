CREATE TABLE IF NOT EXISTS order_products (
    order_id UUID REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(product_id),
    quantity INT NOT NULL,
    PRIMARY KEY(order_id, product_id)
);