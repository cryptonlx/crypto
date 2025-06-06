# Test Plain (HTTP)

## Prerequisites

Ensure API test server is ready.

## User Stories

[US-001] User can deposit money into his/her wallet\
[US-002] User can withdraw money from his/her wallet\
[US-003] User can send money to another user\
[US-004] User can check his/her wallet balance\
[US-005] User can view his/her transaction history

- [ ] [T_0001] User Creation
  User Stories: [US-004]
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

- [ ] [T_0002] - Get User History\
  User Stories: [US-005]
    - [ ] [T_0002_001] Create user with `username`
        - Endpoint: [API-USER-NEW]
        - [ ] Status: 200
    - [ ] [T_0002_002] Get History
        - Endpoint: [API-USER-HST]
        - [ ] Status: 200
        - [ ] Result: []


- [ ] [T_0003] - Deposit Negative Amount Should Fail\
  User Stories: [US-001], [US-005]
    - [ ] [T_0003_001] Create user with `username`
        - Endpoint: [API-USER-NEW]
        - [ ] Status: 200
    - [ ] [T_0003_002] Get History
        - Endpoint: [API-USER-HST]
        - [ ] Status: 200
        - [ ] Result: []
    - [ ] [T_0003_002] Get Balance
        - Endpoint: [API-USER-BAL]
        - [ ] Status: 200
        - [ ] Result: `before.balance`
    - [ ] [T_0003_003] Deposit `amount` less than or equals to 0
        - Endpoint: [API-USER-DEP]
        - [ ] Status: 400
        - [ ] Error Message = `"invalid_amount"`
    - [ ] [T_0003_002] Get Balance
        - Endpoint: [API-USER-BAL]
        - [ ] Status: 200
        - [ ] Result: `after.balance`
        - [ ] Assert: `before.balance` == `after.balance`
    - [ ] [T_0003_002] Get History
        - Endpoint: [API-USER-HST]
        - [ ] Status: 200
        - [ ] Result: [`transaction::deposit`]
- [ ] [T_0004_001] - Withdraw Fail \
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
          Endpoint: [API-USER-HST]

- [ ] [T_0004_002] - Withdraw Success \
  User Stories: [US-002], [US-005]
    - [ ] User `USER_ID`
        - [ ] Withdraw From Sufficient Balance \
          Endpoint: [API-USER-BAL]
            - Get `CURRENT_BALANCE`
            - Withdraw `AMOUNT` = `CURRENT_BALANCE`
            - [ ] Assert `NEW_BALANCE` = `CURRENT_BALANCE` - `AMOUNT`
              Endpoint: [API-USER-BAL], [API-WALL-WDR]
        - [ ] Get history = [`WITHDRAW (OK)`]
          Endpoint: [API-USER-HST]
