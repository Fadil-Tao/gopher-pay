CREATE TABLE IF NOT EXISTS wallets(
    id UUID PRIMARY KEY,
    balance numeric(15,0) default 0,
    user_id int,
    foreign key (user_id) references users(id)
);