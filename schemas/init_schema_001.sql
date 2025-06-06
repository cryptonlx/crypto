drop table if exists  wallets;
drop table if exists transactions;
drop table if exists  user_accounts;

create table public.user_accounts
(
    id      bigint primary key generated always as identity,
    username text not null
        unique
);


create table transactions
(

    id      bigint primary key generated always as identity,
    user_account_id bigint references user_accounts (id) not null,
    nonce text not null,
    status_message text not null ,
    mode text not null
);

create unique index  transactions_nonce_idx on transactions using btree(user_account_id,nonce) ;

comment on column  transactions.nonce is '13 digit epoch i.e 1749199885000';
comment on column  transactions.status_message is 'in_progress, success, error_*';
comment on column  transactions.mode is 'deposit, withdrawal, transfer';


create table public.wallets
(
    id              bigint primary key generated always as identity,
    user_account_id bigint references user_accounts (id) NOT NULL,
    currency_type   text,
    value           bigint
);

create unique index  wallets_currency_idx on wallets using btree(user_account_id,currency_type) ;

comment on column  wallets.currency_type is 'USD';