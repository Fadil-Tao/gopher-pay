CREATE TYPE transaction_type as enum('withdrawal', 'topup', 'transfer');

CREATE TABLE IF NOT EXISTS transactions(
    id UUID primary key,
    source_wallet_id UUID REFERENCES wallets(id),
    destination_wallet_id UUID REFERENCES wallets(id),
    amount numeric(15,2) not null,
    transaction_type transaction_type,
    description text,
    created_at TIMESTAMP DEFAULT current_timestamp
);