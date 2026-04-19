package codexProxy

import (
	"sync"
)

type AccountStatus int
type LoginMode string
type AccoutOption func(*Account) error

type AccountPool struct {
	accountsMap map[string]*Account
	mu          sync.RWMutex
	accountList []*Account
	nextIndex   int
}

func NewAccountPool() *AccountPool {
	Pool := &AccountPool{
		accountsMap: make(map[string]*Account),
		mu:          sync.RWMutex{},
	}
	return Pool
}
func (pool *AccountPool) GetAccountById(id string) (*Account, error) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	if pool.accountsMap == nil {
		return nil, ErrEmptyPool
	}
	account := pool.accountsMap[id]
	if account == nil {
		return nil, ErrAccountNotExists
	}
	return account, nil
}
func (pool *AccountPool) AddAccount(account *Account) error {

	pool.mu.Lock()
	if pool.accountsMap == nil {
		pool.accountsMap = make(map[string]*Account)
	}
	pool.accountsMap[account.ID] = account
	if account.UsagePercent > 0 {
		account.Status = Enabled
	} else {
		account.Status = Disabled
	}
	pool.accountList = append(pool.accountList, account)
	pool.mu.Unlock()

	return nil
}
func (pool *AccountPool) GetAllAccounts() *[]AccountSnap {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	list := make([]AccountSnap, 0, len(pool.accountsMap))
	for _, acc := range pool.accountsMap {
		list = append(list, acc.SnapShot())
	}
	return &list

}
func (pool *AccountPool) GetAccount() (*Account, error) {

	pool.mu.Lock()
	defer pool.mu.Unlock()
	var targetAccount *Account
	if len(pool.accountsMap) == 0 || len(pool.accountList) == 0 {
		return nil, ErrEmptyPool
	}
	for range len(pool.accountList) {
		targetAccount = pool.accountList[pool.nextIndex]

		pool.nextIndex = (pool.nextIndex + 1) % len(pool.accountList)

		if targetAccount != nil {
			if targetAccount.Status == Enabled {
				return targetAccount, nil
			}
		}

	}
	return nil, ErrNotFoundAvailabelAccount

}
