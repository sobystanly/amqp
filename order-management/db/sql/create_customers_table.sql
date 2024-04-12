CREATE TABLE IF NOT EXISTS customers (
    customer_id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL
);