CREATE TABLE IF NOT EXISTS product_images (
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    image_id INT  REFERENCES files(id) ON DELETE CASCADE,
    CONSTRAINT product_images_pkey PRIMARY KEY(product_id, image_id)
)