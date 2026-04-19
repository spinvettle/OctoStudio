package codexProxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyError(t *testing.T) {
	pool := NewAccountPool()
	account, err := pool.GetAccount()
	assert.ErrorIs(t, err, ErrEmptyPool, "expected ErrEmptyPool,get %v", err)
	assert.Nil(t, account, "account should be nil")
}

func setUpPool() *AccountPool {
	account1 := &Account{ID: "1", Name: "account_1", Status: Enabled, UsagePercent: 10.0}
	account2 := &Account{ID: "2", Name: "account_2", Status: Enabled, UsagePercent: 10.0}
	pool := NewAccountPool()
	_ = pool.AddAccount(account1)
	_ = pool.AddAccount(account2)
	return pool
}

func TestAddAccount(t *testing.T) {
	pool := setUpPool()

	assert.Len(t, pool.accountList, 2, "expected size of pool is 2")

}

func TestGetAccountById(t *testing.T) {
	pool := setUpPool()

	acc, err := pool.GetAccountById("1")
	assert.NoError(t, err, "expected nil err,but get err")
	assert.Equal(t, "1", acc.ID, "expected account id is 1,but get %s", acc.ID)
}

func TestGetAccount(t *testing.T) {
	pool := setUpPool()

	acc, err := pool.GetAccount()
	assert.NoError(t, err, "expected no err,but get err")
	assert.Equal(t, "1", acc.ID, "expected account id is 1,but get %s", acc.ID)
	acc.UpdateStatus(Disabled)

	acc, err = pool.GetAccount()
	assert.NoError(t, err, "expected no err,but get err")
	assert.Equal(t, "2", acc.ID, "expected account id is 1,but get %s", acc.ID)
	acc.UpdateStatus(Disabled)

	_, err = pool.GetAccount()
	assert.Error(t, err, ErrNotFoundAvailabelAccount, "expected no available account")

}
