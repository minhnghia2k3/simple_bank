package db

import (
	"context"
	"github.com/minhnghia2k3/simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: 1,
		Amount:    util.RandomBalance(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateEntry(t)
}
func TestGetEntry(t *testing.T) {
	want := CreateEntry(t)
	got, err := testQueries.GetEntry(context.Background(), want.ID)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.AccountID, got.AccountID)
	require.Equal(t, want.Amount, got.Amount)
	require.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Second)
}
func TestListEntry(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateEntry(t)
	}

	arg := ListEntriesParams{
		AccountID: 1,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
