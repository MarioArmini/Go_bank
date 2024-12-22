package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := newAccount(t)
	account2 := newAccount(t)

	fmt.Println(">> before: ", account1.Balance, account2.Balance)

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
	existed := make(map[int]bool)

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

		senderAccount := result.Sender
		require.NotEmpty(t, senderAccount)
		require.Equal(t, senderAccount.Id, account1.Id)

		recipientAccount := result.Recipient
		require.NotEmpty(t, recipientAccount)
		require.Equal(t, recipientAccount.Id, account2.Id)

		fmt.Println(">> during transaction: ", senderAccount.Balance, recipientAccount.Balance)

		// check balances
		diff1 := account1.Balance - senderAccount.Balance
		diff2 := recipientAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check accounts' balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.Id)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.Id)
	require.NoError(t, err)

	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
