package db

import (
	"context"
	"github.com/minhnghia2k3/simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateTransfer(t *testing.T) Transfer {
	arg := CreateTransferParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        util.RandomBalance(),
	}
	entry, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.FromAccountID, entry.FromAccountID)
	require.Equal(t, arg.ToAccountID, entry.ToAccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestCreateTransfer(t *testing.T) {
	CreateTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	want := CreateTransfer(t)
	got, err := testQueries.GetTransfer(context.Background(), want.ID)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.ToAccountID, got.ToAccountID)
	require.Equal(t, want.FromAccountID, got.FromAccountID)
	require.Equal(t, want.Amount, got.Amount)
	require.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Second)
}
func TestListTransfer(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateTransfer(t)
	}

	arg := ListTransfersParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
