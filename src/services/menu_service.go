package services

import (
	"sort"

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
	err := s.DB.Preload("Options.SubMenu").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	for _, menu := range menus {
		sort.Slice(menu.Options, func(i, j int) bool {
			return menu.Options[i].Position < menu.Options[j].Position
		})
	}

	return menus, nil
}

func (s *MenuService) FindMenuByID(id int) (*models.Menu, error) {
	var menu *models.Menu
	err := s.DB.Preload("Options.SubMenu").Take(&menu, "id = ?", id).Order("Options.Position ASC").Error
	if err != nil {
		return nil, err
	}

	sort.Slice(menu.Options, func(i, j int) bool {
		return menu.Options[i].Position < menu.Options[j].Position
	})

	return menu, nil
}

func (s *MenuService) FindMenuBySlug(slug string) (*models.Menu, error) {
	var menu *models.Menu
	err := s.DB.Preload("Options.SubMenu").Take(&menu, "slug = ?", slug).Order("Options.Position ASC").Error
	if err != nil {
		return nil, err
	}

	sort.Slice(menu.Options, func(i, j int) bool {
		return menu.Options[i].Position < menu.Options[j].Position
	})

	return menu, nil
}

func (s *MenuService) FindOptionMenu(id int, position int) (*models.MenuOption, error) {
	var menuOption *models.MenuOption
	err := s.DB.Preload("SubMenu").Take(&menuOption, "menu_id = ? AND position = ?", id, position).Error
	if err != nil {
		return nil, err
	}

	return menuOption, nil
}
