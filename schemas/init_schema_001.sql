drop table if exists transactions;
drop table if exists  wallets;

create table public.user_accounts
(
    id      bigint primary key generated always as identity,
    username text not null
        unique
);


create table public.wallets
(
    id              bigint primary key generated always as identity,
    user_account_id bigint references user_accounts (id) NOT NULL,
    currency_type   text,
    value           bigint
);