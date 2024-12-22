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

		// this check is done to avoid deadlocks
		if arg.SenderId < arg.RecipientId {
			result.Sender, result.Recipient, err = transferMoney(ctx, q, arg.SenderId, -arg.Amount, arg.RecipientId, arg.Amount)
		} else {
			result.Recipient, result.Sender, err = transferMoney(ctx, q, arg.RecipientId, arg.Amount, arg.SenderId, -arg.Amount)
		}

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	accountId1 int64,
	amount1 int64,
	accountId2 int64,
	amount2 int64,
) (account1 Accounts, account2 Accounts, err error) {
	account1, err = q.GetAccountForUpdate(ctx, accountId1)
	if err != nil {
		return
	}

	account1, err = q.UpdateAccount(ctx, UpdateAccountParams{
		Id:      accountId1,
		Balance: account1.Balance + amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.GetAccountForUpdate(ctx, accountId2)
	if err != nil {
		return
	}

	account2, err = q.UpdateAccount(ctx, UpdateAccountParams{
		Id:      accountId2,
		Balance: account2.Balance + amount2,
	})

	return
}
