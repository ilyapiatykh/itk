package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fasthttp/router"
	"github.com/google/uuid"
	"github.com/ilyapiatykh/itk/internal/models"
	"github.com/ilyapiatykh/itk/internal/repo"
	"github.com/valyala/fasthttp"
)

type pathParam = string

const _walletUUID pathParam = "walletUUID"

type WalletOperation struct {
	WalletID      uuid.UUID            `json:"walletId" validate:"required"`
	OperationType models.OperationType `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        float64              `json:"amount" validate:"required,gt=0"`
}

func (r *Router) registerWalletsRoutes(g *router.Group) {
	g.POST("/wallet", r.postWallet())
	g.GET(fmt.Sprintf("/wallets/{%s}", _walletUUID), r.getWallet())
}

func (r *Router) postWallet() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		var body WalletOperation

		err := json.Unmarshal(ctx.PostBody(), &body)
		if err != nil {
			sendError(ctx, "unsupported operation format", fasthttp.StatusBadRequest, nil)
			return
		}

		if err := r.validator.Struct(body); err != nil {
			sendError(ctx, err.Error(), fasthttp.StatusBadRequest, nil)
			return
		}

		wallet, err := r.service.UpdateWallet(ctx, body.WalletID, body.Amount, body.OperationType)
		if err != nil {
			if errors.Is(err, repo.ErrNegativeBalance) {
				sendError(ctx, err.Error(), fasthttp.StatusBadRequest, nil)
				return
			}

			sendError(ctx, "server error", fasthttp.StatusInternalServerError, err)
			return
		}

		sendJSON(ctx, wallet, fasthttp.StatusOK)
	}
}

func (r *Router) getWallet() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		rawID, _ := ctx.UserValue(_walletUUID).(string)
		id, err := uuid.Parse(rawID)
		if err != nil {
			sendError(ctx, "unsupported uuid format", fasthttp.StatusBadRequest, nil)
			return
		}

		wallet, err := r.service.GetWallet(ctx, id)
		if err != nil {
			if errors.Is(err, repo.ErrNoWallet) {
				sendError(ctx, "no such wallet", fasthttp.StatusBadRequest, nil)
				return
			}

			sendError(ctx, "server error", fasthttp.StatusInternalServerError, err)
			return
		}

		sendJSON(ctx, wallet, fasthttp.StatusOK)
	}
}

func sendJSON(ctx *fasthttp.RequestCtx, data any, status int) {
	ctx.SetStatusCode(status)
	ctx.SetContentType("application/json")

	body, err := json.Marshal(data)
	if err != nil {
		sendError(ctx, "server error", fasthttp.StatusInternalServerError, err)
		return
	}

	ctx.SetBody(body)
}

// sendError logs err and send response with passed description in JSON
func sendError(ctx *fasthttp.RequestCtx, description string, status int, err error) {
	if err != nil {
		slog.Error(
			"failed while processing request",
			slog.Any("error", err),
		)
	}

	body := struct {
		Description string `json:"description"`
	}{Description: description}

	sendJSON(ctx, body, status)
}
