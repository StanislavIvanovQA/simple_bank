package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTransaction(ctx context.Context, args TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) executeTransaction(ctx context.Context, fn func(queries *Queries) error) error {
	transaction, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(transaction)
	err = fn(queries)
	if err != nil {
		if rollbackError := transaction.Rollback(); rollbackError != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rollbackError)
		}
	}

	return transaction.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var transactionKey = struct{}{}

func (store *SQLStore) TransferTransaction(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	var err error

	err = store.executeTransaction(ctx, func(queries *Queries) error {
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = transferMoney(
				ctx, queries, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = transferMoney(
				ctx, queries, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount,
			)
		}

		return nil
	})

	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
