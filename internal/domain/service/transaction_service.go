package service

import (
	"go-spe/internal/domain/models"
	"go-spe/internal/domain/repository"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetTransactionStatus(request_id string, bill_number string) (*models.Transaction, error) {
	return s.repo.GetTransactionStatus(request_id, bill_number)
}
