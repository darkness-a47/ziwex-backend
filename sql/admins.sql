CREATE TABLE IF NOT EXISTS admins (
    id Serial UNIQUE PRIMARY KEY ,
    firstname text,
    lastname text,
    email text,
    password text,
    username text
);
