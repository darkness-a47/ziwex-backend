CREATE TABLE IF NOT EXISTS product_categories (
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    CONSTRAINT product_categories_pkey PRIMARY KEY(product_id, category_id)
)