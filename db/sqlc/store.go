package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB // support create a new db transaction.
}

type TransferTxParams struct {
	FromAccountID int64 `json:"form_Account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"form_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

var txKey = struct{}{} // creating new empty object struct{}

// execTx takes context and a callback function.
// It will start a new db transaction, create a new `Queries`
// call the callback function with the created `Queries`
// and finally commit || rollback the transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Start new transaction, returns transaction, err
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// New() get back `Queries` object
	// Instead pass in `sql.DB`, now passing `sql.Tx`
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err %v, rb err: %v\n", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTx take context, and params.
// It will handle transaction callback function then returns TransferTxResult struct and error.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")
		// Create a new transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			ToAccountID:   arg.ToAccountID,
			FromAccountID: arg.FromAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create subtract entry
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create add entry
		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Handle deadlock by update the smaller ID first. ( ID 1 -> ID 2 )
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.ToAccountID, +arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
