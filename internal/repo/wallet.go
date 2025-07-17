package repo

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrNoWallet = errors.New("no wallet with such id")
var ErrNegativeBalance = errors.New("wallet cannot have negative balance")

const (
	queryGetWallet            = "SELECT balance FROM wallets WHERE id = $1"
	queryCreateOrUpdateWallet = `
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
	defer c.mu.RUnlock()

	balance, ok = c.wallets[id]
	return
}

func (c *cache) set(id uuid.UUID, balance float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.wallets[id] = balance
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
	if balance, ok := w.cache.get(id); ok {
		return balance, nil
	}

	err := w.db.QueryRowContext(ctx, queryGetWallet, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoWallet
		}

		return 0, err
	}

	w.cache.set(id, balance)

	return balance, nil
}

func (w *Wallets) UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) (float64, error) {
	var balance float64
	err := w.db.QueryRowContext(ctx, queryGetWallet, id, amount).Scan(&balance)
	if err != nil {
		// TODO
		return 0, ErrNegativeBalance
	}

	w.cache.set(id, balance)

	return balance, nil
}
