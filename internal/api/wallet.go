package api

import (
	"encoding/json"
	"fmt"

	"github.com/fasthttp/router"
	"github.com/google/uuid"
	"github.com/ilyapiatykh/itk/internal/models"
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
			sendError(ctx, "validation error", fasthttp.StatusBadRequest, err)
			return
		}

		wallet, err := r.service.UpdateWallet(ctx, body.WalletID, body.Amount, body.OperationType)
		if err != nil {
			sendError(ctx, "server error", fasthttp.StatusInternalServerError, err)
		}

		sendJSON(ctx, wallet, fasthttp.StatusOK)
	}
}

func (r *Router) getWallet() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		id, ok := ctx.UserValue(_walletUUID).(uuid.UUID)
		if !ok {
			sendError(ctx, "unsupported uuid format", fasthttp.StatusBadRequest, nil)
			return
		}

		wallet, err := r.service.GetWallet(ctx, id)
		if err != nil {
			sendError(ctx, "server error", fasthttp.StatusInternalServerError, err)
		}

		sendJSON(ctx, wallet, fasthttp.StatusOK)
	}
}

func sendJSON(ctx *fasthttp.RequestCtx, data any, statusCode int) {
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")

	body, err := json.Marshal(data)
	if err != nil {
		sendError(ctx, "server error", fasthttp.StatusInternalServerError, err)
		return
	}

	ctx.SetBody(body)
}

func sendError(ctx *fasthttp.RequestCtx, description string, statusCode int, err error) {
	// TODO log error
	body := struct {
		Description string `json:"description"`
	}{Description: description}

	sendJSON(ctx, body, statusCode)
}
