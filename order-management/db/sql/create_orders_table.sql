CREATE TABLE IF NOT EXISTS orders (
    order_id UUID PRIMARY KEY,
    customer_id UUID REFERENCES customers(customer_id),
    payment_status TEXT NOT NULL,
    order_status TEXT NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);