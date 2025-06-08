# Test Plain (HTTP)

## Prerequisites

Ensure API test server is ready.

## User Stories

[US-001] User can deposit money into his/her wallet\
[US-002] User can withdraw money from his/her wallet\
[US-003] User can send money to another user\
[US-004] User can check his/her wallet balance\
[US-005] User can view his/her transaction history

- [x] [T_0001] User Creation\
  User Stories: [US-004]\
  Test if user exist by get balance endpoint.
    - [x] [T_0001_001] Get user by `username` yet to be created.
        - Endpoint: [API-USER-BAL]
        - [x] Status: 400
        - [x] Error Message = `"resource: user not found"`
    - [x] [T_0001_002] Create user with `username`
        - Endpoint: [API-USER-NEW]
        - [x] Status: 200
    - [x] [T_0001_003] Create duplicate user `username` (same as above)
        - Endpoint: [API-USER-NEW]
        - [x] Status: 400
        - [x] Error Message = `"unique_violation"`
    - [x] [T_0001_004] Get existing user `username` (same as above)
        - Endpoint: [API-USER-BAL]
        - [x] Status: 200

- [x] [T_0002] - Get New User Transaction History\
  User Stories: [US-005]
    - [x] [T_0002_001] Get history by user yet to be created.
        - Endpoint: [API-USER-TXH]
        - [x] Status: 400
        - [x] Error Message = `"resource: user not found"`
    - [x] [T_0002_002] Create user with `username`
        - Endpoint: [API-USER-NEW]
        - [x] Status: 200
    - [x] [T_0002_003] Get History
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] Result: []
- [x] [T_0003] - Wallet Creation \
  User Stories: [US-004]
    - [x] [T_0003_001] Create user with `username`
        - Endpoint: [API-USER-NEW]
        - [x] Status: 200
    - [x] [T_0003_002] Get Balance
        - Endpoint: [API-USER-BAL]
        - [x] Status: 200
        - [x] Result: []
    - [x] [T_0003_003] Create `SGD` wallet
        - [x] Status: 200
        - [x] Result: `response.wallet.id` != 0
    - [x] [T_0003_004] Get Balance
        - Endpoint: [API-USER-BAL]
        - [x] Status: 200
        - [x] Result: `response.wallets.len` == 1
        - [x] Result: `response.wallets[0].currency` == `SGD`
- [x] [T_0004] - New User: Deposit Error (Negative Amount)\
  User Stories: [US-001], [US-005]
    - [x] [Setup] Do [T_0003]
    - [x] [T_0004_001] Deposit `amount` less than or equals to 0 should fail
        - Endpoint: [API-USER-DEP]
        - [x] Status: 400
        - [x] Error Message = `"invalid_amount"`
    - [x] [T_0004_002] Get Balance
        - Endpoint: [API-USER-BAL]
        - [x] Status: 200
        - [x] Result: `after.balance`
        - [x] Assert: `before.balance` == `after.balance`
- [x] [T_0005] - New User: Deposit Wallet Success (Positive Amount)\
  User Stories: [US-001], [US-005]
    - [x] [Setup] Do [T_0003]
    - [x] [T_0005_001] Deposit positive `amount`
        - Endpoint: [API-USER-DEP]
        - [x] Status: 200
        - [x] Result: `ledger.entry_type` == `"credit"`
        - [x] Result: `ledger.amount` == `amount`
    - [x] [T_0005_002] Get Balance
        - Endpoint: [API-USER-BAL]
        - [x] Status: 200
        - [x] Result: `after.balance` = `ledger.balance`
        - [x] Assert: `before.balance` + amount == `after.balance`
    - [x] [T_0005_003] Get History
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] Result: Assert `ledgers`= [`deposit`]
- [x] [T_0006] - New User: Withdraw Wallet Fail (Insufficient Funds)\
  User Stories: [US-001], [US-002]
    - [x] [Setup] Do [T_0003]
    - [x] [T_0006_001] Withdraw from new wallet `wdr_amount`
        - Endpoint: [API-USER-DEP]
        - [x] Status: 400
        - [x] Result: Error Message = `insufficient_funds`
- [ ] [T_0006_001] - Withdraw Fail \
  User Stories: [US-002], [US-005]
    - [ ] New User `USER_ID`
        - [ ] 200
          Endpoint: [API-USER-NEW]
        - [ ] 200: Get `CURRENT_BALANCE` \
          Endpoint: [API-USER-BAL]
        - [ ] 400: Withdraw From Insufficient Balance \
          Endpoint: [API-WALL-WDR]
            - Withdraw `AMOUNT` = `CURRENT_BALANCE` + 1
            - [ ] `ERROR_MESSAGE` = `"INSUFFICIENT_FUNDS"`
        - [ ] 200: Get `CURRENT_BALANCE` \
          Endpoint: [API-USER-BAL]
            - [ ] Assert `NEW_BALANCE` = `CURRENT_BALANCE`
        - [ ] 200:Get history = [`WITHDRAW (FAIL)`]\
          Endpoint: [API-USER-TXH]

- [ ] [T_0005_002] - Withdraw Success \
  User Stories: [US-002], [US-005]
    - [ ] User `USER_ID`
        - [ ] Withdraw From Sufficient Balance \
          Endpoint: [API-USER-BAL]
            - Get `CURRENT_BALANCE`
            - Withdraw `AMOUNT` = `CURRENT_BALANCE`
            - [ ] Assert `NEW_BALANCE` = `CURRENT_BALANCE` - `AMOUNT`
              Endpoint: [API-USER-BAL], [API-WALL-WDR]
        - [ ] Get history = [`WITHDRAW (OK)`]
          Endpoint: [API-USER-TXH]
