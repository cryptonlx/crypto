package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/cryptonlx/crypto/cmd/e2e_tests/client"

	"github.com/shopspring/decimal"
)

func main() {
	serverUrl := os.Getenv("SERVER_URL")
	client, err := client.NewClient(serverUrl)
	if err != nil {
		panic(err)
	}

	Ping(client)
	T_0001(client)
	T_0002(client)
	T_0003(client)
	T_0004(client)
	T_0005(client)
	T_0006(client)
}

func Ping(client *client.Client) {
	_, _, err := client.Wallets("")
	if errors.Is(err, syscall.ECONNREFUSED) {
		log.Fatalf("Server not started...")
	}
}

func T_0001(client *client.Client) {
	futureUserName := NewRandomUserName("t00001", 6, 2*24*time.Hour)
	responseBody, responseStatusCode, err := client.Wallets(futureUserName)
	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf(`[T_0001_001] Wallets want status code 400. got code=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error == nil || *responseBody.Error != "resource: user not found" {
		log.Fatalf(`[T_0001_001] Wallets want Response.error="resource: user not found". got err=%v, Response.error=%v, user id=%d`, err, *responseBody.Error, 0)
	}

	username := NewRandomUserName("t00001", 6, 0)
	createUserResponseBody, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[T_0001_002] CreateUser want 200. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[T_0001_002] CreateUser want nil err. err=%v", err)
	}
	if createUserResponseBody.Error != nil {
		log.Fatalf("[T_0001_002] CreateUser want nil err. createUserResponseBody.Error=%v", *createUserResponseBody.Error)
	}

	if createUserResponseBody.Data.Id == 0 {
		log.Fatalf("[T_0001_002] CreateUser want id != 0. got id=%d", createUserResponseBody.Data.Id)
	}

	if createUserResponseBody.Data.Username != username {
		log.Fatalf("[T_0001_002] CreateUser want username=%s. createUserResponseBody.Data.Username=%s", username, createUserResponseBody.Data.Username)
	}

	createUserResponseBody, responseStatusCode, err = client.CreateUser(username)
	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf("[T_0001_003] Duplicate CreateUser want 400. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}

	if createUserResponseBody.Error == nil || *createUserResponseBody.Error != "unique_violation" {
		log.Fatalf(`[T_0001_003] Duplicate CreateUser want "unique_violation". responseStatusCode=%s, err=%v`, *createUserResponseBody.Error, err)
	}

	responseBody, responseStatusCode, err = client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0001_004] Wallets want 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
}

func T_0002(client *client.Client) {
	futureUserName := NewRandomUserName("t00002", 6, 2*24*time.Hour)
	responseBody, responseStatusCode, err := client.Transactions(futureUserName)
	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf(`[T_0002_001] Transactions() want status code 400. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error != nil && *responseBody.Error != "resource: user not found" {
		log.Fatalf(`[T_0002_001] Transactions() want Response.error "resource: user not found". got Response.error %s, user id = %d`, *responseBody.Error, 0)
	}

	username := NewRandomUserName("t00002", 6, 0)
	createUserResponseBody, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[T_0002_002] CreateUser should succeed. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[T_0002_002] CreateUser should succeed. err=%v", err)
	}
	if createUserResponseBody.Error != nil && *createUserResponseBody.Error != "" {
		log.Fatalf("[T_0002_002] CreateUser should succeed. createUserResponseBody.Error=%v", err)
	}

	if createUserResponseBody.Data.Id == 0 {
		log.Fatalf("[T_0002_002] CreateUser success but not return id. createUserResponseBody.Data.Id=%d", createUserResponseBody.Data.Id)
	}

	if createUserResponseBody.Data.Username != username {
		log.Fatalf("[T_0002_002] CreateUser success but not return username. createUserResponseBody.Data.Username=%s", createUserResponseBody.Data.Username)
	}

	responseBody, responseStatusCode, err = client.Transactions(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0002_003] Transactions() want 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
}

func T_0003(client *client.Client) {
	SetupUserAndWalletCreation(client, "T_0003")
}

func T_0004(client *client.Client) {
	username, wallets := SetupUserAndWalletCreation(client, "T_0004")

	wallet := wallets[0]
	balanceBefore := wallet.Balance

	dRespBody, dStatusCode, err := client.Deposit(username, wallet.Id, decimal.NewFromInt(-40))
	if dStatusCode != http.StatusBadRequest {
		log.Fatalf("[T_0004_001] Deposit want err 400. responseStatusCode=%d, err=%v", dStatusCode, err)
	}
	if *(dRespBody.Error) != "invalid_amount" {
		log.Fatalf(`[T_0004_001] Deposit() want Response.error "invalid_amount". got err=%v, Response.error=%s %#v`, err, *dRespBody.Error, dRespBody)
	}

	responseBody, responseStatusCode, err := client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0004_002] SETUP Wallets() want 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		log.Fatalf(`[T_0004_002] SETUP Wallets() want 200. Got responseBody.Error=%s, err=%#v`, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) == 0 {
		log.Fatalf(`[T_0004_002] SETUP Wallets() want some wallets. Got len=%d, err=%#v`, len(responseBody.Data.Wallets), err)
	}
	if responseBody.Data.Wallets[0].Balance != balanceBefore {
		log.Fatalf(`[T_0004_002] SETUP Wallets() balance before and after should be same. want=%s, got=%s`, balanceBefore, responseBody.Data.Wallets[0].Balance)
	}
}

func T_0005(client *client.Client) {
	username, wallets := SetupUserAndWalletCreation(client, "T_0005")

	wallet := wallets[0]

	dRespBody, dStatusCode, cErr := client.Deposit(username, wallet.Id, decimal.NewFromFloat(40.123))
	if cErr != nil {
		log.Fatalf(`[T_0005_001] Deposit transaction want nil err, got error %v`, cErr)
	}
	if dStatusCode != http.StatusOK {
		log.Fatalf("[T_0005_001] Deposit want 200. responseStatusCode=%d, err=%v", dStatusCode, cErr)
	}
	if dRespBody.Error != nil {
		var s string
		if dRespBody.Error != nil {
			s = *dRespBody.Error
		}
		log.Fatalf(`[T_0005_001] Deposit want nil Response.err. got Response.error=%v`, s)
	}
	if dRespBody.Data.Transaction.Id == 0 {
		log.Fatalf(`[T_0005_001] Deposit want non-zero transaction.Id.`)
	}
	if dRespBody.Data.Transaction.Ledgers[0].EntryType != "credit" {
		log.Fatalf(`[T_0005_001] Deposit want ledger.entry_type=%s got %s`, "credit", dRespBody.Data.Transaction.Ledgers[0].EntryType)
	}

	responseBody, responseStatusCode, err := client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0005_002] Wallets want 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		log.Fatalf(`[T_0005_002] Wallets want 200. Got responseBody.Error=%s, err=%#v`, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) == 0 {
		log.Fatalf(`[T_0005_002] Wallets want some wallets. Got len=%d, err=%#v`, len(responseBody.Data.Wallets), err)
	}
	if responseBody.Data.Wallets[0].Balance != "40.123" {
		log.Fatalf(`[T_0005_002] Wallets balance before and after should be same. want=%s, got=%s`, "40.123", responseBody.Data.Wallets[0].Balance)
	}

	tRespBody, tStatusCode, err := client.Transactions(username)
	if tStatusCode != http.StatusOK {
		log.Fatalf(`[T_0005_003] Transactions want 200. Got responseStatusCode=%d, err=%#v`, tStatusCode, err)
	}

	if len(tRespBody.Data.Transactions) == 0 {
		log.Fatalf(`[T_0005_003] Transactions want transactions.len > 0. got 0`)
	}

	transaction0 := tRespBody.Data.Transactions[0]
	if len(transaction0.Ledgers) == 0 {
		log.Fatalf(`[T_0005_003] Transactions want transactions[0].ledgers[0].len > 0. got 0`)
	}
	if transaction0.Operation != "deposit" {
		log.Fatalf(`[T_0005_003] Transactions want transactions[0].ledgers[0].operation=deposit. got %s`, transaction0.Operation)
	}
	ledger := tRespBody.Data.Transactions[0].Ledgers[0]
	if ledger.EntryType != "credit" {
		log.Fatalf(`[T_0005_003] Transactions want transactions[0].ledgers[0].entry_type=credit. got %s`, ledger.EntryType)
	}
	if ledger.Amount != "40.123" {
		log.Fatalf(`[T_0005_003] Transactions want transactions[0].ledgers[0].amount="40.123". got %s`, ledger.Amount)
	}
}

func T_0006(client *client.Client) {
	username, wallets := SetupUserAndWalletCreation(client, "T_0006")

	wallet := wallets[0]
	dRespBody, dStatusCode, cErr := client.Withdraw(username, wallet.Id, decimal.NewFromFloat(50.1))
	if cErr != nil {
		log.Fatalf(`[T_0006_001] Withdraw transaction want nil err, got error %v`, cErr)
	}
	if dStatusCode != http.StatusBadRequest {
		log.Fatalf("[T_0006_001] Withdraw want 400. responseStatusCode=%d, err=%v", dStatusCode, cErr)
	}
	if dRespBody.Error == nil || *dRespBody.Error != "insufficient_funds" {
		var s string
		if dRespBody.Error != nil {
			s = *dRespBody.Error
		}
		log.Fatalf(`[T_0006_001] Withdraw want Response.err == "insufficient_funds". got Response.error=%v`, s)
	}
	if dRespBody.Data.Transaction.Id != 0 {
		log.Fatalf(`[T_0006_001] Withdraw want zero transaction.Id.`)
	}

	tRespBody, tStatusCode, err := client.Transactions(username)
	if tStatusCode != http.StatusOK {
		log.Fatalf(`[T_0006_003] Transactions want 200. Got responseStatusCode=%d, err=%#v`, tStatusCode, err)
	}

	if len(tRespBody.Data.Transactions) == 0 {
		log.Fatalf(`[T_0006_003] Transactions want transactions.len > 0. got 0`)
	}
	transaction0 := tRespBody.Data.Transactions[0]
	if len(transaction0.Ledgers) != 0 {
		log.Fatalf(`[T_0006_003] Transactions want transactions[0].ledgers[0].len == 0 (failed transaction should have 0 ledgers). got %d`, len(transaction0.Ledgers))
	}
	if transaction0.Operation != "withdraw" {
		log.Fatalf(`[T_0006_003] Transactions want transactions[0].ledgers[0].operation=withdraw. got %s`, transaction0.Operation)
	}
	if transaction0.Status != "error_insufficient_funds" {
		log.Fatalf(`[T_0006_003] Transactions want transactions[0].ledgers[0].status="error_insufficient_funds". got %s`, transaction0.Status)
	}
	if transaction0.MetaData.SourceWalletId == nil || *transaction0.MetaData.SourceWalletId != wallet.Id {
		log.Fatalf(`[T_0006_003] Transactions want transactions[0].metadata.source_wallet_id=%d. got %v`, wallet.Id, transaction0.MetaData.SourceWalletId)
	}
	if transaction0.MetaData.Amount == nil || *transaction0.MetaData.Amount != "50.1" {
		log.Fatalf(`[T_0006_003] Transactions want transactions[0].metadata.amount=%s. got %s`, "50.1", transaction0.MetaData.Amount)
	}
}
func SetupUserAndWalletCreation(client *client.Client, logPrefix string) (username string, wallets []client.Wallet) {
	username = NewRandomUserName(logPrefix, 6, 0)
	createUserResponseData, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[%s_001] SETUP CreateUser should succeed. responseStatusCode=%d, err=%v", logPrefix, responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[%s_001] SETUP CreateUser should succeed. err=%v", err)
	}
	if createUserResponseData.Error != nil && *createUserResponseData.Error != "" {
		log.Fatalf("[%s_001] SETUP CreateUser should succeed. createUserResponseData.Error=%v", logPrefix, err)
	}
	if createUserResponseData.Data.Id == 0 {
		log.Fatalf("[%s_001] SETUP CreateUser success but missing return id. createUserResponseData.Data.Id=%d", logPrefix, createUserResponseData.Data.Id)
	}
	if createUserResponseData.Data.Username != username {
		log.Fatalf("[%s_001] SETUP CreateUser success but missing return username. createUserResponseData.Data.Username=%s", logPrefix, createUserResponseData.Data.Username)
	}

	responseBody, responseStatusCode, err := client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[%s_002] SETUP Wallets want 200. Got responseStatusCode=%d, err=%#v`, logPrefix, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		log.Fatalf(`[%s_002] SETUP Wallets want 200. Got responseBody.Error=%s, err=%#v`, logPrefix, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) != 0 {
		log.Fatalf(`[%s_002] SETUP Wallets want no wallets. Got len=%d, err=%#v`, logPrefix, len(responseBody.Data.Wallets), err)
	}

	createWalletResponseBody, createWalletResponseStatusCode, err := client.CreateWallet(username, "SGD")
	if createWalletResponseStatusCode != http.StatusOK {
		log.Fatalf(`[%s_003] SETUP CreateWallet want 200. Got body.Err=%s, statusCode=%d, err=%#v`, logPrefix, *createWalletResponseBody.Error, createWalletResponseStatusCode, err)
	}
	if createWalletResponseBody.Error != nil {
		log.Fatalf(`[%s_003] SETUP CreateWallet want 200. Got createWalletResponseBody.Error=%s, err=%#v`, logPrefix, *createWalletResponseBody.Error, err)
	}
	wallet := createWalletResponseBody.Data.Wallet
	if wallet.Id == 0 {
		log.Fatalf(`[%s_003] SETUP CreateWallet should succeed. Got id=%v, err=%#v`, logPrefix, createWalletResponseBody.Data.Wallet, err)
	}

	responseBody, responseStatusCode, err = client.Wallets(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[%s_004] SETUP Wallets want 200. Got responseStatusCode=%d, err=%#v`, logPrefix, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		log.Fatalf(`[%s_004] SETUP Wallets want 200. Got responseBody.Error=%s, err=%#v`, logPrefix, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) == 0 {
		log.Fatalf(`[%s_004] SETUP Wallets want some wallets. Got len=%d, err=%#v`, logPrefix, len(responseBody.Data.Wallets), err)
	}

	if responseBody.Data.Wallets[0].Currency != "SGD" {
		log.Fatalf(`[%s_004] SETUP Wallets want currency as SGD. Got currency=%s, err=%#v`, logPrefix, responseBody.Data.Wallets[0].Currency, err)
	}

	return username, responseBody.Data.Wallets
}
