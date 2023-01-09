CREATE TABLE IF NOT EXISTS product_recommend_products (
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    recommend_product_id INT REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT product_recommend_products_pkay PRIMARY KEY(product_id, recommend_product_id)
)