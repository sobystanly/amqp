SELECT o.order_id, o.customer_id, o.payment_status, o.order_status,
       op.product_id, op.quantity, p.name, p.price, p.description,
       p.product_type, p.quantity_available
FROM orders AS o
         INNER JOIN order_products AS op
                    ON o.order_id = op.order_id
        INNER JOIN products AS p
        ON op.product_id = p.product_id;