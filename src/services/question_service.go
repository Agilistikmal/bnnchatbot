package services

import (
	"github.com/agilistikmal/bnnchat/src/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuestionService struct {
	DB *gorm.DB
}

func NewQuestionService(db *gorm.DB) *QuestionService {
	return &QuestionService{
		DB: db,
	}
}

func (s *QuestionService) FindQuestions() ([]*models.Question, error) {
	var questions []*models.Question
	err := s.DB.Find(&questions).Error
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func (s *QuestionService) FindQuestion(id string) (*models.Question, error) {
	var question *models.Question
	err := s.DB.Take(&question, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (s *QuestionService) SearchQuestion(triggers []string) (*models.Question, error) {
	var question *models.Question
	err := s.DB.Take(&question, "triggers && ?", pq.Array(triggers)).Error
	if err != nil {
		return nil, err
	}

	return question, nil
}
