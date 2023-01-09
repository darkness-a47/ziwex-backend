CREATE TABLE IF NOT EXISTS products (
    id SERIAL UNIQUE PRIMARY KEY,
    url text UNIQUE,
    title TEXT,
    description TEXT,
    price FLOAT8,
    options JSON[],
    description_key_value JSON[],
    main_image INT
)