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
			return ctx.SendString(fmt.Sprintf("Berhasil menambah data. <a href='/menu/%d' class='text-blue-500'>Lihat %s disini.</a>", menu.ID, menu.Slug))
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

	menus, err := c.MenuService.FindMenus()
	if err != nil {
		return ctx.SendString(fmt.Sprintf("Error Menu List: %v", err.Error()))
	}

	// Dapatkan next position untuk form tambah sub menu
	nextPosition, err := c.MenuService.GetLastPosition(id)
	if err != nil {
		nextPosition = 1 // fallback ke 1 jika terjadi error
	}

	binding := fiber.Map{
		"menu":         menu,
		"menus":        menus,
		"nextPosition": nextPosition,
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
		err := c.MenuService.DeleteMenu(menu.ID)
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error: %v", err.Error()))
		}

		ctx.Append("HX-Redirect", "/")
		return ctx.SendString("Menu berhasil dihapus beserta semua opsi/sub menu terkait")
	default:
		return ctx.Render("pages/menu/detail", binding, "layouts/base")
	}
}

func (c *MenuController) SubMenu(ctx *fiber.Ctx) error {
	menuID, _ := strconv.Atoi(ctx.Params("menuID"))
	menu, err := c.MenuService.FindMenuByID(menuID)
	if err != nil {
		return ctx.SendString(fmt.Sprintf("Error Menu Detail: %v", err.Error()))
	}

	switch ctx.Method() {
	case fiber.MethodPost:
		position, _ := strconv.Atoi(ctx.FormValue("position"))
		subMenuID, _ := strconv.Atoi(ctx.FormValue("sub_menu_id"))
		subMenu, err := c.MenuService.FindMenuByID(subMenuID)
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error SubMenu Detail: %v", err.Error()))
		}

		option := &models.MenuOption{
			MenuID:    menu.ID,
			SubMenuID: subMenu.ID,
			Position:  position,
		}
		err = c.MenuService.DB.Create(&option).Error
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error Option Create: %v", err.Error()))
		} else {
			ctx.Append("HX-Redirect", fmt.Sprintf("/menu/%d", menu.ID))
			return ctx.SendString("Berhasil menambah opsi/sub menu")
		}
	case fiber.MethodDelete:
		subMenuId, _ := strconv.Atoi(ctx.FormValue("sub_menu_id"))
		subMenu, err := c.MenuService.FindMenuByID(subMenuId)
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error Menu Detail: %v", err.Error()))
		}

		var option *models.MenuOption
		err = c.MenuService.DB.Delete(&option, "menu_id = ? and sub_menu_id = ?", menu.ID, subMenu.ID).Error
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error Delete SubMenu: %v", err.Error()))
		}

		ctx.Append("HX-Redirect", fmt.Sprintf("/menu/%d", menu.ID))
		return ctx.SendString("Berhasil menghapus opsi/sub menu")
	default:
		return nil
	}
}

func (c *MenuController) SubMenuPosition(ctx *fiber.Ctx) error {
	menuID, _ := strconv.Atoi(ctx.Params("menuID"))
	menu, err := c.MenuService.FindMenuByID(menuID)
	if err != nil {
		return ctx.SendString(fmt.Sprintf("Error Menu Detail: %v", err.Error()))
	}

	switch ctx.Method() {
	case fiber.MethodPut:
		optionID, _ := strconv.Atoi(ctx.FormValue("option_id"))
		position, _ := strconv.Atoi(ctx.FormValue("position"))
		subMenuID, _ := strconv.Atoi(ctx.FormValue("sub_menu_id"))
		subMenu, err := c.MenuService.FindMenuByID(subMenuID)
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error SubMenu Detail: %v", err.Error()))
		}

		option := &models.MenuOption{
			ID:        optionID,
			MenuID:    menu.ID,
			SubMenuID: subMenu.ID,
			Position:  position,
		}
		err = c.MenuService.DB.Save(&option).Error
		if err != nil {
			return ctx.SendString(fmt.Sprintf("Error Option Update: %v", err.Error()))
		} else {
			return ctx.SendString(fmt.Sprintf("Berhasil mengubah posisi %s menjadi %d", subMenu.Slug, position))
		}
	default:
		return nil
	}
}
