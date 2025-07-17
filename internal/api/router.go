package api

import (
	"context"

	"github.com/fasthttp/router"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ilyapiatykh/itk/config"
	"github.com/ilyapiatykh/itk/internal/models"
	"github.com/valyala/fasthttp"
)

type walletsProvider interface {
	GetWallet(context.Context, uuid.UUID) (models.Wallet, error)
	UpdateWallet(context.Context, uuid.UUID, float64, models.OperationType) (models.Wallet, error)
}

type serviceProvider interface {
	walletsProvider
}

type Router struct {
	service   serviceProvider
	server    *fasthttp.Server
	validator *validator.Validate
	port      string
}

func NewRouter(cfg *config.Server, service serviceProvider) *Router {
	r := router.New()

	router := &Router{
		service: service,
		server: &fasthttp.Server{
			ReadTimeout: cfg.ReadTimeout,
			Handler:     r.Handler,
		},
		validator: validator.New(),
		port: cfg.Port,
	}

	r.GET("/status", func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	})

	api := r.Group("/api")

	v1 := api.Group("/v1")
	router.registerWalletsRoutes(v1)

	return router
}

func (r *Router) Start() error {
	return r.server.ListenAndServe(r.port)
}

func (r *Router) Stop(ctx context.Context) error {
	return r.server.ShutdownWithContext(ctx)
}
