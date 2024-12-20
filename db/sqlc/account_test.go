package db

import (
	"context"
	"database/sql"
	"go_bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newAccount(t *testing.T) Accounts {
	arg := CreateAccountParams{}

	arg.Owner = utils.GetRandomString(5)
	arg.Balance = utils.GetRandomInt()
	arg.Currency = utils.GetRandomCurrency()
	arg.InterestRate = utils.GetRandomInterestRate(0, 0.4)

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.Id)
	require.NotZero(t, account.CreationTime)

	require.Equal(t, account.InterestRate, arg.InterestRate)

	return account
}

func TestCreateAccount(t *testing.T) {
	newAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := newAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.Id)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Id, account2.Id)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreationTime, account2.CreationTime, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := newAccount(t)
	arg := UpdateAccountParams{
		Id:      account.Id,
		Balance: 100,
	}

	accountUpdated, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)

	require.Equal(t, accountUpdated.Id, account.Id)
	require.Equal(t, accountUpdated.Balance, arg.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account := newAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.Id)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account.Id)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestGetAllAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		newAccount(t)
	}

	arg := GetAllAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.GetAllAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, int(arg.Limit))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestGetBalance(t *testing.T) {
	account := newAccount(t)

	balance, err := testQueries.GetBalance(context.Background(), account.Id)
	require.NoError(t, err)
	require.NotEmpty(t, balance)
}
