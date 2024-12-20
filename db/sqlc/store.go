package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v - rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	SenderId    int64 `json:"senderId"`
	RecipientId int64 `json:"recipientId"`
	Amount      int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer       Transfers `json:"transfer"`
	Sender         Accounts  `json:"sender"`
	Recipient      Accounts  `json:"recipient"`
	SenderEntry    Entries   `json:"senderEntry"`
	RecipientEntry Entries   `json:"recipientEntry"`
}

func (store *Store) CreateTransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.SenderEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountId: arg.SenderId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.RecipientEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountId: arg.RecipientId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: update balances

		return nil
	})

	return result, err
}
