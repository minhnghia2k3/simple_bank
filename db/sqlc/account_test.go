package db

import (
	"context"
	"database/sql"
	"github.com/minhnghia2k3/simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	want := createRandomAccount(t)
	got, err := testQueries.GetAccount(context.Background(), want.ID)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Owner, got.Owner)
	require.Equal(t, want.Currency, got.Currency)
	require.Equal(t, want.Balance, got.Balance)
	require.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	want := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      want.ID,
		Balance: util.RandomBalance(),
	}

	got, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Owner, got.Owner)
	require.Equal(t, arg.Balance, got.Balance)
	require.Equal(t, want.Currency, got.Currency)
	require.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
