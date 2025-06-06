package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"syscall"

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
	_, _, err := client.GetWalletBalance(0)
	if errors.Is(err, syscall.ECONNREFUSED) {
		log.Fatalf("Server not started...")
	}
}

func T_0001(client *client.Client) {
	_, responseStatusCode, err := client.GetWalletBalance(0)

	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf(`GetWalletBalance should fail with status code 400. Got %d`, responseStatusCode)
	}

	if err == nil || !(err.Error() == "user id cannot be zero") {
		log.Fatalf(`GetWalletBalance should fail with error "user id cannot be zero". got err %#v, user id = %d`, err.Error(), 0)
	}

	responseData, responseStatusCode, err := client.GetWalletBalance(1)
	if responseStatusCode != http.StatusBadRequest {
		log.Fatalf(`GetWalletBalance should fail with status code 400. Got %d`, responseStatusCode)
	}
	if err == nil {
		log.Fatalf(`GetWalletBalance should fail with error "resource: user not found". got err %#v, user id = %d`, err.Error(), 0)
	}
	if responseData.Error != "resource: user not found" {
		log.Fatalf(`GetWalletBalance should fail with Response.error "resource: user not found". got Response.error %#v, user id = %d`, responseData.Error, 0)
	}
}
