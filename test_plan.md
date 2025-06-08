# Test Plan (HTTP)

## Prerequisites

Ensure API test server is [ready](./readme.md#setup-local-environment) and [execute](./readme.md#run-e2e-tests) tests.

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

- [x] [T_0002] - New User Transaction History\
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
- [x] [T_0003] - Wallet Creation\
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
- [x] [T_0004] - New User: Deposit Error on Pre-Validation (Negative Amount)\
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
        - [x] Result: Assert `ledgers`= [`deposit.status=success`]
- [x] [T_0006] - New User: Withdraw Wallet Fail (Insufficient Funds)\
  User Stories: [US-002], [US-005]
    - [x] [Setup] Do [T_0003]
    - [x] [T_0006_001] Withdraw from new wallet `wdr_amount`
        - Endpoint: [API-USER-WDR]
        - [x] Status: 400
        - [x] Result: Error Message = `insufficient_funds`
    - [x] [T_0006_002] Get History
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] Result: Assert `ledgers`= [`withdraw.status=error, metadata.source_wallet_id, metadata.amount`]
- [x] [T_0007] - New User: Withdraw Wallet Success\
  User Stories: [US-001], [US-002], [US-004], [US-005]
    - [x] [Setup] Do [T_0003]
    - [x] [T_0007_001] Deposit `deposit_amount`=60.2
        - Endpoint: [API-USER-DEP]
        - [x] Status: 200
    - [x] [T_0007_002] Withdraw `wdr_amount`=50.1
        - Endpoint: [API-USER-WDR]
        - [x] Status: 200
    - [x] [T_0007_003] Get Balance
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] Result: `wallet.balance`=10.1
    - [x] [T_0007_004] Get History
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] Result: Assert in order: `ledgers`= [`withdraw.status=success`, `deposit.status=success`]
- [x] [T_0008] - Transfer Fail (Currency Mismatch)\
  User Stories: [US-001], [US-002], [US-003], [US-004], [US-005]
    - [x] [Setup]
        - [x] get `user1.wallet` <- Do [T_0003] curr=SGD
        - [x] get `user2.wallet` <- Do [T_0003] curr=USD
    - [x] [T_0008_001] Deposit to `user1.wallet`. `amount`=60.2
        - Endpoint: [API-USER-DEP]
        - [x] Status: 200
    - [x] [T_0008_002] Transfer `amount` to `user2.wallet`
        - Endpoint: [API-USER-TRF]
        - [x] Status: 400
        - [x] Error Message = `currency_mismatch`
    - [x] [T_0008_003] Get `user1` History
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] `ledgers` = [`transfer.status=error_currency_mismatch`, `deposit.status=success`]
- [x] [T_0009] - Transfer Fail (Insufficient Funds)\
  User Stories: [US-001], [US-002], [US-003], [US-004], [US-005]
    - [x] [Setup]
        - [x] get `user1.wallet` <- Do [T_0003] curr=SGD
        - [x] get `user2.wallet` <- Do [T_0003] curr=SGD
    - [x] [T_0009_001] Deposit to `user1.wallet`. `amount`=60.2
        - Endpoint: [API-USER-DEP]
        - [x] Status: 200
    - [x] [T_0009_002] Transfer `transfer_amount`=90.2 to `user2.wallet`
        - Endpoint: [API-USER-TRF]
        - [x] Status: 400
        - [x] Error Message = `currency_mismatch`
    - [x] [T_0009_003] Get `user1` History
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] `ledgers` = [`transfer.status=error_currency_mismatch`, `deposit.status=success`]
- [x] [T_0010] - Transfer Success\
  User Stories: [US-001], [US-002], [US-003], [US-004], [US-005]
    - [x] [Setup]
        - [x] get `user1.wallet` <- Do [T_0003] curr=SGD
        - [x] get `user2.wallet` <- Do [T_0003] curr=SGD
    - [x] [T_0010_001] Deposit to `user1.wallet`. `amount`=60.2
        - Endpoint: [API-USER-DEP]
        - [x] Status: 200
    - [x] [T_0010_002] Transfer `amount` to `user2.wallet`
        - Endpoint: [API-USER-TRF]
        - [x] Status: 200
    - [x] [T_0010_003] Get `user1` history
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] Result: Assert in order: `ledgers`= []
    - [x] [T_0010_004] Get `user2` history
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] Result: Assert in order: `ledgers`= [`transfer.status=success`]
    - [x] [T_0010_005] Get `user1.wallet` balance
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] Result: `wallet.balance`=0
    - [x] [T_0010_006] Get `user2.wallet` balance
        - Endpoint: [API-USER-TXH]
        - [x] Status: 200
        - [x] Result: `wallet.balance`=60.2