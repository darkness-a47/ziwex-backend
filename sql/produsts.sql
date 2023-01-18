CREATE TABLE IF NOT EXISTS products (
    id SERIAL UNIQUE PRIMARY KEY,
    url text UNIQUE,
    title TEXT,
    description TEXT,
    price FLOAT8,
    options JSONB[],
    description_key_value JSONB[],
    main_image_index INT
)