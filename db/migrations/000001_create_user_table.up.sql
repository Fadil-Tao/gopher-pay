CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL unique,
    name VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    salt text not null,
    phone varchar(255) not null unique,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);