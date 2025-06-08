package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cryptonlx/crypto/src/repositories/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
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

func (r *Repo) user(ctx context.Context, tx pgx.Tx, username string) (*User, error) {
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

func (r *Repo) userWalletByWalletId(ctx context.Context, tx pgx.Tx, walletId int64) (*UserWallet, error) {
	if tx == nil {
		return nil, utils.NilTxError
	}
	rows, err := tx.Query(ctx, "select ua.id, ua.username, w.id, w.user_account_id, w.currency, w.balance from user_accounts ua join wallets w on w.user_account_id = ua.id  where w.id=$1 FOR UPDATE OF w", walletId)
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

	return wallets, nil
}

type TransactionLedgers struct {
	Transaction Transaction
	//Ledgers     []Ledger
	Ledgers []Ledger
}

func (r *Repo) getTransactionsByUserId(ctx context.Context, tx pgx.Tx, userId int64) ([]TransactionLedgers, error) {
	if tx == nil {
		return []TransactionLedgers{}, utils.NilTxError
	}

	rows, err := tx.Query(ctx, `select t.id,t.requestor_id, t.nonce, t.status, t.operation,t.created_at, t.metadata, COALESCE(json_agg(json_build_object('id',l.wallet_id,'transaction_id',l.transaction_id,'entry_type', l.entry_type,'amount', l.amount,'created_at', l.created_at,'balance', l.balance)) filter (where l.id is not null), '[]'::json)
										from transactions t
												 left join user_accounts ua on ua.id = t.requestor_id
												 left join ledgers l on t.id = l.transaction_id
											where t.requestor_id = $1
										group by t.id;
`, userId)
	if err != nil {
		return []TransactionLedgers{}, err
	}
	defer rows.Close()

	var transactionLedgers []TransactionLedgers
	for rows.Next() {
		var t Transaction
		var ledgersRaw json.RawMessage
		rows.Scan(&t.Id,
			&t.RequestorId,
			&t.Nonce,
			&t.Status,
			&t.Operation,
			&t.CreatedAt,
			&t.MetaData,
			&ledgersRaw)
		var ledgers []Ledger
		log.Printf("METADATA %v", t.MetaData)
		err := json.Unmarshal(ledgersRaw, &ledgers)
		if err != nil {
			return []TransactionLedgers{}, err
		}
		if err := rows.Err(); err != nil {
			return []TransactionLedgers{}, err
		}
		transactionLedgers = append(transactionLedgers, TransactionLedgers{
			Transaction: t,
			Ledgers:     ledgers,
		})
	}
	log.Printf("transactionLedgers %+v\n", transactionLedgers)
	return transactionLedgers, nil
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

	user, err := r.user(ctx, tx, username)
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

func (r *Repo) User(ctx context.Context, username string) (*User, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	user, err := r.user(ctx, tx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserTransactions struct {
	User         User
	Transactions []Transaction
}

func (r *Repo) Transactions(ctx context.Context, username string) ([]TransactionLedgers, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return []TransactionLedgers{}, err
	}
	defer tx.Rollback(ctx)

	user, err := r.user(ctx, tx, username)
	if err != nil {
		return []TransactionLedgers{}, err
	}

	transactions, err := r.getTransactionsByUserId(ctx, tx, user.Id)
	if err != nil {
		return []TransactionLedgers{}, err
	}
	return transactions, nil
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
	user, err := r.user(ctx, tx, username)
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

func (r *Repo) Deposit(requestor string, ctx context.Context, nonce int64, walletId int64, amount decimal.Decimal) (Transaction, Ledger, error) {
	if !amount.IsPositive() {
		return Transaction{}, Ledger{}, errors.New("amount negative")
	}

	user, err := r.User(ctx, requestor)
	transaction, err := r.newTransaction(ctx, nonce, user.Id, "deposit", nil)
	if err != nil {
		log.Printf("err%v\n", err)
		return Transaction{}, Ledger{}, err
	}

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return Transaction{}, Ledger{}, err
	}
	defer tx.Rollback(ctx)

	userWallet, err := r.userWalletByWalletId(ctx, tx, walletId)
	if err != nil {
		return Transaction{}, Ledger{}, err
	}
	if requestor != userWallet.User.Username {
		return Transaction{}, Ledger{}, errors.New("requestor and wallet owner mismatch")
	}

	newBalance := userWallet.Wallet.Balance.Add(amount)
	err = r.updateBalance(ctx, tx, walletId, newBalance)
	if err != nil {
		return Transaction{}, Ledger{}, err
	}

	ledger, err := r.appendLedger(ctx, tx, nonce, walletId, transaction.Id, "credit", amount, newBalance)
	if err != nil {
		return Transaction{}, Ledger{}, err
	}

	err = r.updateTransactionStatus(ctx, tx, transaction.Id, "success")
	if err != nil {
		return Transaction{}, Ledger{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Transaction{}, Ledger{}, err
	}
	return transaction, ledger, nil
}

func (r *Repo) Withdraw(requestor string, ctx context.Context, nonce int64, walletId int64, amount decimal.Decimal) (Transaction, Ledger, error) {
	if !amount.IsPositive() {
		return Transaction{}, Ledger{}, errors.New("amount negative")
	}

	user, err := r.User(ctx, requestor)
	transaction, err := r.newTransaction(ctx, nonce, user.Id, "withdraw", map[string]any{
		"amount":           amount.String(),
		"source_wallet_id": walletId,
	})
	if err != nil {
		return Transaction{}, Ledger{}, err
	}

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return Transaction{}, Ledger{}, err
	}
	defer tx.Rollback(ctx)

	userWallet, err := r.userWalletByWalletId(ctx, tx, walletId)
	if err != nil {
		return Transaction{}, Ledger{}, err
	}
	if requestor != userWallet.User.Username {
		return Transaction{}, Ledger{}, errors.New("requestor and wallet owner mismatch")
	}

	newBalance := userWallet.Wallet.Balance.Sub(amount)
	err = r.updateBalance(ctx, tx, walletId, newBalance)
	if err != nil {
		tsErr := r.UpdateTransactionStatus(context.Background(), transaction.Id, fmt.Sprintf("error_%s", err.Error()))
		if tsErr != nil {
			return Transaction{}, Ledger{}, errors.Join(tsErr, err)
		}
		return Transaction{}, Ledger{}, err
	}

	ledger, err := r.appendLedger(ctx, tx, nonce, walletId, transaction.Id, "debit", amount, newBalance)
	if err != nil {
		return Transaction{}, Ledger{}, err
	}

	err = r.updateTransactionStatus(ctx, tx, transaction.Id, "success")
	if err != nil {
		return Transaction{}, Ledger{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Transaction{}, Ledger{}, err
	}
	return transaction, ledger, nil
}

type Ledger struct {
	Id            int64           `json:"id"`
	WalletId      int64           `json:"wallet_id"`
	EntryType     string          `json:"entry_type"`
	Amount        decimal.Decimal `json:"amount"`
	CreatedAt     time.Time       `json:"created_at"`
	Balance       decimal.Decimal `json:"balance"`
	TransactionId int64           `json:"transaction_id"`
}

func (r *Repo) appendLedger(ctx context.Context, tx pgx.Tx, nonce int64, walletId int64, transactionId int64, entryType string, amount, balance decimal.Decimal) (Ledger, error) {
	if tx == nil {
		return Ledger{}, utils.NilTxError
	}

	row := tx.QueryRow(ctx, `insert into ledgers(wallet_id, entry_type, amount, balance, transaction_id, created_at)
		values ($1,$2,$3,$4,$5,now()) returning id, wallet_id, entry_type, amount, created_at, balance, transaction_id`,
		walletId, entryType, amount, balance, transactionId)

	var l Ledger
	err := row.Scan(&l.Id, &l.WalletId, &l.EntryType, &l.Amount, &l.CreatedAt, &l.Balance, &l.TransactionId)
	if err != nil {
		return Ledger{}, err
	}
	return l, nil
}

type TransactionMetaData struct {
	SourceWalletId *int64           `json:"source_wallet_id" example:"1"`
	Amount         *decimal.Decimal `json:"amount" example:"1"`
}

type Transaction struct {
	Id          int64
	RequestorId int64
	Nonce       int64
	Status      string
	Operation   string
	CreatedAt   time.Time
	MetaData    TransactionMetaData
}

func (r *Repo) newTransaction(ctx context.Context, nonce, requestorId int64, operation string, metaData map[string]any) (Transaction, error) {
	if metaData == nil {
		metaData = make(map[string]any)
	}
	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		return Transaction{}, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `insert into transactions(requestor_id, nonce, status, operation, metadata, created_at) values ($1,$2,$3,$4,$5,now())
	returning id, requestor_id, nonce, status, operation, created_at, metadata`,
		requestorId, nonce, "pending", operation, metaData)

	var l Transaction
	err = row.Scan(&l.Id, &l.RequestorId, &l.Nonce, &l.Status, &l.Operation, &l.CreatedAt, &l.MetaData)
	if err != nil {
		return Transaction{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return Transaction{}, err
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

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		err = utils.ToError(pgErr)
	}
	return err
}

func (r *Repo) updateTransactionStatus(ctx context.Context, tx pgx.Tx, id int64, status string) error {
	if tx == nil {
		return utils.NilTxError
	}
	_, err := tx.Exec(ctx, `
    UPDATE transactions
    SET status = $1 where id = $2`, status, id)
	return err
}

func (r *Repo) UpdateTransactionStatus(ctx context.Context, id int64, status string) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `
    UPDATE transactions
    SET status = $1 where id = $2`, status, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
