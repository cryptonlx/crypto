package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Params struct {
	DatabaseParams
	ServerParams
}

type DatabaseParams struct {
	ConnString string
}

type ServerParams struct {
	Port string
}

func LoadParams() (c Params, e error) {
	err := godotenv.Load()
	if !(os.Getenv("OPTIONAL_LOAD_ENV_FILE") == "TRUE") && err != nil {
		return Params{}, fmt.Errorf("error. loading .env file is compulsory. OPTIONAL_LOAD_ENV_FILE=%s.\n", os.Getenv("OPTIONAL_LOAD_ENV_FILE"))
	}
	c.ServerParams.Port = os.Getenv("LISTENING_PORT")

	// database
	dbUrl := os.Getenv("DATABASE_URL")
	c.DatabaseParams.ConnString = dbUrl

	// server
	return c, nil
}
