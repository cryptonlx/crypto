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
}

func Ping(client *client.Client) {
	_, _, err := client.GetWalletBalance("")
	if errors.Is(err, syscall.ECONNREFUSED) {
		log.Fatalf("Server not started...")
	}
}

func T_0001(client *client.Client) {
	futureUserName := NewRandomUserName("t00001", 6, 2*24*time.Hour)
	responseData, responseStatusCode, err := client.GetWalletBalance(futureUserName)
	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf(`[T_0001] GetWalletBalance() should fail with status code 400. Got responseStatusCode=%d, err=%#v`, responseStatusCode, err)
	}
	if err == nil {
		log.Fatalf(`[T_0001] GetWalletBalance() should fail with error "resource: user not found". got err %#v, user id = %d`, err.Error(), 0)
	}
	if responseData.Error != nil && *responseData.Error != "resource: user not found" {
		log.Fatalf(`[T_0001] GetWalletBalance() should fail with Response.error "resource: user not found". got Response.error %s, user id = %d`, *responseData.Error, 0)
	}

	username := NewRandomUserName("t00001", 6, 0)
	createUserResponseData, responseStatusCode, err := client.CreateUser(username)
	if responseStatusCode != http.StatusOK {
		log.Fatalf("[T_0001] Create User should succeed. responseStatusCode=%v, err=%v", responseStatusCode, err)
	}
	if err != nil {
		log.Fatalf("[T_0001] Create User should succeed. err=%v", err)
	}
	if createUserResponseData.Error != nil && *createUserResponseData.Error != "" {
		log.Fatalf("[T_0001] Create User should succeed. createUserResponseData.Error=%v", err)
	}
}
