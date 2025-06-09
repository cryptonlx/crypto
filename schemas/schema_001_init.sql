DROP TABLE if EXISTS ledgers;
DROP TABLE if EXISTS transactions;
DROP TABLE if EXISTS wallets;
DROP TABLE if EXISTS user_accounts;

CREATE TABLE public.user_accounts
(
    id       BIGINT PRIMARY key generated always AS IDENTITY,
    username text NOT NULL UNIQUE
);
CREATE TABLE public.wallets
(
    id              BIGINT PRIMARY key generated always AS IDENTITY,
    user_account_id BIGINT REFERENCES user_accounts (id) NOT NULL,
    currency        text                                 NOT NULL,
    balance         NUMERIC(20, 6)                       NOT NULL CHECK (balance >= 0)
);

CREATE UNIQUE index wallets_currency_idx ON wallets USING btree (user_account_id, currency);
comment ON COLUMN wallets.currency IS 'ISO4217 compliant USD';


create table transactions
(
    id      bigint primary key generated always as identity,
    requestor_id bigint references user_accounts (id) not null,
    nonce bigint not null,
    status text not null,
    operation  text                           NOT NULL,
    metadata JSONB default '{}'::jsonb,
    created_at TIMESTAMP WITH TIME zone
);

create unique index  transactions_nonce_idx on transactions using btree(requestor_id,nonce) ;
comment ON COLUMN transactions.operation IS 'deposit, withdrawal, transfer';
comment on column  transactions.nonce is '13 digit epoch i.e 1749199885000';
comment on column  transactions.status is 'pending, success, error_*';

CREATE TABLE ledgers
(
    id         BIGINT PRIMARY key generated always AS IDENTITY,
    wallet_id  BIGINT REFERENCES wallets (id) NOT NULL,
    transaction_id      bigint references transactions(id) not null,
    entry_type text                           NOT NULL,
    amount     NUMERIC(20, 6)                 NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP WITH TIME zone,
    balance    numeric(20, 6) CHECK (balance >= 0)
);

comment ON COLUMN ledgers.entry_type IS 'debit,credit';
comment ON COLUMN ledgers.balance IS 'wallet balance after operation';