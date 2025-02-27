package controllers

import (
	"fmt"
	"strconv"

	"github.com/agilistikmal/bnnchat/src/services"
	"github.com/gofiber/fiber/v2"
)

type MenuController struct {
	MenuService *services.MenuService
}

func NewMenuController(menuService *services.MenuService) *MenuController {
	return &MenuController{
		MenuService: menuService,
	}
}

func (c *MenuController) Detail(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	menu, err := c.MenuService.FindMenuByID(id)
	if err != nil {
		return ctx.SendString(fmt.Sprintf("Error Menu Detail: %v", err.Error()))
	}

	binding := fiber.Map{
		"menu": menu,
	}
	return ctx.Render("pages/menu/detail", binding, "layouts/base")
}
