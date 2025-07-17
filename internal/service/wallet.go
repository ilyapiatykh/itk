package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/ilyapiatykh/itk/internal/models"
)

type WalletsStorage interface {
	GetBalance(ctx context.Context, id uuid.UUID) (float64, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) (float64, error)
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
	if err != nil {
		slog.Error(
			"failed to load balance",
			slog.String("id", id.String()),
			slog.Any("error", err),
		)
	}

	return
}

func (w *Wallets) UpdateWallet(ctx context.Context, id uuid.UUID, amount float64, operationType models.OperationType) (wallet models.Wallet, err error) {
	if operationType == models.Withdraw {
		amount = -amount
	}

	wallet.ID = id
	wallet.Balance, err = w.s.UpdateBalance(ctx, id, amount)
	if err != nil {
		slog.Error(
			"failed to update balance",
			slog.String("id", id.String()),
			slog.Any("error", err),
		)
	}

	return
}
