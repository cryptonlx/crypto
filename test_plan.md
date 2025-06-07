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
        - Endpoint: [API-WALL-BAL]
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
        - Endpoint: [API-WALL-BAL]
        - [x] Status: 200

- [x] [T_0002] - Get New User Transaction History\
  User Stories: [US-005]
    - [x] [T_0002_001] Get history by user yet to be created.
        - Endpoint: [API-WALL-HST]
        - [x] Status: 400
        - [x] Error Message = `"resource: user not found"`
    - [x] [T_0002_002] Create user with `username`
        - Endpoint: [API-USER-NEW]
        - [x] Status: 200
    - [x] [T_0002_003] Get History
        - Endpoint: [API-WALL-HST]
        - [x] Status: 200
        - [x] Result: []


- [x] [T_0003] - Wallet Creation \
  User Stories: [US-004]
    - [x] [T_0003_001] Create user with `username`
        - Endpoint: [API-USER-NEW]
        - [x] Status: 200
    - [x] [T_0003_002] Get Balance
        - Endpoint: [API-WALL-BAL]
        - [x] Status: 200
        - [x] Result: []
    - [x] [T_0003_003] Create `SGD` wallet
        - [x] Status: 200
        - [x] Result: `response.wallet.id` != 0
    - [x] [T_0003_004] Get Balance
        - Endpoint: [API-WALL-BAL]
        - [x] Status: 200
        - [x] Result: `response.wallets.len` == 1
        - [x] Result: `response.wallets[0].currency` == `SGD`
- [x] [T_0004] - New User: Deposit Error\
  User Stories: [US-001], [US-005]
    - [x] [Setup] Do [T_0003]
    - [x] [T_0004_001] Deposit `amount` less than or equals to 0 should fail
        - Endpoint: [API-USER-DEP]
        - [x] Status: 400
        - [x] Error Message = `"invalid_amount"`
    - [x] [T_0004_002] Get Balance
        - Endpoint: [API-WALL-BAL]
        - [x] Status: 200
        - [x] Result: `after.balance`
        - [x] Assert: `before.balance` == `after.balance`
- [ ] [T_0005] - New User: Deposit Success\
  User Stories: [US-001], [US-005]
    - [ ] [Setup] Do [T_0003]
    - [ ] [T_0005_001] Deposit positive `amount`
        - Endpoint: [API-USER-DEP]
        - [ ] Status: 200
        - [ ] Result: `ledger.entry_type` == `"credit"`
        - [ ] Result: `ledger.amount` == `amount`
    - [ ] [T_0005_002] Get Balance
        - Endpoint: [API-WALL-BAL]
        - [ ] Status: 200
        - [ ] Result: `after.balance`
        - [ ] Assert: `before.balance` + amount == `after.balance`

- [ ] [T_0005] - New User: Deposit Success\
  User Stories: [US-001], [US-005]
    - [ ] [T_0004_001] Create user with `username`
        - Endpoint: [API-USER-NEW]
        - [ ] Status: 200
    - [ ] [T_0004_002] Get History
        - Endpoint: [API-WALL-HST]
        - [ ] Status: 200
        - [ ] Result: []
    - [ ] [T_0004_002] Get Balance
        - Endpoint: [API-WALL-BAL]
        - [ ] Status: 200
        - [ ] Result: wallets length > 1
        - [ ] Result: `target_wallet`=`wallets[0]`
        - [ ] Result: `before.amount`=`wallets[0].amount`
    - [ ] [T_0004_004] Deposit `amount` less than or equals to 0 should fail
        - Endpoint: [API-USER-DEP]
        - [ ] Status: 400
        - [ ] Error Message = `"invalid_amount"`
    - [ ] [T_0004_005] Get Balance
        - Endpoint: [API-WALL-BAL]
        - [ ] Status: 200
        - [ ] Result: `after.balance`
        - [ ] Assert: `before.balance` == `after.balance`
    - [ ] [T_0004_006] Get History
        - Endpoint: [API-WALL-HST]
        - [ ] Status: 200
        - [ ] Result: [`transaction::deposit::error`]
- [ ] [T_0005_001] - Withdraw Fail \
  User Stories: [US-002], [US-005]
    - [ ] New User `USER_ID`
        - [ ] 200
          Endpoint: [API-USER-NEW]
        - [ ] 200: Get `CURRENT_BALANCE` \
          Endpoint: [API-WALL-BAL]
        - [ ] 400: Withdraw From Insufficient Balance \
          Endpoint: [API-WALL-WDR]
            - Withdraw `AMOUNT` = `CURRENT_BALANCE` + 1
            - [ ] `ERROR_MESSAGE` = `"INSUFFICIENT_FUNDS"`
        - [ ] 200: Get `CURRENT_BALANCE` \
          Endpoint: [API-WALL-BAL]
            - [ ] Assert `NEW_BALANCE` = `CURRENT_BALANCE`
        - [ ] 200:Get history = [`WITHDRAW (FAIL)`]\
          Endpoint: [API-WALL-HST]

- [ ] [T_0005_002] - Withdraw Success \
  User Stories: [US-002], [US-005]
    - [ ] User `USER_ID`
        - [ ] Withdraw From Sufficient Balance \
          Endpoint: [API-WALL-BAL]
            - Get `CURRENT_BALANCE`
            - Withdraw `AMOUNT` = `CURRENT_BALANCE`
            - [ ] Assert `NEW_BALANCE` = `CURRENT_BALANCE` - `AMOUNT`
              Endpoint: [API-WALL-BAL], [API-WALL-WDR]
        - [ ] Get history = [`WITHDRAW (OK)`]
          Endpoint: [API-WALL-HST]
