package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var ErrNoWallet = errors.New("no wallet with such id")
var ErrNegativeBalance = errors.New("wallet cannot have negative balance")

const (
	queryGetWallet = "SELECT balance FROM wallets WHERE id = $1"
	queryWithdraw  = `
		UPDATE wallets
			SET balance = wallets.balance - $2
			WHERE id = $1 AND wallets.balance - $2 >= 0
			RETURNING balance;`
	queryDeposit = `
		INSERT INTO wallets (id, balance) VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE
			SET balance = wallets.balance + $2
			RETURNING balance;`
)

type cache struct {
	mu      sync.RWMutex
	wallets map[uuid.UUID]float64
}

func (c *cache) get(id uuid.UUID) (balance float64, ok bool) {
	c.mu.RLock()
	balance, ok = c.wallets[id]
	c.mu.RUnlock() // without defer faster and more dangerous
	return
}

func (c *cache) set(id uuid.UUID, balance float64) {
	c.mu.Lock()
	c.wallets[id] = balance
	c.mu.Unlock()
}

type Wallets struct {
	cache *cache
	db    *sql.DB
}

func NewWallets(db *sql.DB) *Wallets {
	cache := &cache{wallets: make(map[uuid.UUID]float64)}
	return &Wallets{cache: cache, db: db}
}

func (w *Wallets) GetBalance(ctx context.Context, id uuid.UUID) (float64, error) {
	var balance float64

	balance, ok := w.cache.get(id)
	if ok {
		return balance, nil
	}

	if err := w.db.QueryRowContext(ctx, queryGetWallet, id).Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoWallet
		}

		return 0, fmt.Errorf("getting balance from db: %v", err)
	}

	w.cache.set(id, balance)

	return balance, nil
}

func (w *Wallets) Withdraw(ctx context.Context, id uuid.UUID, amount float64) (float64, error) {
	var balance float64
	if err := w.db.QueryRowContext(ctx, queryWithdraw, id, amount).Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNegativeBalance
		}

		return 0, fmt.Errorf("withdraw from wallet in db: %v", err)
	}

	w.cache.set(id, balance)

	return balance, nil
}

func (w *Wallets) Deposit(ctx context.Context, id uuid.UUID, amount float64) (float64, error) {
	var balance float64
	if err := w.db.QueryRowContext(ctx, queryDeposit, id, amount).Scan(&balance); err != nil {
		return 0, fmt.Errorf("deposit to wallet in db: %v", err)
	}

	w.cache.set(id, balance)

	return balance, nil
}
