package web

import "github.com/gofiber/fiber/v2"

type WebHandler struct {
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

func (h *WebHandler) HomePage(ctx *fiber.Ctx) error {
	binding := fiber.Map{}
	return ctx.Render("home", binding)
}
