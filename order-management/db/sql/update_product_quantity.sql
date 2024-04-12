UPDATE products AS p
SET quantity_available = quantity_available - (
    SELECT quantity
    FROM order_products AS op
    WHERE order_id = $1
    AND p.product_id = op.product_id
)
WHERE product_id IN (
    SELECT product_id
    FROM order_products
    WHERE order_id = $1
)