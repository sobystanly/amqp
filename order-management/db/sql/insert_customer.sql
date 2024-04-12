INSERT INTO customers(customer_id, first_name, last_name, email) VALUES ($1, $2, $3, $4)
ON CONFLICT(customer_id) DO NOTHING;