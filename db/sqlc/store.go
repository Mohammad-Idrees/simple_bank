package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the another
// It creates a transfer record, and new account entries and update account balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			arg.FromAccountId,
			arg.ToAccountId,
			arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// get account -> update its balance
		/*
			fmt.Println(txName, "get account 1")
			account1, err := store.GetAccountForUpdate(ctx, arg.FromAccountId)
			if err != nil {
				return err
			}
			fmt.Println(txName, "update account 1")
			result.FromAccount, err = store.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountId,
				Balance: account1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}

			fmt.Println(txName, "get account 2")
			account2, err := store.GetAccountForUpdate(ctx, arg.ToAccountId)
			if err != nil {
				return err
			}
			fmt.Println(txName, "update account 2")
			result.ToAccount, err = store.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountId,
				Balance: account2.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}
		*/

		result.FromAccount, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: -arg.Amount,
			ID:     arg.FromAccountId,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: arg.Amount,
			ID:     arg.ToAccountId,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
