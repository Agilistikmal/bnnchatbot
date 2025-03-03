package controllers

import (
	"fmt"
	"strconv"

	"github.com/agilistikmal/bnnchat/src/models"
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

func (c *MenuController) Add(ctx *fiber.Ctx) error {
	binding := fiber.Map{}

	switch ctx.Method() {
	case fiber.MethodPost:
		slug := ctx.FormValue("slug")
		header := ctx.FormValue("header")
		content := ctx.FormValue("content")
		footer := ctx.FormValue("footer")

		menu := &models.Menu{
			Slug:    slug,
			Header:  header,
			Content: content,
			Footer:  footer,
		}
		err := c.MenuService.DB.Create(&menu).Error
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error: %v", err.Error()))
		} else {
			return ctx.SendString("Berhasil menambah data")
		}
	default:
		return ctx.Render("pages/menu/add", binding, "layouts/base")
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

	switch ctx.Method() {
	case fiber.MethodPut:
		slug := ctx.FormValue("slug")
		header := ctx.FormValue("header")
		content := ctx.FormValue("content")
		footer := ctx.FormValue("footer")

		menu.Slug = slug
		menu.Header = header
		menu.Content = content
		menu.Footer = footer

		err := c.MenuService.DB.Updates(&menu).Error
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error: %v", err.Error()))
		} else {
			return ctx.SendString("Berhasil menyimpan data")
		}
	case fiber.MethodDelete:
		err := c.MenuService.DB.Delete(&menu, "id = ?", menu.ID).Error
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error: %v", err.Error()))
		} else {
			ctx.Append("HX-Redirect", "/")
			return ctx.SendString("Berhasil menghapus data")
		}
	default:
		return ctx.Render("pages/menu/detail", binding, "layouts/base")
	}
}
