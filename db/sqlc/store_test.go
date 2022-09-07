package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	s := NewStore(testDB)

	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)

	// 交易次数
	times := 5
	// 每次交易金额
	amount := int64(10)
	// 交易错误管道
	errs := make(chan error, times)
	// 交易结果管道
	results := make(chan TransferTxResult, times)

	for i := 0; i < times; i++ {
		go func() {
			res, err := s.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: a1.ID,
				ToAccountID:   a2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- res

		}()
	}

	for i := 0; i < times; i++ {
		err := <-errs
		assert.NoError(t, err)

		res := <-results

		// 管道交易记录
		tf := res.Transfer
		assert.NotZero(t, tf.ID)
		assert.NotZero(t, tf.CreatedAt)
		assert.Equal(t, a1.ID, tf.FromAccountID)
		assert.Equal(t, a2.ID, tf.ToAccountID)
		assert.Equal(t, amount, tf.Amount)

		// 获取交易记录
		rtf, err := s.GetTransfer(context.Background(), tf.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, rtf)

		// 获取进项
		fe, err := s.GetEntry(context.Background(), res.FromEntry.ID)
		assert.NoError(t, err)
		assert.NotZero(t, fe.ID)
		assert.NotZero(t, fe.CreatedAt)
		assert.Equal(t, a1.ID, fe.AccountID)
		assert.Equal(t, -amount, fe.Amount)

		// 获取出项
		te, err := s.GetEntry(context.Background(), res.ToEntry.ID)
		assert.NoError(t, err)
		assert.NotZero(t, te.ID)
		assert.NotZero(t, te.CreatedAt)
		assert.Equal(t, a2.ID, te.AccountID)
		assert.Equal(t, amount, te.Amount)

		// 检查流出账户
		fa := res.FromAccount
		assert.NotEmpty(t, fa)
		assert.Equal(t, a1.ID, fa.ID)

		// 检查流入账户
		ta := res.ToAccount
		assert.NotEmpty(t, ta)
		assert.Equal(t, a2.ID, ta.ID)

		// 检查余额
		assert.Equal(t, (a1.Balance - fa.Balance), (ta.Balance - a2.Balance))

		// 打印每次交易后余额
		fmt.Println(">> tx:", fa.Balance, ta.Balance)
	}

	// 检查流出账户最终余额
	a11, err := s.GetAccount(context.Background(), a1.ID)
	assert.NoError(t, err)
	assert.Equal(t, a1.Balance, a11.Balance+amount*int64(times))

	// 检查流入账户最终余额
	a22, err := s.GetAccount(context.Background(), a2.ID)
	assert.NoError(t, err)
	assert.Equal(t, a2.Balance, a22.Balance-amount*int64(times))

	// 检查流入金额和流出金额是否相等
	assert.Equal(t, (a1.Balance - a11.Balance), (a22.Balance - a2.Balance))
}

func TestTransferTxDeadlock(t *testing.T) {

	s := NewStore(testDB)

	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)

	fmt.Println(">> before:", a1.Balance, a2.Balance)

	// 交易次数
	times := 10
	// 每次交易金额
	amount := int64(10)
	// 交易错误管道
	errs := make(chan error, times)

	for i := 0; i < times; i++ {
		fromID := a1.ID
		toID := a2.ID

		// id为单数，调转进项账户和出项账户
		if i%2 == 1 {
			fromID = a2.ID
			toID = a1.ID
		}

		go func() {
			_, err := s.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromID,
				ToAccountID:   toID,
				Amount:        amount,
			})

			errs <- err

		}()
	}

	// 排除错误
	for i := 0; i < times; i++ {
		err := <-errs
		assert.NoError(t, err)
	}

	// 检查流出账户最终余额
	a11, err := s.GetAccount(context.Background(), a1.ID)
	assert.NoError(t, err)
	assert.Equal(t, a1.Balance, a11.Balance)

	// 检查流入账户最终余额
	a22, err := s.GetAccount(context.Background(), a2.ID)
	assert.NoError(t, err)
	assert.Equal(t, a2.Balance, a22.Balance)

	fmt.Println(">> after:", a11.Balance, a22.Balance)
}
