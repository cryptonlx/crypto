package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	serverconfig "github.com/cryptonlx/crypto/cmd/server/config"

	usermux "github.com/cryptonlx/crypto/src/controller/mux/user"
	userrepo "github.com/cryptonlx/crypto/src/repositories/user"
	userservice "github.com/cryptonlx/crypto/src/service/user"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/cryptonlx/crypto/docs"
	"github.com/swaggo/http-swagger"
)

func main() {
	configParams, dbConnPool, err := Init()
	if err != nil {
		log.Fatal(err)
	}

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT /*keyboard input*/, syscall.SIGTERM /*process kill*/)
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	taskRepo := userrepo.New(dbConnPool)
	taskService := userservice.New(taskRepo)
	taskHandlers := usermux.NewHandlers(taskService)

	mux.HandleFunc("GET /user/{username}/balance", taskHandlers.GetWalletBalances)
	mux.HandleFunc("GET /user/{username}/transactions", taskHandlers.GetWalletTransactions)
	mux.HandleFunc("POST /user", taskHandlers.CreateUser)

	go func() {
		log.Println("Listening on " + configParams.ServerParams.Port)
		http.ListenAndServe(configParams.ServerParams.Port, mux)
	}()

	recvSig := <-interruptSignal
	log.Println("Received signal: " + recvSig.String() + " ; tearing down...")
	log.Println("Terminating hepmilserver::main()...")
}

func Init() (serverconfig.Params, *pgxpool.Pool, error) {
	params, err := serverconfig.LoadParams()
	if err != nil {
		return serverconfig.Params{}, nil, err
	}

	// Db Connection Pool
	if params.ConnString == "" {
		return serverconfig.Params{}, nil, errors.New("no db connection string provided")
	}
	pgConfig, err := pgxpool.ParseConfig(params.ConnString)
	if err != nil {
		return serverconfig.Params{}, nil, err
	}

	pgConfig.MaxConns = 10
	pgConfig.MinConns = 2
	pgConfig.MaxConnIdleTime = 5 * time.Minute
	ctx := context.Background()

	// GetWalletBalances the pool
	dbConnPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return serverconfig.Params{}, nil, err
	}

	err = dbConnPool.Ping(ctx)
	if err != nil {
		panic(err)
	}
	return params, dbConnPool, nil
}
