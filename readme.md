
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

#### start server `go run ./cmd/server`
#### tests `<server_url> go run ./cmd/e2e_tests`
Spins up a client that executes the [???test_plan](???).

# Design Approach

The HTTP APIs will be drafted and tests will be written accordingly to verify the behavior via the API contract.
The tests are end-to-end and will require external connections (db etc.).

# APIs
1. Deposit to specify user wallet

    `/POST /user/deposit_intent`

    - IDEMPOTENT
2. Withdraw from specify user wallet

    `/POST /user/withdrawal`

    - IDEMPOTENT
3. Transfer from one user to another user

    `/POST /user/transfer`

    - IDEMPOTENT
4. Get specify user balance

    `/GET /user`

5. Get specify user transaction history

    `/GET /user/wallet_history`

### Functional requirements

### Non-functional requirements

#### Wallet Modification
- Idempotency for wallet modification requests.
  - Request will require a nonce field as timestamp id that is applicable. Requests with same nonce will be treated as duplicitous.
  - There will be no operation retries.
    - For example, user request (nonce:`001`) deposit of $50. Outcome must be success or failure (subsequent request will receive the same response).
      - If success, return status: `SUCCESS`.
      - If fail, return status: `FAILURE_MESSAGE`.
      - To retry, send another user request (nonce:`002`).

- Operations should be atomic and serialized across affected tables to ensure data integrity.
- A response will be sent to indicate if an operation failed or successful.

### Things to Improve on (Current and Future Scope)
- Scalability
  - Consider service availability/maintainability for massive operations.
    - Upgrade in-mem cache to Redis so that server is stateless.
    - Set rate limiting per endpoint basis to stabilise server.
  - Support for asynchronous communications
    - For example, notify on operation success, balance change etc.
- Payload selection
  - List responses should have pagination parameters to return subset as result.
- Greater API Flexibility
  - Currency Unit
  - Conversion and Broker Fees Calculation (Effective Transaction Value)
- Security
  - User authentication via token issuance or session.
  - Request authentication via payload signing.