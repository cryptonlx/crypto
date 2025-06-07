## Setup (Local Environment)

### Requirements

1. Fetch Repository

``` /bin/sh
git clone https://github.com/cryptonlx/crypto.git
cd crypto
```

2. PostgreSQL instance

- execute DDL on a new database: `./schemas/init_schema_001.sql`

### Commands

#### start server `DATABASE_URL=<> go run ./cmd/server`

#### tests `SERVER_URL=<server_url> go run ./cmd/e2e_tests`

Spins up a client that executes the [???test_plan](???).

# Design/Development Approach

The HTTP APIs will be drafted and tests will be written accordingly to verify the behavior via the API contract.
The tests are end-to-end and will require external connections (db etc.).

- Understand requirements.
- Create [Test Plan](./test_plan.md).
- Write failing test cases iteratively @ [e2etests](./cmd/e2e_tests)
    - Implement controller, application logic and database layer.
    - Verify, refine and update API contract.

# API Endpoints

1. **[API-WALL-DEP]** Deposit to user's wallet.

   `/POST /wallet/deposit`

    - See [Wallet Idempotency](#wallet-idempotency)
2. **[API-WALL-WDR]** Withdraw user's wallet

   `/POST /wallet/withdrawal`

    - See [Wallet Idempotency](#wallet-idempotency)
3. **[API-WALL-TRF]** Transfer from one user's wallet to another user's wallet.

   `/POST /wallet/transfer`

    - See [Wallet Idempotency](#wallet-idempotency)
    - currency type of wallets must match.
4. **[API-USER-BAL]** Get balances of user's wallets.

   `/GET /user/{username}/wallets`

5. **[API-USER-TXH]** Get user's transaction history.

   `/GET /user/{username}/transactions`

6. **[API-USER-NEW]** Create new user.

   `/POST /user`

   Fails on conflict with existing user. User identification by `request.username`.
7. **[API-WALL-NEW]** Create new wallet for user.

   `/POST /wallet`

### Glossary

- User: an account that can own wallets.
- Wallet: Value store of a currency owned by a User.
- Transaction: A record that changes a wallet's balance.
- Ledger: Authoritative set of records for wallet debit/credit.

### Functional Requirements

- Create new user if not exists.
- Each user can have multiple wallets.
- Supports deposit and withdrawal.
- Supports transfer from/to wallets.
- Viewing of wallet balance.
- Viewing of transaction history.

### Non-functional Requirements

#### Wallet Idempotency

- Idempotency for deposit/withdraw/transfer requests:
    - Include a 13-digit unix timestamp as nonce field for request identification.
    - Subsequent requests from same user with same `nonce` will be treated as duplicitous.
    - Each request can succeed at most once. Retries are allowed.

#### Atomicity

- Each request should be processed atomically and serialized across affected tables to ensure data integrity.

### Areas to Improve on

- Testing
    - Add table-driven unit tests to test in packages in isolation for more confidence.
  - Scalability
      - Consider service availability/maintainability for massive operations.
          - Set rate limiting per endpoint basis to stabilise server. Use Redis to store rate limiter's state a server
            cluster.
      - Support for asynchronous services.
          - For example, notify on operation fail/success, balance change etc.
- Payload selection
    - List responses should have pagination parameters to return subset as result.
- Greater API Flexibility
    - Currency Value and Unit Type
        - Support for cross-currency transfer.
        - Decide on better value type assignment in PostgreSQL for accuracy. There are a few options to choose from:
            - Multiply value by 1000x and store as `bigint`.
            - Store as floating point.
            - Store as `money`.
    - Conversion and Broker Fees Calculation for Effective Transaction Value)
- Security
    - Principal authorization for wallet actions via token issuance or session.
    - Ensure request integrity via payload signing.
- Observability
    - Record failed transactions for auditing.
    - Request tracing and logging for easy debugging.
