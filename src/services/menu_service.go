package services

import (
	"github.com/agilistikmal/bnnchat/src/models"
	"gorm.io/gorm"
)

type MenuService struct {
	DB *gorm.DB
}

func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{
		DB: db,
	}
}

func (s *MenuService) FindMenus() ([]*models.Menu, error) {
	var menus []*models.Menu
	err := s.DB.Preload("Options.Menu").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (s *MenuService) FindMenuByID(id int) (*models.Menu, error) {
	var menu *models.Menu
	err := s.DB.Preload("Options.Menu").Take(&menu, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *MenuService) FindMenuBySlug(slug string) (*models.Menu, error) {
	var menu *models.Menu
	err := s.DB.Preload("Options.Menu").Take(&menu, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *MenuService) FindOptionMenu(id int, position int) (*models.MenuOption, error) {
	var menuOption *models.MenuOption
	err := s.DB.Preload("Menu").Take(&menuOption, "id = ? AND position = ?", id, position).Error
	if err != nil {
		return nil, err
	}

	return menuOption, nil
}
