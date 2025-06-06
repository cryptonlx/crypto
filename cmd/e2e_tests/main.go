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
		log.Fatalf(`[T_0001_001] GetWalletBalances() should fail with status code 400. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if err == nil {
		log.Fatalf(`[T_0001_001] GetWalletBalances() should fail with error "resource: user not found". got err %#v, user id = %d`, err.Error(), 0)
	}
	if responseBody.Error != nil && *responseBody.Error != "resource: user not found" {
		log.Fatalf(`[T_0001_001] GetWalletBalances() should fail with Response.error "resource: user not found". got Response.error %s, user id = %d`, *responseBody.Error, 0)
	}

	username := NewRandomUserName("t00001", 6, 0)
	createUserResponseBody, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[T_0001_002] Create User should succeed. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[T_0001_002] Create User should succeed. err=%v", err)
	}
	if createUserResponseBody.Error != nil && *createUserResponseBody.Error != "" {
		log.Fatalf("[T_0001_002] Create User should succeed. createUserResponseBody.Error=%v", err)
	}

	if createUserResponseBody.Data.Id == 0 {
		log.Fatalf("[T_0001_002] Create User success but not return id. createUserResponseBody.Data.Id=%d", createUserResponseBody.Data.Id)
	}

	if createUserResponseBody.Data.Username != username {
		log.Fatalf("[T_0001_002] Create User success but not return username. createUserResponseBody.Data.Username=%s", createUserResponseBody.Data.Username)
	}

	createUserResponseBody, responseStatusCode, err = client.CreateUser(username)
	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf("[T_0001_003] Create Duplicate User should fail with 400. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}

	if createUserResponseBody.Error == nil || *createUserResponseBody.Error != "unique_violation" {
		log.Fatalf(`[T_0001_003] Create Duplicate User should fail with "unique_violation". responseStatusCode=%d, err=%v`, *createUserResponseBody.Error, err)
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
		log.Fatalf(`[T_0002_001] GetTransactionHistory() should fail with status code 400. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if err == nil {
		log.Fatalf(`[T_0002_001] GetTransactionHistory() should fail with error "resource: user not found". got err %#v, user id = %d`, err.Error(), 0)
	}
	if responseBody.Error != nil && *responseBody.Error != "resource: user not found" {
		log.Fatalf(`[T_0002_001] GetTransactionHistory() should fail with Response.error "resource: user not found". got Response.error %s, user id = %d`, *responseBody.Error, 0)
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
	username := NewRandomUserName("t00003", 6, 0)
	createUserResponseData, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[T_0003_001] Create User should succeed. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[T_0003_001] Create User should succeed. err=%v", err)
	}
	if createUserResponseData.Error != nil && *createUserResponseData.Error != "" {
		log.Fatalf("[T_0003_001] Create User should succeed. createUserResponseData.Error=%v", err)
	}
	if createUserResponseData.Data.Id == 0 {
		log.Fatalf("[T_0003_001] Create User success but missing return id. createUserResponseData.Data.Id=%d", createUserResponseData.Data.Id)
	}
	if createUserResponseData.Data.Username != username {
		log.Fatalf("[T_0003_001] Create User success but missing return username. createUserResponseData.Data.Username=%s", createUserResponseData.Data.Username)
	}

	responseBody, responseStatusCode, err := client.GetWalletBalances(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0003_002] GetWalletBalances() should succeed with 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		log.Fatalf(`[T_0003_002] GetWalletBalances() should succeed with 200. Got responseBody.Error=%s, err=%#v`, *responseBody.Error, err)
	}
	if len(responseBody.Data.WalletBalances) != 0 {
		log.Fatalf(`[T_0003_002] GetWalletBalances() should succeed with no wallets. Got len=%d, err=%#v`, len(responseBody.Data.WalletBalances), err)
	}

	createWalletResponseBody, createWalletResponseStatusCode, err := client.CreateWallet(username, "SGD")
	if createWalletResponseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0003_003] CreateWallet() should succeed with 200. Got createWalletResponseStatusCode=%d, err=%#v`, createWalletResponseStatusCode, err)
	}
	if createWalletResponseBody.Error != nil {
		log.Fatalf(`[T_0003_003] CreateWallet() should succeed with 200. Got createWalletResponseBody.Error=%s, err=%#v`, *createWalletResponseBody.Error, err)
	}

	wallet := createWalletResponseBody.Data.Wallet
	if wallet.Id == 0 {
		log.Fatalf(`[T_0003_003] CreateWallet() should succeed. Got id=%v, err=%#v`, createWalletResponseBody.Data.Wallet, err)
	}

	responseBody, responseStatusCode, err = client.GetWalletBalances(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0003_004] GetWalletBalances() should succeed with 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseBody.Error != nil {
		log.Fatalf(`[T_0003_004] GetWalletBalances() should succeed with 200. Got responseBody.Error=%s, err=%#v`, *responseBody.Error, err)
	}
	if len(responseBody.Data.WalletBalances) == 0 {
		log.Fatalf(`[T_0003_004] GetWalletBalances() should succeed with some wallets. Got len=%d, err=%#v`, len(responseBody.Data.WalletBalances), err)
	}

	if responseBody.Data.WalletBalances[0].Currency != "SGD" {
		log.Fatalf(`[T_0003_004] GetWalletBalances() should succeed with currency as SGD. Got currency=%s, err=%#v`, responseBody.Data.WalletBalances[0].Currency, err)
	}
}

func T_0004_WIP(client *client.Client) {
	username := NewRandomUserName("t00003", 6, 0)
	createUserResponseData, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[T_0003_002] Create User should succeed. responseStatusCode=%d, err=%v", responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[T_0003_002] Create User should succeed. err=%v", err)
	}
	if createUserResponseData.Error != nil && *createUserResponseData.Error != "" {
		log.Fatalf("[T_0003_002] Create User should succeed. createUserResponseData.Error=%v", err)
	}
	if createUserResponseData.Data.Id == 0 {
		log.Fatalf("[T_0003_002] Create User success but missing return id. createUserResponseData.Data.Id=%d", createUserResponseData.Data.Id)
	}
	if createUserResponseData.Data.Username != username {
		log.Fatalf("[T_0002_002] Create User success but missing return username. createUserResponseData.Data.Username=%s", createUserResponseData.Data.Username)
	}

	responseData, responseStatusCode, err := client.GetTransactionHistory(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf(`[T_0002_003] GetTransactionHistory() should succeed with 200. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if responseData.Error != nil {
		log.Fatalf(`[T_0002_003] GetTransactionHistory() should succeed with 200. Got responseData.Error=%s, err=%#v`, *responseData.Error, err)
	}
	if len(responseData.Data.Transactions) != 0 {
		log.Fatalf(`[T_0002_003] GetTransactionHistory() should succeed with empty transactions. Got len=%s, err=%#v`, len(responseData.Data.Transactions), err)
	}

	depositAmount := -50
	client.Deposit(username, int64(depositAmount))

}
