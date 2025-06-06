package user

import (
	"context"

	"github.com/cryptonlx/crypto/src/repositories/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repo {
	return &Repo{
		conn: conn,
	}
}

type User struct {
	Id       int64
	Username string
}

func (r *Repo) getUser(ctx context.Context, tx pgx.Tx, userId int64) (*User, error) {
	if tx == nil {
		return nil, utils.NilTxError
	}
	rows, err := tx.Query(ctx, "select id, username from user_accounts where id=$1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var t User
		rows.Scan(&t.Id, &t.Username)
		if err := rows.Err(); err != nil {
			return nil, err
		}
		users = append(users, t)
	}

	if len(users) == 0 {
		return nil, utils.NotFoundErrorF("user")
	}
	if len(users) > 1 {
		return nil, utils.RowLengthShouldBeAtMost1Error
	}
	return &users[0], nil
}

func (r *Repo) getWalletsByUserId(ctx context.Context, tx pgx.Tx, userId int64) ([]Wallet, error) {
	if tx == nil {
		return []Wallet{}, utils.NilTxError
	}
	rows, err := tx.Query(ctx, "select id, user_account_id, currency_type, value from wallets where user_account_id=$1", userId)
	if err != nil {
		return []Wallet{}, err
	}
	defer rows.Close()

	var wallets []Wallet
	for rows.Next() {
		var t Wallet
		rows.Scan(&t.Id, &t.CurrencyType, &t.Value)
		if err := rows.Err(); err != nil {
			return []Wallet{}, err
		}
		wallets = append(wallets, t)
	}
	return wallets, nil
}

type CurrencyType = string

const CurrencyTypeUSD CurrencyType = "USD"

type Wallet struct {
	Id           int64
	CurrencyType string
	Value        int64
}

type UserBalance struct {
	User    User
	Wallets []Wallet
}

func (r *Repo) GetUserBalance(ctx context.Context, userId int64) (UserBalance, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return UserBalance{}, err
	}
	defer tx.Rollback(ctx)

	user, err := r.getUser(ctx, tx, userId)
	if err != nil {
		return UserBalance{}, err
	}

	wallets, err := r.getWalletsByUserId(ctx, tx, user.Id)
	if err != nil {
		return UserBalance{}, err
	}
	return UserBalance{
		User:    *user,
		Wallets: wallets,
	}, nil
}
