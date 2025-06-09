package e2e_tests

import (
	"errors"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"

	testclient "github.com/cryptonlx/crypto/cmd/e2e_tests/client"

	"github.com/shopspring/decimal"
)

func TestE2E(t *testing.T) {
	serverUrl := os.Getenv("SERVER_URL")

	var n int
	_n := os.Getenv("N")
	if _n == "" {
		n = 1
	} else {
		_n, _err := strconv.Atoi(_n)
		n = _n
		if _err != nil {
			panic(_err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			client, err := testclient.NewClient(serverUrl)
			if err != nil {
				panic(err)
			}
			defer wg.Done()
			Exec(t, client)
		}()
	}
	wg.Wait()
}

func Exec(t *testing.T, client *testclient.Client) {
	Ping(t, client)
	T_0001(t, client)
	T_0002(t, client)
	T_0003(t, client)
	T_0004(t, client)
	T_0005(t, client)
	T_0006(t, client)
	T_0007(t, client)
	T_0008(t, client)
	T_0009(t, client)
	T_0010(t, client)
	T_0011(t, client)
}

func Ping(t *testing.T, client *testclient.Client) {
	_, _, err := client.Wallets("")
	if errors.Is(err, syscall.ECONNREFUSED) {
		t.Fatalf("Server not started...")
	}
}

func T_0001(t *testing.T, client *testclient.Client) {
	futureUserName := NewRandomUserName("t00001", 12, 2*24*time.Hour)
	responseBody, responseStatusCode, err := client.Wallets(futureUserName)
	if responseStatusCode != http.StatusBadRequest {
		t.Fatalf(`[T_0001_001] Wallets want status code 400. got code=%d, err=%#v`, responseStatusCode, err.Error())
	}
	if responseBody.Error == nil || *responseBody.Error != "resource: user not found" {
		t.Fatalf(`[T_0001_001] Wallets want Response.error="resource: user not found". got err=%v, Response.error=%v, user id=%d`, err, *responseBody.Error, 0)
	}

	username := NewRandomUserName("t00001", 12, 0)
	createUserResponseBody, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		t.Fatalf("[T_0001_002] CreateUser want 200. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}
	if err != nil {
		t.Fatalf("[T_0001_002] CreateUser want nil err. err=%v", err)
	}
	if createUserResponseBody.Error != nil {
		t.Fatalf("[T_0001_002] CreateUser want nil err. createUserResponseBody.Error=%v", *createUserResponseBody.Error)
	}

	if createUserResponseBody.Data.Id == 0 {
		t.Fatalf("[T_0001_002] CreateUser want id != 0. got id=%d", createUserResponseBody.Data.Id)
	}

	if createUserResponseBody.Data.Username != username {
		t.Fatalf("[T_0001_002] CreateUser want username=%s. createUserResponseBody.Data.Username=%s", username, createUserResponseBody.Data.Username)
	}

	createUserResponseBody, responseStatusCode, err = client.CreateUser(username)
	if responseStatusCode != http.StatusBadRequest {
		t.Fatalf("[T_0001_003] Duplicate CreateUser want 400. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}

	if createUserResponseBody.Error == nil || *createUserResponseBody.Error != "unique_violation" {
		t.Fatalf(`[T_0001_003] Duplicate CreateUser want "unique_violation". responseStatusCode=%s, err=%v`, *createUserResponseBody.Error, err)
	}

	responseBody, responseStatusCode, err = client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		t.Fatalf(`[T_0001_004] Wallets want 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
}

func T_0002(t *testing.T, client *testclient.Client) {
	futureUserName := NewRandomUserName("t00002", 12, 2*24*time.Hour)
	responseBody, responseStatusCode, err := client.Transactions(futureUserName)
	if responseStatusCode != http.StatusBadRequest {
		t.Fatalf(`[T_0002_001] Transactions() want status code 400. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error != nil && *responseBody.Error != "resource: user not found" {
		t.Fatalf(`[T_0002_001] Transactions() want Response.error "resource: user not found". got Response.error %s, user id = %d`, *responseBody.Error, 0)
	}

	username := NewRandomUserName("t00002", 12, 0)
	createUserResponseBody, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		t.Fatalf("[T_0002_002] CreateUser should succeed. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}
	if err != nil {
		t.Fatalf("[T_0002_002] CreateUser should succeed. err=%v", err)
	}
	if createUserResponseBody.Error != nil && *createUserResponseBody.Error != "" {
		t.Fatalf("[T_0002_002] CreateUser should succeed. createUserResponseBody.Error=%v", err)
	}

	if createUserResponseBody.Data.Id == 0 {
		t.Fatalf("[T_0002_002] CreateUser success but not return id. createUserResponseBody.Data.Id=%d", createUserResponseBody.Data.Id)
	}

	if createUserResponseBody.Data.Username != username {
		t.Fatalf("[T_0002_002] CreateUser success but not return username. createUserResponseBody.Data.Username=%s", createUserResponseBody.Data.Username)
	}

	responseBody, responseStatusCode, err = client.Transactions(username)
	if responseStatusCode != http.StatusOK {
		t.Fatalf(`[T_0002_003] Transactions() want 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, *responseBody.Error)
	}
}

func T_0003(t *testing.T, client *testclient.Client) {
	SetupUserAndWalletCreation(t, client, "T_0003", []string{"SGD"})
}

func T_0004(t *testing.T, client *testclient.Client) {
	username, wallets := SetupUserAndWalletCreation(t, client, "T_0004", []string{"SGD"})

	wallet := wallets[0]
	balanceBefore := wallet.Balance

	dRespBody, dStatusCode, err := client.Deposit(username, wallet.Id, decimal.NewFromInt(-40))
	if dStatusCode != http.StatusBadRequest {
		t.Fatalf("[T_0004_001] Deposit want err 400. responseStatusCode=%d, err=%v", dStatusCode, err)
	}
	if *(dRespBody.Error) != "invalid_amount" {
		t.Fatalf(`[T_0004_001] Deposit() want Response.error "invalid_amount". got err=%v, Response.error=%s %#v`, err, *dRespBody.Error, dRespBody)
	}

	responseBody, responseStatusCode, err := client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		t.Fatalf(`[T_0004_002] SETUP Wallets() want 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		t.Fatalf(`[T_0004_002] SETUP Wallets() want 200. Got responseBody.Error=%s, err=%#v`, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) == 0 {
		t.Fatalf(`[T_0004_002] SETUP Wallets() want some wallets. Got len=%d, err=%#v`, len(responseBody.Data.Wallets), err)
	}
	if responseBody.Data.Wallets[0].Balance != balanceBefore {
		t.Fatalf(`[T_0004_002] SETUP Wallets() balance before and after should be same. want=%s, got=%s`, balanceBefore, responseBody.Data.Wallets[0].Balance)
	}
}

func T_0005(t *testing.T, client *testclient.Client) {
	username, wallets := SetupUserAndWalletCreation(t, client, "T_0005", []string{"SGD"})

	wallet := wallets[0]

	dRespBody, dStatusCode, cErr := client.Deposit(username, wallet.Id, decimal.NewFromFloat(40.123))
	if cErr != nil {
		t.Fatalf(`[T_0005_001] Deposit transaction want nil err, got error %v`, cErr)
	}
	if dStatusCode != http.StatusOK {
		t.Fatalf("[T_0005_001] Deposit want 200. responseStatusCode=%d, err=%v", dStatusCode, cErr)
	}
	if dRespBody.Error != nil {
		var s string
		if dRespBody.Error != nil {
			s = *dRespBody.Error
		}
		t.Fatalf(`[T_0005_001] Deposit want nil Response.err. got Response.error=%v`, s)
	}
	if dRespBody.Data.Transaction.Id == 0 {
		t.Fatalf(`[T_0005_001] Deposit want non-zero transaction.Id.`)
	}
	if dRespBody.Data.Transaction.Ledgers[0].EntryType != "credit" {
		t.Fatalf(`[T_0005_001] Deposit want ledger.entry_type=%s got %s`, "credit", dRespBody.Data.Transaction.Ledgers[0].EntryType)
	}

	responseBody, responseStatusCode, err := client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		t.Fatalf(`[T_0005_002] Wallets want 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		t.Fatalf(`[T_0005_002] Wallets want 200. Got responseBody.Error=%s, err=%#v`, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) == 0 {
		t.Fatalf(`[T_0005_002] Wallets want some wallets. Got len=%d, err=%#v`, len(responseBody.Data.Wallets), err)
	}
	if responseBody.Data.Wallets[0].Balance != "40.123" {
		t.Fatalf(`[T_0005_002] Wallets balance before and after should be same. want=%s, got=%s`, "40.123", responseBody.Data.Wallets[0].Balance)
	}

	tRespBody, tStatusCode, err := client.Transactions(username)
	if tStatusCode != http.StatusOK {
		t.Fatalf(`[T_0005_003] Transactions want 200. Got responseStatusCode=%d, err=%#v`, tStatusCode, err)
	}

	if len(tRespBody.Data.Transactions) == 0 {
		t.Fatalf(`[T_0005_003] Transactions want transactions.len > 0. got 0`)
	}

	transaction0 := tRespBody.Data.Transactions[0]
	if len(transaction0.Ledgers) == 0 {
		t.Fatalf(`[T_0005_003] Transactions want transactions[0].ledgers[0].len > 0. got 0`)
	}
	if transaction0.Operation != "deposit" {
		t.Fatalf(`[T_0005_003] Transactions want transactions[0].ledgers[0].operation=deposit. got %s`, transaction0.Operation)
	}
	ledger := tRespBody.Data.Transactions[0].Ledgers[0]
	if ledger.EntryType != "credit" {
		t.Fatalf(`[T_0005_003] Transactions want transactions[0].ledgers[0].entry_type=credit. got %s`, ledger.EntryType)
	}
	if ledger.Amount != "40.123" {
		t.Fatalf(`[T_0005_003] Transactions want transactions[0].ledgers[0].amount="40.123". got %s`, ledger.Amount)
	}
}

func T_0006(t *testing.T, client *testclient.Client) {
	username, wallets := SetupUserAndWalletCreation(t, client, "T_0006", []string{"SGD"})

	wallet := wallets[0]
	dRespBody, dStatusCode, cErr := client.Withdraw(username, wallet.Id, decimal.NewFromFloat(50.1))
	if cErr != nil {
		t.Fatalf(`[T_0006_001] Withdraw transaction want nil err, got error %v`, cErr)
	}
	if dStatusCode != http.StatusBadRequest {
		t.Fatalf("[T_0006_001] Withdraw want 400. responseStatusCode=%d, err=%v", dStatusCode, cErr)
	}
	if dRespBody.Error == nil || *dRespBody.Error != "insufficient_funds" {
		var s string
		if dRespBody.Error != nil {
			s = *dRespBody.Error
		}
		t.Fatalf(`[T_0006_001] Withdraw want Response.err == "insufficient_funds". got Response.error=%v`, s)
	}
	if dRespBody.Data.Transaction.Id != 0 {
		t.Fatalf(`[T_0006_001] Withdraw want zero transaction.Id.`)
	}

	tRespBody, tStatusCode, err := client.Transactions(username)
	if tStatusCode != http.StatusOK {
		t.Fatalf(`[T_0006_002] Transactions want 200. Got responseStatusCode=%d, err=%#v`, tStatusCode, err)
	}
	if len(tRespBody.Data.Transactions) == 0 {
		t.Fatalf(`[T_0006_002] Transactions want transactions.len > 0. got 0`)
	}
	transaction0 := tRespBody.Data.Transactions[0]
	if len(transaction0.Ledgers) != 0 {
		t.Fatalf(`[T_0006_002] Transactions want transactions[0].ledgers[0].len == 0 (failed transaction should have 0 ledgers). got %d`, len(transaction0.Ledgers))
	}
	if transaction0.Operation != "withdraw" {
		t.Fatalf(`[T_0006_002] Transactions want transactions[0].ledgers[0].operation=withdraw. got %s`, transaction0.Operation)
	}
	if transaction0.Status != "error_insufficient_funds" {
		t.Fatalf(`[T_0006_002] Transactions want transactions[0].ledgers[0].status="error_insufficient_funds". got %s`, transaction0.Status)
	}
	if transaction0.MetaData.SourceWalletId == nil || *transaction0.MetaData.SourceWalletId != wallet.Id {
		t.Fatalf(`[T_0006_002] Transactions want transactions[0].metadata.source_wallet_id=%d. got %v`, wallet.Id, transaction0.MetaData.SourceWalletId)
	}
	if transaction0.MetaData.Amount == nil || *transaction0.MetaData.Amount != "50.1" {
		t.Fatalf(`[T_0006_002] Transactions want transactions[0].metadata.amount=%s. got %s`, "50.1", *transaction0.MetaData.Amount)
	}
}

func T_0007(t *testing.T, client *testclient.Client) {
	username, wallets := SetupUserAndWalletCreation(t, client, "T_0007", []string{"SGD"})

	wallet := wallets[0]

	_, dStatusCode, cErr := client.Deposit(username, wallet.Id, decimal.NewFromFloat(60.2))
	if cErr != nil {
		t.Fatalf(`[T_0007_001] Deposit transaction want nil err, got error %v`, cErr)
	}
	if dStatusCode != http.StatusOK {
		t.Fatalf("[T_0007_001] Deposit want 200. responseStatusCode=%d, err=%v", dStatusCode, cErr)
	}

	wRespBody, wStatusCode, cErr := client.Withdraw(username, wallet.Id, decimal.NewFromFloat(50.1))
	if cErr != nil {
		t.Fatalf(`[T_0007_002] Withdraw transaction want nil err, got error %v`, cErr)
	}
	if wStatusCode != http.StatusOK {
		t.Fatalf("[T_0007_002] Withdraw want 400. responseStatusCode=%d, err=%v", wStatusCode, cErr)
	}
	if wRespBody.Data.Transaction.Id == 0 {
		t.Fatalf(`[T_0007_002] Withdraw want transaction.Id. got 0.`)
	}

	responseBody, walStatusCode, cErr := client.Wallets(username)
	if cErr != nil {
		t.Fatalf("[T_0007_003] Wallets want nil err. responseStatusCode=%d, err=%v", wStatusCode, cErr)
	}
	if walStatusCode != http.StatusOK {
		t.Fatalf("[T_0007_003] Wallets want 400. responseStatusCode=%d, err=%v", wStatusCode, cErr)
	}

	walletAfterTx := responseBody.Data.Wallets[0]
	if walletAfterTx.Balance != "10.1" {
		t.Fatalf(`[T_0007_003] Wallets after transactions want balance=%s. got %s.`, "10.1", walletAfterTx.Balance)
	}

	tRespBody, tStatusCode, cErr := client.Transactions(username)
	if cErr != nil {
		t.Fatalf("[T_0007_004] Transactions want nil err. responseStatusCode=%d, err=%v", tStatusCode, cErr)
	}
	if tStatusCode != http.StatusOK {
		t.Fatalf(`[T_0007_004] Transactions want 200. Got responseStatusCode=%d, err=%#v`, tStatusCode, cErr)
	}
	if len(tRespBody.Data.Transactions) != 2 {
		t.Fatalf(`[T_0007_004] Transactions want transactions.len = 2. got %d`, len(tRespBody.Data.Transactions))
	}
	if tRespBody.Data.Transactions[0].Operation != "withdraw" {
		t.Fatalf(`[T_0007_004] Transactions want transactions[0].operation="withdraw". got %s`, tRespBody.Data.Transactions[0].Operation)
	}
	if tRespBody.Data.Transactions[1].Operation != "deposit" {
		t.Fatalf(`[T_0007_004] Transactions want transactions[0].operation="deposit". got %s`, tRespBody.Data.Transactions[0].Operation)
	}
}

func T_0008(t *testing.T, client *testclient.Client) {
	username0, user0Wallets := SetupUserAndWalletCreation(t, client, "T_0008", []string{"SGD"})
	_, user1Wallets := SetupUserAndWalletCreation(t, client, "T_0008", []string{"USD"})

	user0wallet0 := user0Wallets[0]
	_, dStatusCode, cErr := client.Deposit(username0, user0wallet0.Id, decimal.NewFromFloat(60.2))
	if cErr != nil {
		t.Fatalf(`[T_0008_001] Deposit transaction want nil err, got error %v`, cErr)
	}
	if dStatusCode != http.StatusOK {
		t.Fatalf("[T_0008_001] Deposit want 200. responseStatusCode=%d, err=%v", dStatusCode, cErr)
	}

	user1wallet0 := user1Wallets[0]
	wRespBody, wStatusCode, cErr := client.Transfer(username0, user0wallet0.Id, user1wallet0.Id, decimal.NewFromFloat(60.2))
	if cErr != nil {
		t.Fatalf(`[T_0008_002] Transfer transaction want nil err, got error %v`, cErr)
	}
	if wStatusCode != http.StatusBadRequest {
		t.Fatalf("[T_0008_002] Transfer want 400. responseStatusCode=%d, err=%s", wStatusCode, cErr.Error())
	}

	if wRespBody.Error == nil {
		t.Fatalf(`[T_0008_002] Transfer want non-nil responseBody.Error. got nil`)
	}
	if *wRespBody.Error != "currency_mismatch" {
		t.Fatalf(`[T_0008_002] Transfer want responseBody.Error="currency_mismatch". got %s`, *wRespBody.Error)
	}

	tRespBody, tStatusCode, cErr := client.Transactions(username0)
	if cErr != nil {
		t.Fatalf("[T_0008_003] Transactions want nil err. responseStatusCode=%d, err=%v", tStatusCode, cErr)
	}
	if tStatusCode != http.StatusOK {
		t.Fatalf(`[T_0008_003] Transactions want 200. Got responseStatusCode=%d, err=%#v`, tStatusCode, cErr)
	}
	if len(tRespBody.Data.Transactions) != 2 {
		t.Fatalf(`[T_0008_003] Transactions want transactions.len = 2. got %d`, len(tRespBody.Data.Transactions))
	}
	if tRespBody.Data.Transactions[0].Operation != "transfer" {
		t.Fatalf(`[T_0008_003] Transactions want transactions[0].operation="transfer". got %s`, tRespBody.Data.Transactions[0].Operation)
	}
	if tRespBody.Data.Transactions[0].Status != "error_currency_mismatch" {
		t.Fatalf(`[T_0008_003] Transactions want transactions[0].status="error_currency_mismatch". got %s`, tRespBody.Data.Transactions[0].Status)
	}
	if tRespBody.Data.Transactions[1].Operation != "deposit" {
		t.Fatalf(`[T_0008_003] Transactions want transactions[1].operation="deposit". got %s`, tRespBody.Data.Transactions[0].Operation)
	}
	if tRespBody.Data.Transactions[1].Status != "success" {
		t.Fatalf(`[T_0008_003] Transactions want transactions[1].status="success". got %s`, tRespBody.Data.Transactions[0].Status)
	}
}

func T_0009(t *testing.T, client *testclient.Client) {
	username0, user0Wallets := SetupUserAndWalletCreation(t, client, "T_0009", []string{"SGD"})
	_, user1Wallets := SetupUserAndWalletCreation(t, client, "T_0009", []string{"SGD"})

	user0wallet0 := user0Wallets[0]
	_, dStatusCode, cErr := client.Deposit(username0, user0wallet0.Id, decimal.NewFromFloat(60.2))
	if cErr != nil {
		t.Fatalf(`[T_0009_001] Deposit transaction want nil err, got error %v`, cErr)
	}
	if dStatusCode != http.StatusOK {
		t.Fatalf("[T_0009_001] Deposit want 200. responseStatusCode=%d, err=%v", dStatusCode, cErr)
	}

	user1wallet0 := user1Wallets[0]
	wRespBody, wStatusCode, cErr := client.Transfer(username0, user0wallet0.Id, user1wallet0.Id, decimal.NewFromFloat(90.2))
	if cErr != nil {
		t.Fatalf(`[T_0009_002] Transfer transaction want nil cErr, got error %v`, cErr)
	}
	if wStatusCode != http.StatusBadRequest {
		t.Fatalf("[T_0009_002] Transfer want 400. responseStatusCode=%d, err=%v", wStatusCode, cErr)
	}

	if wRespBody.Error == nil {
		t.Fatalf(`[T_0009_002] Transfer want non-nil responseBody.Error. got nil`)
	}
	if *wRespBody.Error != "insufficient_funds" {
		t.Fatalf(`[T_0009_002] Transfer want responseBody.Error="currency_mismatch". got %s`, *wRespBody.Error)
	}

	tRespBody, tStatusCode, cErr := client.Transactions(username0)
	if cErr != nil {
		t.Fatalf("[T_0009_003] Transactions want nil err. responseStatusCode=%d, err=%v", tStatusCode, cErr)
	}
	if tStatusCode != http.StatusOK {
		t.Fatalf(`[T_0009_003] Transactions want 200. Got responseStatusCode=%d, err=%#v`, tStatusCode, cErr)
	}
	if len(tRespBody.Data.Transactions) != 2 {
		t.Fatalf(`[T_0009_003] Transactions want transactions.len = 2. got %d`, len(tRespBody.Data.Transactions))
	}
	if tRespBody.Data.Transactions[0].Operation != "transfer" {
		t.Fatalf(`[T_0009_003] Transactions want transactions[0].operation="transfer". got %s`, tRespBody.Data.Transactions[0].Operation)
	}
	if tRespBody.Data.Transactions[0].Status != "error_insufficient_funds" {
		t.Fatalf(`[T_0009_003] Transactions want transactions[0].status="error_insufficient_funds". got %s`, tRespBody.Data.Transactions[0].Status)
	}
	if tRespBody.Data.Transactions[1].Operation != "deposit" {
		t.Fatalf(`[T_0009_003] Transactions want transactions[1].operation="deposit". got %s`, tRespBody.Data.Transactions[1].Operation)
	}
	if tRespBody.Data.Transactions[1].Status != "success" {
		t.Fatalf(`[T_0009_003] Transactions want transactions[1].status="success". got %s`, tRespBody.Data.Transactions[1].Status)
	}
}

func T_0010(t *testing.T, client *testclient.Client) {
	username0, user0Wallets := SetupUserAndWalletCreation(t, client, "T_0010", []string{"SGD"})
	username1, user1Wallets := SetupUserAndWalletCreation(t, client, "T_0010", []string{"SGD"})

	user0wallet0 := user0Wallets[0]
	_, dStatusCode, cErr := client.Deposit(username0, user0wallet0.Id, decimal.NewFromFloat(60.2))
	if cErr != nil {
		t.Fatalf(`[T_0010_001] Deposit transaction want nil err, got error %v`, cErr)
	}
	if dStatusCode != http.StatusOK {
		t.Fatalf("[T_0010_001] Deposit want 200. responseStatusCode=%d, err=%v", dStatusCode, cErr)
	}

	user1wallet0 := user1Wallets[0]
	wRespBody, wStatusCode, cErr := client.Transfer(username0, user0wallet0.Id, user1wallet0.Id, decimal.NewFromFloat(60.2))
	if cErr != nil {
		t.Fatalf(`[T_0010_002] Transfer transaction want nil cErr, got error %v`, cErr)
	}
	if wStatusCode != http.StatusOK {
		t.Fatalf("[T_0010_002] Transfer want 200. responseStatusCode=%d, err=%v", wStatusCode, cErr)
	}

	if wRespBody.Error != nil {
		t.Fatalf(`[T_0010_002] Transfer want nil responseBody.Error. got %s`, *wRespBody.Error)
	}

	tRespBody, tStatusCode, cErr := client.Transactions(username0)
	if cErr != nil {
		t.Fatalf("[T_0010_003] Transactions want nil err. responseStatusCode=%d, err=%v", tStatusCode, cErr)
	}
	if tStatusCode != http.StatusOK {
		t.Fatalf(`[T_0010_003] Transactions want 200. Got responseStatusCode=%d, err=%#v`, tStatusCode, cErr)
	}
	if len(tRespBody.Data.Transactions) != 2 {
		t.Fatalf(`[T_0010_003] Transactions want transactions.len = 2 (transfer+deposit). got %d`, len(tRespBody.Data.Transactions))
	}
	if tRespBody.Data.Transactions[0].Operation != "transfer" {
		t.Fatalf(`[T_0010_003] Transactions want transactions[0].operation="transfer". got %s`, tRespBody.Data.Transactions[0].Operation)
	}
	if tRespBody.Data.Transactions[0].Status != "success" {
		t.Fatalf(`[T_0010_003] Transactions want transactions[0].status="success". got %s`, tRespBody.Data.Transactions[0].Status)
	}
	if len(tRespBody.Data.Transactions[0].Ledgers) != 1 {
		t.Fatalf(`[T_0010_003] Transactions want transactions[0].ledgers.len=1. got %d`, len(tRespBody.Data.Transactions[0].Ledgers))
	}

	if tRespBody.Data.Transactions[1].Operation != "deposit" {
		t.Fatalf(`[T_0010_003] Transactions want transactions[1].operation="deposit". got %s`, tRespBody.Data.Transactions[1].Operation)
	}
	if tRespBody.Data.Transactions[1].Status != "success" {
		t.Fatalf(`[T_0010_003] Transactions want transactions[1].status="success". got %s`, tRespBody.Data.Transactions[1].Status)
	}
	if len(tRespBody.Data.Transactions[1].Ledgers) != 1 {
		t.Fatalf(`[T_0010_003] Transactions want transactions[0].ledgers.len=1. got %d`, len(tRespBody.Data.Transactions[0].Ledgers))
	}

	tRespBody, tStatusCode, cErr = client.Transactions(username1)
	if cErr != nil {
		t.Fatalf("[T_0010_004] Transactions want nil err. responseStatusCode=%d, err=%v", tStatusCode, cErr)
	}
	if tStatusCode != http.StatusOK {
		t.Fatalf(`[T_0010_004] Transactions want 200. Got responseStatusCode=%d, err=%#v`, tStatusCode, cErr)
	}
	if len(tRespBody.Data.Transactions) != 1 {
		t.Fatalf(`[T_0010_004] Transactions want transactions.len = 1 (transfer by others). got %d`, len(tRespBody.Data.Transactions))
	}
	if tRespBody.Data.Transactions[0].Operation != "transfer" {
		t.Fatalf(`[T_0010_004] Transactions want transactions[0].operation="transfer". got %s`, tRespBody.Data.Transactions[0].Operation)
	}
	if tRespBody.Data.Transactions[0].Status != "success" {
		t.Fatalf(`[T_0010_004] Transactions want transactions[0].status="success". got %s`, tRespBody.Data.Transactions[0].Status)
	}
	if len(tRespBody.Data.Transactions[0].Ledgers) != 1 {
		t.Fatalf(`[T_0010_004] Transactions want transactions[0].ledgers.len=1 (transfer by others). got %d`, len(tRespBody.Data.Transactions[0].Ledgers))
	}

	wallRespBody, wallStatusCode, cErr := client.Wallets(username0)
	if cErr != nil {
		t.Fatalf("[T_0010_005] Wallets want nil err. responseStatusCode=%d, err=%v", wallStatusCode, cErr)
	}
	if wallStatusCode != http.StatusOK {
		t.Fatalf(`[T_0010_005] Wallets want 200. Got responseStatusCode=%d, err=%#v`, wallStatusCode, cErr)
	}
	wallet := wallRespBody.Data.Wallets[0]
	if wallet.Balance != "0" {
		t.Fatalf("[T_0010_005] Wallets want balance=0. got=%s", wallet.Balance)
	}

	wallRespBody, wallStatusCode, cErr = client.Wallets(username1)
	if cErr != nil {
		t.Fatalf("[T_0010_006] Wallets want nil err. responseStatusCode=%d, err=%v", wallStatusCode, cErr)
	}
	if wallStatusCode != http.StatusOK {
		t.Fatalf(`[T_0010_006] Wallets want 200. Got responseStatusCode=%d, err=%#v`, wallStatusCode, cErr)
	}
	wallet = wallRespBody.Data.Wallets[0]
	if wallet.Balance != "60.2" {
		t.Fatalf("[T_0010_006] Wallets want balance=0. got=%s", wallet.Balance)
	}
}

func T_0011(t *testing.T, client *testclient.Client) {
	SetupUserAndWalletCreation(t, client, "T_0011", []string{"SGD", "USD", "MYD"})
}

func SetupUserAndWalletCreation(t *testing.T, client *testclient.Client, logPrefix string, currencies []string) (username string, wallets []testclient.Wallet) {
	username = NewRandomUserName(logPrefix, 12, 0)
	createUserResponseData, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		t.Fatalf("[%s_001] SETUP CreateUser should succeed. responseStatusCode=%d, err=%v", logPrefix, responseStatusCode, err)
	}
	if err != nil {
		t.Fatalf("[%s_001] SETUP CreateUser should succeed. err=%s", logPrefix, err.Error())
	}
	if createUserResponseData.Error != nil && *createUserResponseData.Error != "" {
		t.Fatalf("[%s_001] SETUP CreateUser should succeed. createUserResponseData.Error=%v", logPrefix, err)
	}
	if createUserResponseData.Data.Id == 0 {
		t.Fatalf("[%s_001] SETUP CreateUser success but missing return id. createUserResponseData.Data.Id=%d", logPrefix, createUserResponseData.Data.Id)
	}
	if createUserResponseData.Data.Username != username {
		t.Fatalf("[%s_001] SETUP CreateUser success but missing return username. createUserResponseData.Data.Username=%s", logPrefix, createUserResponseData.Data.Username)
	}

	responseBody, responseStatusCode, err := client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		t.Fatalf(`[%s_002] SETUP Wallets want 200. Got responseStatusCode=%d, err=%#v`, logPrefix, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		t.Fatalf(`[%s_002] SETUP Wallets want 200. Got responseBody.Error=%s, err=%#v`, logPrefix, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) != 0 {
		t.Fatalf(`[%s_002] SETUP Wallets want no wallets. Got len=%d, err=%#v`, logPrefix, len(responseBody.Data.Wallets), err)
	}

	for _, currency := range currencies {
		createWalletResponseBody, createWalletResponseStatusCode, err := client.CreateWallet(username, currency)
		if createWalletResponseStatusCode != http.StatusOK {
			t.Fatalf(`[%s_003] SETUP CreateWallet want 200. Got body.Err=%s, statusCode=%d, err=%#v`, logPrefix, *createWalletResponseBody.Error, createWalletResponseStatusCode, err)
		}
		if createWalletResponseBody.Error != nil {
			t.Fatalf(`[%s_003] SETUP CreateWallet want 200. Got createWalletResponseBody.Error=%s, err=%#v`, logPrefix, *createWalletResponseBody.Error, err)
		}
		wallet := createWalletResponseBody.Data.Wallet
		if wallet.Id == 0 {
			t.Fatalf(`[%s_003] SETUP CreateWallet should succeed. Got id=%v, err=%#v`, logPrefix, createWalletResponseBody.Data.Wallet, err)
		}
	}

	responseBody, responseStatusCode, err = client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		t.Fatalf(`[%s_004] SETUP Wallets want 200. Got responseStatusCode=%d, err=%#v`, logPrefix, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		t.Fatalf(`[%s_004] SETUP Wallets want 200. Got responseBody.Error=%s, err=%#v`, logPrefix, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) != len(currencies) {
		t.Fatalf(`[%s_004] SETUP Wallets want some wallets. Got len=%d, err=%#v`, logPrefix, len(responseBody.Data.Wallets), err)
	}

	gotWallets := responseBody.Data.Wallets

	slices.Sort(currencies)
	slices.SortFunc(gotWallets, func(a, b testclient.Wallet) int {
		if a.Currency < b.Currency {
			return -1
		}
		if a.Currency > b.Currency {
			return 1
		}
		return 0
	})

	for i, _ := range gotWallets {
		if gotWallets[i].Currency != currencies[i] {
			t.Fatalf(`[%s_004] SETUP Wallets want currency=%s. Got currency=%s, err=%#v`, logPrefix, currencies, responseBody.Data.Wallets[0].Currency, i)
		}
	}
	return username, responseBody.Data.Wallets
}
