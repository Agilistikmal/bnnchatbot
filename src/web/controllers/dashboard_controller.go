package controllers

import "github.com/gofiber/fiber/v2"

type DashboardController struct {
}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (c *DashboardController) Dashboard(ctx *fiber.Ctx) error {
	binding := fiber.Map{}
	return ctx.Render("pages/dashboard/index", binding, "layouts/base")
}
