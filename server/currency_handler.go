package server

import (
	"context"
	"log"
	"strconv"

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

		c.JSON(ErrorResponse{
			ID:      "Latest",
			Message: err.Error(),
		})
		return
	}

	// Map results
	output := []map[string]interface{}{}
	for _, curr := range currencies {
		output = append(output, map[string]interface{}{
			"symbol":  curr.Symbol,
			"value":   curr.Value,
			"value_f": (float64(curr.Value) * 0.000001),
		})
	}

	c.JSON(output)
}

func (h *CurrencyHandler) BySymbol(c *fiber.Ctx) {
	symbol := c.Params("id")

	// Limit query param
	limit := 20
	limitString := c.Query("limit")
	if limitString != "" {
		limitInt, err := strconv.ParseInt(limitString, 10, 32)
		if err != nil {
			log.Println(err)
		} else {
			limit = int(limitInt)
		}
	}

	// Get currencies by symbol
	currencies, err := h.currencyService.BySymbol(symbol, limit)
	if err != nil {
		log.Println(err)

		c.JSON(ErrorResponse{
			ID:      "BySymbol",
			Message: err.Error(),
		})
		return
	}

	// Map results
	output := []map[string]interface{}{}
	for _, curr := range currencies {
		output = append(output, map[string]interface{}{
			"symbol":  curr.Symbol,
			"value":   curr.Value,
			"value_f": (float64(curr.Value) * 0.000001),
			"date":    curr.Date.Time,
		})
	}

	c.JSON(output)
}
