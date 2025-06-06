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