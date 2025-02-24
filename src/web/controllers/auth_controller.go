package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	req := ctx.Request()
	if string(req.Header.Method()) == fiber.MethodPost {
		return ctx.SendString("POST")
	}

	binding := fiber.Map{}
	return ctx.Render("pages/auth/login", binding, "layouts/auth")
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	binding := fiber.Map{}
	return ctx.Render("pages/auth/register", binding, "layouts/auth")
}
