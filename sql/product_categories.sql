CREATE TABLE IF NOT EXISTS product_categories (
    product_id INT,
    category_id INT,
    PRIMARY KEY(product_id, category_id),
    CONSTRAINT fk_product
        FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_category
        FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
)