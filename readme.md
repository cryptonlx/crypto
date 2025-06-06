
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

# Design Approach

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
3. **[API-WALL-TRF]** Transfer from one user's account to another user's account.

    `/POST /wallet/transfer`

    - See [Wallet Idempotency](#wallet-idempotency)
    - currency type of wallets must match.
4. **[API-WALL-BAL]** Get balances of user's wallets.

    `/GET /user/{username}/balance`

5. **[API-WALL-HST]** Get user's transaction history.

    `/GET /user/{username}/transactions`

6. **[API-USER-NEW]** Create new user.

   `/POST /user`

    Fails on conflict with existing user. User identification by `request.username`.
7. **[API-WALL-NEW]** Create new wallet for user.

    `/POST /wallet`

### Glossary

- User: an account that can own wallets.
- Wallet: Value store of a currency owned by a User.
- Transaction: An set of operations to process a deposit/withdraw/transfer request.
- Ledger: Authoritative set of records for wallet debit/credit.

### Functional Requirements

- Create new user if not exists.
- Each user can have multiple wallets.
- Supports deposit and withdrawal.
- Supports transfer from/to wallets.
- Enable viewing of wallet balance and transaction history.

### Non-functional Requirements

#### Wallet Idempotency
- Idempotency for modification requests:
  - Include a 13-digit unix timestamp as nonce field for request identification. Subsequent requests with same `username` and `nonce` will be treated as duplicitous.
  - There will be no operation retries. Subsequent request will receive the same response as the first.

#### Atomicity
- Operations should be atomic and serialized across affected tables to ensure data integrity.

### Things to Improve on (Future Scope)
- Testing
  - Can add table-driven unit tests to test in isolation for more confidence.
- Scalability
  - Consider service availability/maintainability for massive operations.
    - Set rate limiting per endpoint basis to stabilise server. Can use Redis to store rate limiter's state a server cluster.
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
  - User authentication via token issuance or session.
  - Request authentication via payload signing.
- Observability
  - Request tracing and logging for easy debugging.