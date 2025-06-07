package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/cryptonlx/crypto/cmd/e2e_tests/client"
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
}

func Ping(client *client.Client) {
	_, _, err := client.GetWalletBalances("")
	if errors.Is(err, syscall.ECONNREFUSED) {
		log.Fatalf("Server not started...")
	}
}

func T_0001(client *client.Client) {
	futureUserName := NewRandomUserName("t00001", 6, 2*24*time.Hour)
	responseBody, responseStatusCode, err := client.GetWalletBalances(futureUserName)
	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf(`[T_0001_001] GetWalletBalances() want status code 400. got code=%d, err=%#v`, responseStatusCode, err)
	}
	if err == nil || responseBody.Error == nil || *responseBody.Error != "resource: user not found" {
		log.Fatalf(`[T_0001_001] GetWalletBalances() want Response.error "resource: user not found". got err=%v, Response.error%v, user id=%d`, err, responseBody.Error, 0)
	}

	username := NewRandomUserName("t00001", 6, 0)
	createUserResponseBody, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[T_0001_002] Create User want 200. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[T_0001_002] Create User want nil err. err=%v", err)
	}
	if createUserResponseBody.Error != nil {
		log.Fatalf("[T_0001_002] Create User want nil err. createUserResponseBody.Error=%v", *createUserResponseBody.Error)
	}

	if createUserResponseBody.Data.Id == 0 {
		log.Fatalf("[T_0001_002] Create User success. want id. got id=%d", createUserResponseBody.Data.Id)
	}

	if createUserResponseBody.Data.Username != username {
		log.Fatalf("[T_0001_002] Create User success but not return username. createUserResponseBody.Data.Username=%s", createUserResponseBody.Data.Username)
	}

	createUserResponseBody, responseStatusCode, err = client.CreateUser(username)
	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf("[T_0001_003] Create Duplicate User want 400. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}

	if createUserResponseBody.Error == nil || *createUserResponseBody.Error != "unique_violation" {
		log.Fatalf(`[T_0001_003] Create Duplicate User want "unique_violation". responseStatusCode=%d, err=%v`, *createUserResponseBody.Error, err)
	}

	responseBody, responseStatusCode, err = client.GetWalletBalances(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0001_004] GetWalletBalances() should succeed with 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
}

func T_0002(client *client.Client) {
	futureUserName := NewRandomUserName("t00002", 6, 2*24*time.Hour)
	responseBody, responseStatusCode, err := client.GetTransactionHistory(futureUserName)
	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf(`[T_0002_001] GetTransactionHistory() want status code 400. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if err == nil {
		log.Fatalf(`[T_0002_001] GetTransactionHistory() want error "resource: user not found". got err %#v, user id = %d`, err.Error(), 0)
	}
	if responseBody.Error != nil && *responseBody.Error != "resource: user not found" {
		log.Fatalf(`[T_0002_001] GetTransactionHistory() want Response.error "resource: user not found". got Response.error %s, user id = %d`, *responseBody.Error, 0)
	}

	username := NewRandomUserName("t00002", 6, 0)
	createUserResponseBody, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[T_0002_002] Create User should succeed. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[T_0002_002] Create User should succeed. err=%v", err)
	}
	if createUserResponseBody.Error != nil && *createUserResponseBody.Error != "" {
		log.Fatalf("[T_0002_002] Create User should succeed. createUserResponseBody.Error=%v", err)
	}

	if createUserResponseBody.Data.Id == 0 {
		log.Fatalf("[T_0002_002] Create User success but not return id. createUserResponseBody.Data.Id=%d", createUserResponseBody.Data.Id)
	}

	if createUserResponseBody.Data.Username != username {
		log.Fatalf("[T_0002_002] Create User success but not return username. createUserResponseBody.Data.Username=%s", createUserResponseBody.Data.Username)
	}

	responseBody, responseStatusCode, err = client.GetTransactionHistory(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0002_003] GetTransactionHistory() should succeed with 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
}

func T_0003(client *client.Client) {
	SetupUserAndWalletCreation(client, "T_0003")
}

func T_0004(client *client.Client) {
	username, wallets := SetupUserAndWalletCreation(client, "T_0004")

	wallet := wallets[0]
	balanceBefore := wallet.Balance

	dRespBody, dStatusCode, err := client.Deposit(wallet.Id, -40)
	if dStatusCode != http.StatusBadRequest {
		log.Fatalf("[T_0004_001] Deposit want err 400. responseStatusCode=%d, err=%v", dStatusCode, err)
	}
	if err != nil {
		log.Fatalf(`[T_0004_001] Deposit transaction want nil err, got error %v`, err)
	}
	if err != nil || *dRespBody.Error != "invalid_amount" {
		log.Fatalf(`[T_0004_001] Deposit() want Response.error "invalid_amount". got err=%v, Response.error=%v`, err, *dRespBody.Error)
	}

	responseBody, responseStatusCode, err := client.GetWalletBalances(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0004_002] SETUP GetWalletBalances() should succeed with 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		log.Fatalf(`[T_0004_002] SETUP GetWalletBalances() should succeed with 200. Got responseBody.Error=%s, err=%#v`, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) == 0 {
		log.Fatalf(`[T_0004_002] SETUP GetWalletBalances() should succeed with some wallets. Got len=%d, err=%#v`, len(responseBody.Data.Wallets), err)
	}
	if responseBody.Data.Wallets[0].Balance != balanceBefore {
		log.Fatalf(`[T_0004_002] SETUP GetWalletBalances() balance before and after should be same. want=%s, got=%s`, balanceBefore, responseBody.Data.Wallets[0].Balance)
	}
}

func SetupUserAndWalletCreation(client *client.Client, logPrefix string) (username string, wallets []client.Wallet) {
	username = NewRandomUserName(logPrefix, 6, 0)
	createUserResponseData, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[%s_001] SETUP Create User should succeed. responseStatusCode=%d, err=%v", logPrefix, responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[%s_001] SETUP Create User should succeed. err=%v", err)
	}
	if createUserResponseData.Error != nil && *createUserResponseData.Error != "" {
		log.Fatalf("[%s_001] SETUP Create User should succeed. createUserResponseData.Error=%v", logPrefix, err)
	}
	if createUserResponseData.Data.Id == 0 {
		log.Fatalf("[%s_001] SETUP Create User success but missing return id. createUserResponseData.Data.Id=%d", logPrefix, createUserResponseData.Data.Id)
	}
	if createUserResponseData.Data.Username != username {
		log.Fatalf("[%s_001] SETUP Create User success but missing return username. createUserResponseData.Data.Username=%s", logPrefix, createUserResponseData.Data.Username)
	}

	responseBody, responseStatusCode, err := client.GetWalletBalances(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[%s_002] SETUP GetWalletBalances() should succeed with 200. Got responseStatusCode=%d, err=%#v`, logPrefix, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		log.Fatalf(`[%s_002] SETUP GetWalletBalances() should succeed with 200. Got responseBody.Error=%s, err=%#v`, logPrefix, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) != 0 {
		log.Fatalf(`[%s_002] SETUP GetWalletBalances() should succeed with no wallets. Got len=%d, err=%#v`, logPrefix, len(responseBody.Data.Wallets), err)
	}

	createWalletResponseBody, createWalletResponseStatusCode, err := client.CreateWallet(username, "SGD")
	if createWalletResponseStatusCode != http.StatusOK {
		log.Fatalf(`[%s_003] SETUP CreateWallet() should succeed with 200. Got body.Err=%s, statusCode=%d, err=%#v`, logPrefix, *createWalletResponseBody.Error, createWalletResponseStatusCode, err)
	}
	if createWalletResponseBody.Error != nil {
		log.Fatalf(`[%s_003] SETUP CreateWallet() should succeed with 200. Got createWalletResponseBody.Error=%s, err=%#v`, logPrefix, *createWalletResponseBody.Error, err)
	}
	wallet := createWalletResponseBody.Data.Wallet
	if wallet.Id == 0 {
		log.Fatalf(`[%s_003] SETUP CreateWallet() should succeed. Got id=%v, err=%#v`, logPrefix, createWalletResponseBody.Data.Wallet, err)
	}

	responseBody, responseStatusCode, err = client.GetWalletBalances(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[%s_004] SETUP GetWalletBalances() should succeed with 200. Got responseStatusCode=%d, err=%#v`, logPrefix, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		log.Fatalf(`[%s_004] SETUP GetWalletBalances() should succeed with 200. Got responseBody.Error=%s, err=%#v`, logPrefix, *responseBody.Error, err)
	}
	if len(responseBody.Data.Wallets) == 0 {
		log.Fatalf(`[%s_004] SETUP GetWalletBalances() should succeed with some wallets. Got len=%d, err=%#v`, logPrefix, len(responseBody.Data.Wallets), err)
	}

	if responseBody.Data.Wallets[0].Currency != "SGD" {
		log.Fatalf(`[%s_004] SETUP GetWalletBalances() should succeed with currency as SGD. Got currency=%s, err=%#v`, logPrefix, responseBody.Data.Wallets[0].Currency, err)
	}

	return username, responseBody.Data.Wallets
}
