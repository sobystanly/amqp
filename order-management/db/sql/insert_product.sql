INSERT INTO products (product_id, name, price, description, quantity_available, product_type) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT(product_id) DO NOTHING;