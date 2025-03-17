package controllers

import (
	"time"

	"github.com/agilistikmal/bnnchat/src/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HelpController struct {
	DB *gorm.DB
}

func NewHelpController(db *gorm.DB) *HelpController {
	return &HelpController{
		DB: db,
	}
}

func (c *HelpController) HelpPart(ctx *fiber.Ctx) error {
	var helps []models.Help
	c.DB.Find(&helps)

	binding := fiber.Map{
		"helps":       helps,
		"lastRefresh": time.Now().Format("15:04:05"),
	}
	return ctx.Render("pages/help/help_part", binding)
}

func (c *HelpController) Delete(ctx *fiber.Ctx) error {
	jid := ctx.Params("jid")

	err := c.DB.Where("j_id = ?", jid).Delete(&models.Help{}).Error
	if err != nil {
		return ctx.SendString("Error: " + err.Error())
	}

	return ctx.SendString("Success: " + jid)
}
