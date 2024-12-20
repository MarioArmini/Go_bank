package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := newAccount(t)
	account2 := newAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.CreateTransferTx(context.Background(), TransferTxParams{
				SenderId:    account1.Id,
				RecipientId: account2.Id,
				Amount:      amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.Id, transfer.SenderId)
		require.Equal(t, account2.Id, transfer.RecipientId)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.Id)
		require.NotZero(t, transfer.CreationTime)

		_, err = store.GetTransfer(context.Background(), transfer.Id)
		require.NoError(t, err)

		senderEntry := result.SenderEntry
		require.NotEmpty(t, senderEntry)
		require.Equal(t, account1.Id, senderEntry.AccountId)
		require.Equal(t, -amount, senderEntry.Amount)
		require.NotZero(t, senderEntry.Id)
		require.NotZero(t, senderEntry.CreationTime)

		_, err = store.GetEntry(context.Background(), senderEntry.Id)
		require.NoError(t, err)

		recipientEntry := result.RecipientEntry
		require.NotEmpty(t, recipientEntry)
		require.Equal(t, account2.Id, recipientEntry.AccountId)
		require.Equal(t, amount, recipientEntry.Amount)
		require.NotZero(t, recipientEntry.Id)
		require.NotZero(t, recipientEntry.CreationTime)

		_, err = store.GetEntry(context.Background(), recipientEntry.Id)
		require.NoError(t, err)

		// TODO: check balances

	}
}
