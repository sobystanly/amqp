UPDATE orders
SET payment_status = $1, order_status = $2
WHERE order_id = $3;