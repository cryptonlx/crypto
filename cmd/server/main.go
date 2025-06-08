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
	"github.com/cryptonlx/crypto/src/httplog"

	usermux "github.com/cryptonlx/crypto/src/controllers/mux/user"
	userrepo "github.com/cryptonlx/crypto/src/repositories/user"
	userservice "github.com/cryptonlx/crypto/src/service/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"

	_ "github.com/cryptonlx/crypto/docs"
	"github.com/swaggo/http-swagger"
)

// @securityDefinitions.basic BasicAuth
func main() {
	configParams, dbConnPool, err := Init()
	if err != nil {
		log.Fatal(err)
	}

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT /*keyboard input*/, syscall.SIGTERM /*process kill*/)
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	userRepo := userrepo.New(dbConnPool)
	userService := userservice.New(userRepo)
	userHandlers := usermux.NewHandlers(userService)

	mux.HandleFunc("GET /user/{username}/wallets", userHandlers.Wallets)
	mux.HandleFunc("GET /user/{username}/transactions", userHandlers.Transactions)
	mux.HandleFunc("POST /user", userHandlers.CreateUser)
	mux.HandleFunc("POST /wallet", userHandlers.CreateWallet)
	mux.HandleFunc("POST /wallet/{wallet_id}/deposit", userHandlers.Deposit)
	mux.HandleFunc("POST /wallet/{wallet_id}/withdraw", userHandlers.Withdraw)
	mux.HandleFunc("POST /wallet/{wallet_id}/transfer", userHandlers.Transfer)

	go func() {
		log.Println("Listening on " + configParams.ServerParams.Port)
		limiter := rate.NewLimiter(1000, 1000)

		server := &http.Server{
			Addr: configParams.ServerParams.Port,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !limiter.Allow() {
					http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
					return
				}
				r = httplog.ContextualizeHttpRequest(r)
				log.Printf("%s [request received]\n", httplog.SPrintHttpRequestPrefix(r))
				mux.ServeHTTP(w, r)
			}),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			log.Printf("ListenAndServe err %v\n", err)
		}
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
