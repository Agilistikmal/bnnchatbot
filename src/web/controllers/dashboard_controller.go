package controllers

import (
	"fmt"

	"github.com/agilistikmal/bnnchat/src/services"
	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow"
)

type DashboardController struct {
	WAClient    *whatsmeow.Client
	MenuService *services.MenuService
}

func NewDashboardController(waClient *whatsmeow.Client, menuService *services.MenuService) *DashboardController {
	return &DashboardController{
		WAClient:    waClient,
		MenuService: menuService,
	}
}

func (c *DashboardController) Dashboard(ctx *fiber.Ctx) error {
	binding := fiber.Map{}
	return ctx.Render("pages/dashboard/index", binding, "layouts/base")
}

func (c *DashboardController) MenuPart(ctx *fiber.Ctx) error {
	menus, err := c.MenuService.FindMenus()
	if err != nil {
		return ctx.SendString(fmt.Sprintf("Menu Part Error: %v", err.Error()))
	}

	binding := fiber.Map{
		"menus": menus,
	}
	return ctx.Render("pages/dashboard/menu_part", binding)
}

func (c *DashboardController) ChatPart(ctx *fiber.Ctx) error {
	binding := fiber.Map{
		"chats": "",
	}
	return ctx.Render("pages/dashboard/chat_part", binding, "layouts/base")
}
