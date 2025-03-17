package controllers

import (
	"fmt"

	"github.com/agilistikmal/bnnchat/src/lib"
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
	if c.WAClient != nil {
		if c.WAClient.Store.ID != nil {
			binding["waID"] = c.WAClient.Store.ID
		}
	}
	return ctx.Render("pages/dashboard/index", binding, "layouts/base")
}

func (c *DashboardController) QrCode(ctx *fiber.Ctx) error {
	return ctx.SendString(`<img src="/public/qr.png">`)
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

func (c *DashboardController) Logout(ctx *fiber.Ctx) error {
	c.WAClient.Logout()
	lib.RestartApp()
	ctx.Append("HX-Redirect", "/")
	return ctx.SendString("Berhasil Keluar WhatsApp")
}
