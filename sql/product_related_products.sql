CREATE TABLE IF NOT EXISTS product_related_products (
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    related_product_id INT REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT product_related_products_pkey PRIMARY KEY(product_id, related_product_id)
)