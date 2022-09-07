package db

import (
	"context"
	"testing"
	"time"

	"github.com/qss888888/simple_bank/util"
	"github.com/stretchr/testify/assert"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	e, err := testQueries.CreateEntry(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, e)
	assert.Equal(t, arg.AccountID, e.AccountID)
	assert.Equal(t, arg.Amount, e.Amount)
	assert.NotZero(t, e.ID)
	assert.NotZero(t, e.CreatedAt)

	return e
}

func TestCreateEntry(t *testing.T) {
	a := createRandomAccount(t)
	createRandomEntry(t, a)
}

func TestGetEntry(t *testing.T) {
	a := createRandomAccount(t)
	e1 := createRandomEntry(t, a)

	e2, err := testQueries.GetEntry(context.Background(), e1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, e2)

	assert.Equal(t, e1.ID, e2.ID)
	assert.Equal(t, e1.AccountID, e2.AccountID)
	assert.Equal(t, e1.Amount, e2.Amount)
	assert.WithinDuration(t, e1.CreatedAt, e2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	a := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, a)
	}

	arg := ListEntriesParams{
		AccountID: a.ID,
		Limit:     5,
		Offset:    5,
	}

	e, err := testQueries.ListEntries(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, e, 5)

	for _, entry := range e {
		assert.NotEmpty(t, entry)
		assert.Equal(t, arg.AccountID, entry.AccountID)
	}
}
