package user

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"log"
	"time"

	"github.com/cryptonlx/crypto/src/repositories/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *Repo) getUser(ctx context.Context, tx pgx.Tx, username string) (*User, error) {
	if tx == nil {
		return nil, utils.NilTxError
	}
	rows, err := tx.Query(ctx, "select id, username from user_accounts where username=$1", username)
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

type UserWallet struct {
	User   User
	Wallet Wallet
}

func (r *Repo) getUserAndWalletByWalletId(ctx context.Context, tx pgx.Tx, walletId int64) (*UserWallet, error) {
	if tx == nil {
		return nil, utils.NilTxError
	}
	rows, err := tx.Query(ctx, "select ua.id, ua.username, w.id, w.user_account_id, w.currency, w.amount from user_accounts ua join wallets w on w.user_account_id = ua.id  where w.id=$1 FOR UPDATE OF w", walletId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserWallet
	for rows.Next() {
		var t UserWallet
		rows.Scan(&t.User.Id, &t.User.Username, &t.Wallet.Id, &t.Wallet.UserAccountId, &t.Wallet.Currency, &t.Wallet.Balance)
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

func (r *Repo) createUser(ctx context.Context, tx pgx.Tx, username string) (User, error) {
	if tx == nil {
		return User{}, utils.NilTxError
	}
	row := tx.QueryRow(ctx, "insert into user_accounts(username) VALUES ($1) RETURNING id, username", username)

	var user User
	err := row.Scan(&user.Id, &user.Username)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			err = utils.ToError(pgErr)
		}
		return User{}, err
	}
	return user, nil
}

func (r *Repo) createWallet(ctx context.Context, tx pgx.Tx, username int64, currency CurrencyType) (Wallet, error) {
	row := tx.QueryRow(ctx, "insert into wallets(user_account_id, currency, balance) VALUES ($1,$2,$3) RETURNING id, user_account_id, currency, balance", username, currency, decimal.Zero)

	var wallet Wallet
	err := row.Scan(&wallet.Id, &wallet.UserAccountId, &wallet.Currency, &wallet.Balance)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			err = utils.ToError(pgErr)
		}
		return Wallet{}, err
	}

	return wallet, nil
}

func (r *Repo) getWalletsByUserId(ctx context.Context, tx pgx.Tx, userId int64) ([]Wallet, error) {
	if tx == nil {
		return []Wallet{}, utils.NilTxError
	}

	rows, err := tx.Query(ctx, "select id, user_account_id, currency, balance from wallets where user_account_id=$1", userId)
	if err != nil {
		return []Wallet{}, err
	}
	defer rows.Close()

	var wallets []Wallet
	for rows.Next() {
		var t Wallet
		rows.Scan(&t.Id, &t.UserAccountId, &t.Currency, &t.Balance)
		if err := rows.Err(); err != nil {
			return []Wallet{}, err
		}
		wallets = append(wallets, t)
	}

	log.Printf("%d l%d", userId, len(wallets))

	return wallets, nil
}

type Transaction struct{}

func (r *Repo) getTransactionsByUserId(ctx context.Context, tx pgx.Tx, userId int64) ([]Transaction, error) {
	if tx == nil {
		return []Transaction{}, utils.NilTxError
	}
	// TODO
	rows, err := tx.Query(ctx, "select id, user_account_id, currency, balance from wallets where user_account_id=$1", userId)
	if err != nil {
		return []Transaction{}, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		//rows.Scan(&t.Id, &t.CurrencyType, &t.Balance)
		//if err := rows.Err(); err != nil {
		//	return []Transaction{}, err
		//}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

type Wallet struct {
	Id            int64
	UserAccountId int64
	Currency      string
	Balance       decimal.Decimal
}

type UserWallets struct {
	User    User
	Wallets []Wallet
}

func (r *Repo) UserWallets(ctx context.Context, username string) (UserWallets, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return UserWallets{}, err
	}
	defer tx.Rollback(ctx)

	user, err := r.getUser(ctx, tx, username)
	if err != nil {
		return UserWallets{}, err
	}

	wallets, err := r.getWalletsByUserId(ctx, tx, user.Id)
	if err != nil {
		return UserWallets{}, err
	}
	return UserWallets{
		User:    *user,
		Wallets: wallets,
	}, nil
}

type UserTransactions struct {
	User         User
	Transactions []Transaction
}

func (r *Repo) Transactions(ctx context.Context, username string) (UserTransactions, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return UserTransactions{}, err
	}
	defer tx.Rollback(ctx)

	user, err := r.getUser(ctx, tx, username)
	if err != nil {
		return UserTransactions{}, err
	}

	transactions, err := r.getTransactionsByUserId(ctx, tx, user.Id)
	if err != nil {
		return UserTransactions{}, err
	}
	return UserTransactions{
		User:         *user,
		Transactions: transactions,
	}, nil
}

func (r *Repo) CreateUser(ctx context.Context, username string) (User, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback(ctx)

	user, err := r.createUser(ctx, tx, username)
	if err != nil {
		return User{}, err
	}

	tx.Commit(ctx)
	return user, nil
}

// CurrencyType
// Subset of ISO4217
type CurrencyType string

const CurrencyTypeUSD CurrencyType = "USD"
const CurrencyTypeSGD CurrencyType = "SGD"

func (r *Repo) CreateWallet(ctx context.Context, username string, currency CurrencyType) (Wallet, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return Wallet{}, err
	}
	defer tx.Rollback(ctx)
	user, err := r.getUser(ctx, tx, username)
	if err != nil {
		return Wallet{}, err
	}

	wallet, err := r.createWallet(ctx, tx, user.Id, currency)
	if err != nil {
		return Wallet{}, err
	}

	tx.Commit(ctx)
	return wallet, nil
}

func (r *Repo) Deposit(ctx context.Context, nonce int64, walletId int64, amount decimal.Decimal) (Ledger, error) {
	if !amount.IsPositive() {
		return Ledger{}, errors.New("amount negative")
	}

	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		return Ledger{}, err
	}
	defer tx.Rollback(context.Background())

	userWallet, err := r.getUserAndWalletByWalletId(ctx, tx, walletId)
	if err != nil {
		return Ledger{}, err
	}

	balance := userWallet.Wallet.Balance.Add(amount)
	err = r.updateBalance(ctx, tx, walletId, balance)
	if err != nil {

		return Ledger{}, err
	}
	ledger, err := r.addLedger(ctx, tx, nonce, walletId, "deposit", "credit", amount, balance)
	if err != nil {
		return Ledger{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return Ledger{}, err
	}
	return ledger, nil
}

type Ledger struct {
	Id        int64
	WalletId  int64
	Nonce     int64
	Operation string
	EntryType string
	Amount    decimal.Decimal
	CreatedAt time.Time
	Balance   decimal.Decimal
}

func (r *Repo) addLedger(ctx context.Context, tx pgx.Tx, nonce, walletId int64, operation, entryType string, amount, balance decimal.Decimal) (Ledger, error) {
	if tx == nil {
		return Ledger{}, utils.NilTxError
	}

	//entry_type
	//amount
	//created_at
	row := tx.QueryRow(ctx, "insert into ledgers(wallet_id, nonce, operation, entry_type, amount, created_at, balance) values ($1,$2,$3,$4,$5,now(), $6) returning id, wallet_id, nonce, operation, entry_type, amount, created_at, balance", walletId, nonce, operation, entryType, amount, balance)

	var l Ledger
	err := row.Scan(&l.Id, &l.WalletId, &l.Nonce, &l.Operation, &l.EntryType, &l.Amount, &l.CreatedAt, &l.Balance)
	if err != nil {
		return Ledger{}, err
	}
	return l, nil
}

func (r *Repo) updateBalance(ctx context.Context, tx pgx.Tx, id int64, balance decimal.Decimal) error {
	if tx == nil {
		return utils.NilTxError
	}
	_, err := tx.Exec(ctx, `
    UPDATE wallets
    SET balance = $1 where id = $2`, balance, id)
	return err
}
