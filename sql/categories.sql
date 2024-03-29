CREATE TABLE IF NOT EXISTS categories (
    id SERIAL UNIQUE PRIMARY KEY,
    title text,
    image_id uuid,
    description text,
    parent_category_id int,
    Tags text[],
    CONSTRAINT fk_parent_category 
        FOREIGN KEY(parent_category_id) 
        REFERENCES categories(id)
        ON DELETE SET NULL,
    CONSTRAINT fk_image
        FOREIGN KEY(image_id)
        REFERENCES files(file_id)
        ON DELETE SET NULL
);