# Test Plain (HTTP)

## Prerequisites

Ensure API server is ready.

## User Stories

[US-001] User can deposit money into his/her wallet\
[US-002] User can withdraw money from his/her wallet\
[US-003] User can send money to another user\
[US-004] User can check his/her wallet balance\
[US-005] User can view his/her transaction history

- [ ] [T-0001] - Get User Balance\
  User Stories: [US-004], [US-005]
    - [ ] Non-existing user `USER_ID`
        - [ ] Get balance of 0\
          Endpoints: [API-USER-BAL]
        - [ ] Get history = []\
          Endpoints: [API-USER-HST]
    - [ ] Existing user `USER_ID` (same as above)
        - [ ] Get balance of 0\
          Endpoints: [API-USER-BAL]
        - [ ] Get history = []\
          Endpoints: [API-USER-HST]

- [ ] [T-0002] - Deposit\
  User Stories: [US-001], [US-005]
    - [ ] New User `USER_ID`
        - [ ] Deposit\
          Endpoints: [API-USER-BAL], [API-USER-DEP]
            - Get `CURRENT_BALANCE`
            - Deposit `AMOUNT`
            - [ ] Assert `NEW_BALANCE` = `CURRENT_BALANCE` + `AMOUNT`
        - [ ] Get history = [`Deposit`]\
          Endpoints: [API-USER-HST]

- [ ] [T-0003-001] - Withdraw Fail \
  User Stories: [US-002], [US-005]
    - [ ] User `USER_ID`
        - [ ] Withdraw From Insufficient Balance \
          Endpoints: [API-USER-BAL], [API-USER-WDR]
            - Get `CURRENT_BALANCE`
            - Withdraw `AMOUNT` = `CURRENT_BALANCE` + 1
                - [ ] `ERROR_MESSAGE` = `INSUFFICIENT FUNDS`
            - [ ] Assert `NEW_BALANCE` = `CURRENT_BALANCE`
        - [ ] Get history = [`WITHDRAW (FAIL)`]\
          Endpoints: [API-USER-HST]

- [ ] [T-0003-002] - Withdraw Success \
  User Stories: [US-002], [US-005]
    - [ ] User `USER_ID`
        - [ ] Withdraw From Sufficient Balance \
          Endpoints: [API-USER-BAL], [API-USER-WDR]
            - Get `CURRENT_BALANCE`
            - Withdraw `AMOUNT` = `CURRENT_BALANCE`
            - [ ] Assert `NEW_BALANCE` = `CURRENT_BALANCE` - `AMOUNT`
        - [ ] Get history = [`WITHDRAW (OK)`]
          Endpoints: [API-USER-HST]
