package e2e_tests

import (
	"os"
)

func main() {
	serverUrl := os.Getenv("SERVER_URL")
	client, err := NewClient(serverUrl)
	if err != nil {
		panic(err)
	}

	T_0001(client)
}

func T_0001(client *Client) {

}

type Client struct {
	serverUrl string
}

func NewClient(serverUrl string) (*Client, error) {
	return &Client{
		serverUrl: serverUrl,
	}, nil
}
