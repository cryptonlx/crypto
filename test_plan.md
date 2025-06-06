# Test Plain (HTTP)

## Prerequisites

Ensure API test server is ready.

## User Stories

[US-001] User can deposit money into his/her wallet\
[US-002] User can withdraw money from his/her wallet\
[US-003] User can send money to another user\
[US-004] User can check his/her wallet balance\
[US-005] User can view his/her transaction history

- [ ] [T-0001] User Creation
    - [ ] Get user `USER_ID`
        - Endpoint: [API-USER-BAL]
        - [ ] Error: `ERROR_MESSAGE` = `"USER_NOT_FOUND"`
        - [ ] Status: 400

    - [ ] Create user `USER_ID` (same as above)
        - Endpoint: [API-USER-NEW]
        - [ ] Status: 200
    - [ ] Create duplicate user `USER_ID` (same as above)
        - Endpoint: [API-USER-NEW]
        - [ ] Error: `ERROR_MESSAGE` = `"USER_EXISTS"`
        - [ ] Status: 400
    - [ ] Get user `USER_ID` (same as above)\
        - Endpoint: [API-USER-BAL]
        - Status: 200

- [ ] [T-0002] - Get User Balance\
  User Stories: [US-004], [US-005]
    - [ ] Non-existing user `USER_ID`
        - [ ] 404: `ERROR_MESSAGE` = `"USER NOT FOUND"`\
          Endpoint: [API-USER-BAL]
        - [ ] 404: `ERROR_MESSAGE` = `"USER NOT FOUND"`\
          Endpoint: [API-USER-HST]
    - [ ] Create user `USER_ID` (same as above)
        - [ ] 200
          Endpoint: [API-USER-NEW]
        - [ ] 404: Get balance of 0\
          Endpoint: [API-USER-BAL]
        - [ ] 200: Get history = []\
          Endpoint: [API-USER-HST]

- [ ] [T-0003] - Deposit\
  User Stories: [US-001], [US-005]
    - [ ] New User `USER_ID`
        - [ ] 200
          Endpoint: [API-USER-NEW]
        - [ ] 400 Deposit `AMOUNT` = negative value\
          Endpoint: [API-WALL-DEP]
            - Deposit `AMOUNT` = negative value
            - [ ] Assert `NEW_BALANCE` = `CURRENT_BALANCE` + `AMOUNT`
        - [ ] Get history = [`Deposit`]\
          Endpoint: [API-USER-HST]

- [ ] [T-0004-001] - Withdraw Fail \
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

- [ ] [T-0004-002] - Withdraw Success \
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
