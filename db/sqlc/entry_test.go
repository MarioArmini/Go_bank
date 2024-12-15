package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func newEntry(t *testing.T, accountId int64) Entries {
	arg := CreateEntryParams{
		AccountId: accountId,
		Amount:    50,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountId, entry.AccountId)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.Id)
	require.NotZero(t, entry.CreationTime)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := newAccount(t)
	newEntry(t, account.Id)
}

func TestUpdateEntry(t *testing.T) {
	account := newAccount(t)
	entry := newEntry(t, account.Id)

	arg := UpdateEntryParams{
		Id:     entry.Id,
		Amount: 65,
	}
	entryUpdated, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entryUpdated)

	require.Equal(t, entryUpdated.Id, entry.Id)
	require.Equal(t, entryUpdated.Amount, arg.Amount)
	require.Equal(t, entryUpdated.AccountId, entry.AccountId)

}

func TestDeleteEntry(t *testing.T) {
	account := newAccount(t)
	entry := newEntry(t, account.Id)

	err := testQueries.DeleteEntry(context.Background(), entry.Id)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry.Id)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestGetEntry(t *testing.T) {
	account := newAccount(t)
	entry := newEntry(t, account.Id)

	entry2, err := testQueries.GetEntry(context.Background(), entry.Id)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry2.Id, entry.Id)
	require.Equal(t, entry2.Amount, entry.Amount)
	require.Equal(t, entry2.AccountId, entry.AccountId)
}

func TestGetAllEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		account := newAccount(t)
		newEntry(t, account.Id)
	}

	arg := GetAllEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.GetAllEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, int(arg.Limit))

	for _, account := range entries {
		require.NotEmpty(t, account)
	}
}
