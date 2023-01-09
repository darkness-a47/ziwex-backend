CREATE TABLE IF NOT EXISTS product_images (
    product_id INT,
    image_id INT,
    PRIMARY KEY(product_id, image_id),
    CONSTRAINT fk_product
        FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_image
        FOREIGN KEY(image_id) REFERENCES files(id) ON DELETE CASCADE
)