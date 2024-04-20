CREATE TABLE users (
    id serial PRIMARY KEY ,
    phone_number VARCHAR (13) UNIQUE NOT NULL,
    full_name VARCHAR (60) NOT NULL,
    password VARCHAR (64) NOT NULL,
    successful_login INT DEFAULT 0
)