CREATE TYPE transaction_type as enum('withdrawal', 'topup', 'transfer');

CREATE TABLE IF NOT EXISTS transactions(
    id UUID primary key,
    sender_id int references users(id),
    receiver_id int references users(id),
    amount numeric(15,2) not null,
    transaction_type transaction_type,
    description text,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);