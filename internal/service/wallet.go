package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ilyapiatykh/itk/internal/models"
)

type WalletsStorage interface {
	GetBalance(ctx context.Context, id uuid.UUID) (float64, error)
	Deposit(ctx context.Context, id uuid.UUID, amount float64) (float64, error)
	Withdraw(ctx context.Context, id uuid.UUID, amount float64) (float64, error)
}

type Wallets struct {
	s WalletsStorage
}

func NewWallets(s WalletsStorage) *Wallets {
	return &Wallets{s: s}
}

func (w *Wallets) GetWallet(ctx context.Context, id uuid.UUID) (wallet models.Wallet, err error) {
	wallet.ID = id
	wallet.Balance, err = w.s.GetBalance(ctx, id)
	return
}

func (w *Wallets) UpdateWallet(ctx context.Context, id uuid.UUID, amount float64, operationType models.OperationType) (wallet models.Wallet, err error) {
	wallet.ID = id

	switch operationType {
	case models.Deposit:
		wallet.Balance, err = w.s.Deposit(ctx, id, amount)
	case models.Withdraw:
		wallet.Balance, err = w.s.Withdraw(ctx, id, amount)
	}

	return
}
