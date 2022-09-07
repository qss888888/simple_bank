package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/qss888888/simple_bank/util"
	"github.com/stretchr/testify/assert"
)

// 创建随机账户
func createRandomAccount(t *testing.T) Account {
	// 创建账户参数
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// 数据库查询账户
	a, err := testQueries.CreateAccount(context.Background(), arg)

	assert.NoError(t, err)
	assert.Equal(t, arg.Owner, a.Owner)
	assert.Equal(t, arg.Balance, a.Balance)
	assert.Equal(t, arg.Currency, a.Currency)
	assert.NotZero(t, a.ID)
	assert.NotZero(t, a.CreatedAt)

	return a

}

// 测试创建账户
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

// 测试获取账户
func TestGetAccount(t *testing.T) {
	a := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), a.ID)

	assert.NoError(t, err)
	assert.Equal(t, a.ID, account.ID)
	assert.Equal(t, a.Owner, account.Owner)
	assert.Equal(t, a.Balance, account.Balance)
	assert.Equal(t, a.Currency, account.Currency)
	assert.WithinDuration(t, a.CreatedAt, account.CreatedAt, time.Second)

}

// 测试更新账户
func TestUpdateAccount(t *testing.T) {
	a := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      a.ID,
		Balance: util.RandomMoney(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)

	assert.NoError(t, err)
	assert.Equal(t, a.ID, account.ID)
	assert.Equal(t, a.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, a.Currency, account.Currency)
	assert.WithinDuration(t, a.CreatedAt, account.CreatedAt, time.Second)
}

// 测试删除账户
func TestDeleteAccount(t *testing.T) {
	a := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), a.ID)

	assert.NoError(t, err)

	a, err = testQueries.GetAccount(context.Background(), a.ID)

	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, a)
}

// 测试获取用户列表
func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	assert.NoError(t, err)

	for _, account := range accounts {
		assert.NotEmpty(t, account.ID)
		assert.NotEmpty(t, account.Owner)
		assert.NotEmpty(t, account.Balance)
		assert.NotEmpty(t, account.Currency)
		assert.NotEmpty(t, account.CreatedAt)
	}

}
