package db

import (
	"context"
	"testing"
	"time"

	"github.com/qss888888/simple_bank/util"
	"github.com/stretchr/testify/assert"
)

func createRandomTransfer(t *testing.T, a1, a2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: a1.ID,
		ToAccountID:   a2.ID,
		Amount:        util.RandomMoney(),
	}

	tf, err := testQueries.CreateTransfer(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotZero(t, tf.ID)
	assert.NotZero(t, tf.CreatedAt)
	assert.Equal(t, arg.FromAccountID, tf.FromAccountID)
	assert.Equal(t, arg.ToAccountID, tf.ToAccountID)
	assert.Equal(t, arg.Amount, tf.Amount)

	return tf
}

func TestCreateTransfer(t *testing.T) {
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)
	createRandomTransfer(t, a1, a2)
}

func TestGetTransfer(t *testing.T) {
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)
	t1 := createRandomTransfer(t, a1, a2)

	t2, err := testQueries.GetTransfer(context.Background(), t1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, t2)

	assert.Equal(t, t1.ID, t2.ID)
	assert.Equal(t, t1.FromAccountID, t2.FromAccountID)
	assert.Equal(t, t1.ToAccountID, t2.ToAccountID)
	assert.Equal(t, t1.Amount, t2.Amount)
	assert.WithinDuration(t, t1.CreatedAt, t2.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, a1, a2)
		createRandomTransfer(t, a2, a1)
	}

	arg := ListTransfersParams{
		FromAccountID: a1.ID,
		ToAccountID:   a1.ID,
		Limit:         5,
		Offset:        5,
	}

	tfs, err := testQueries.ListTransfers(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, tfs, 5)

	for _, tf := range tfs {
		assert.NotEmpty(t, tf)
		assert.True(t, tf.FromAccountID == a1.ID || tf.ToAccountID == a1.ID)
	}
}
