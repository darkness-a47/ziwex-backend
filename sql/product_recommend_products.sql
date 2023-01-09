CREATE TABLE IF NOT EXISTS product_recommend_products (
    product_id INT,
    recommend_product_id INT,
    PRIMARY KEY(product_id, recommend_product_id),
    CONSTRAINT fk_product
        FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_recommend_product
        FOREIGN KEY(recommend_product_id) REFERENCES products(id) ON DELETE CASCADE
)