package server

import (
	"context"
	"log"

	"github.com/gofiber/fiber"

	"github.com/matiss/go-microservice-test/services"
)

type CurrencyHandler struct {
	ctx             context.Context
	currencyService *services.CurrencyService
}

func NewCurrencyHandler(ctx context.Context, currencyService *services.CurrencyService) *CurrencyHandler {
	handler := CurrencyHandler{
		ctx:             ctx,
		currencyService: currencyService,
	}

	return &handler
}

func (h *CurrencyHandler) Latest(c *fiber.Ctx) {
	currencies, err := h.currencyService.Latest()
	if err != nil {
		log.Println(err)
	}

	c.JSON(currencies)
}

func (h *CurrencyHandler) BySymbol(c *fiber.Ctx) {
	symbol := c.Params("id")

	currencies, err := h.currencyService.BySymbol(symbol, 20)
	if err != nil {
		log.Println(err)
	}

	c.JSON(currencies)
}
