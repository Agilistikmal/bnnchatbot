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
	err := s.DB.Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (s *MenuService) FindMenu(id string) (*models.Menu, error) {
	var menu *models.Menu
	err := s.DB.Take(&menu, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return menu, nil
}
