CREATE TABLE IF NOT EXISTS product_related_products (
    product_id INT,
    related_product_id INT,
    PRIMARY KEY(product_id, related_product_id),
    CONSTRAINT fk_product
        FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_related_product
        FOREIGN KEY(related_product_id) REFERENCES products(id) ON DELETE CASCADE
)