CREATE TABLE IF NOT EXISTS files (
    id serial unique primary key,
    filename text,
    content_type text,
    file_id UUID unique,
    hash_md5 UUID,
    size bigint
);