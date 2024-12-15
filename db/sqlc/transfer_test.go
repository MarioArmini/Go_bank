package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTransfer(t *testing.T, senderId int64, recipiendId int64) Transfers {
	arg := CreateTransferParams{
		SenderId:    senderId,
		RecipientId: recipiendId,
		Amount:      50,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.SenderId, transfer.SenderId)
	require.Equal(t, arg.RecipientId, transfer.RecipientId)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.Id)
	require.NotZero(t, transfer.CreationTime)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	sender := newAccount(t)
	recipient := newAccount(t)
	newTransfer(t, sender.Id, recipient.Id)
}

func TestGetTransfer(t *testing.T) {
	sender := newAccount(t)
	recipient := newAccount(t)
	transfer := newTransfer(t, sender.Id, recipient.Id)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.Id)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer2.Id, transfer.Id)
	require.Equal(t, transfer2.SenderId, transfer.SenderId)
	require.Equal(t, transfer2.RecipientId, transfer.RecipientId)
	require.Equal(t, transfer2.Amount, transfer.Amount)
}

func TestGetAllTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		sender := newAccount(t)
		recipient := newAccount(t)
		newTransfer(t, sender.Id, recipient.Id)
	}

	arg := GetAllTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.GetAllTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, int(arg.Limit))

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestDeleteTransfer(t *testing.T) {
	sender := newAccount(t)
	recipient := newAccount(t)
	transfer := newTransfer(t, sender.Id, recipient.Id)

	err := testQueries.DeleteTransfer(context.Background(), transfer.Id)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.Id)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}
