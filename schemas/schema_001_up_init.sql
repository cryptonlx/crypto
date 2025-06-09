CREATE TABLE public.user_accounts
(
    id       bigint GENERATED always AS IDENTITY PRIMARY KEY,
    username text NOT NULL UNIQUE
);


CREATE TABLE public.wallets
(
    id              bigint GENERATED always AS IDENTITY PRIMARY KEY,
    user_account_id bigint         NOT NULL REFERENCES public.user_accounts,
    currency        text           NOT NULL,
    balance         numeric(20, 6) NOT NULL
        CONSTRAINT wallets_balance_check CHECK (balance >= (0)::numeric)
);

COMMENT ON COLUMN public.wallets.currency IS 'ISO4217 compliant USD';

CREATE UNIQUE INDEX wallets_currency_idx ON public.wallets (user_account_id, currency);


CREATE TABLE public.transactions
(
    id           bigint GENERATED always AS IDENTITY PRIMARY KEY,
    requestor_id bigint                   NOT NULL REFERENCES public.user_accounts,
    nonce        bigint                   NOT NULL,
    status       text                     NOT NULL,
    OPERATION    text                     NOT NULL,
    metadata     JSONB DEFAULT '{}'::JSONB,
    created_at   timestamp WITH TIME ZONE NOT NULL
);

COMMENT ON COLUMN public.transactions.nonce IS '13 digit epoch i.e 1749199885000';
COMMENT ON COLUMN public.transactions.status IS 'pending, success, error_*';
COMMENT ON COLUMN public.transactions.operation IS 'deposit, withdrawal, transfer';

CREATE UNIQUE INDEX transactions_nonce_idx ON public.transactions (requestor_id, nonce);
CREATE INDEX transactions_requestor_id_index ON public.transactions (requestor_id);


CREATE TABLE public.ledgers
(
    id             bigint GENERATED always AS IDENTITY PRIMARY KEY,
    wallet_id      bigint                   NOT NULL REFERENCES public.wallets,
    transaction_id bigint                   NOT NULL REFERENCES public.transactions,
    entry_type     text                     NOT NULL,
    amount         numeric(20, 6)           NOT NULL
        CONSTRAINT ledgers_amount_check CHECK (amount > (0)::numeric),
    created_at     timestamp WITH TIME ZONE NOT NULL,
    balance        numeric(20, 6)
        CONSTRAINT ledgers_balance_check CHECK (balance >= (0)::numeric)
);

COMMENT ON COLUMN public.ledgers.entry_type IS 'debit,credit';
COMMENT ON COLUMN public.ledgers.balance IS 'wallet balance after operation';

CREATE INDEX ledgers_wallet_id_index ON public.ledgers (wallet_id);