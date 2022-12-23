package usecase

import (
	"testcase/server/internal/domain"
	"testcase/server/internal/entity/models"
)

type InsertBatchUseCase struct {
	repo domain.Repository
}

func NewInsertBatchUseCase(repo domain.Repository) *InsertBatchUseCase {
	return &InsertBatchUseCase{repo: repo}
}

func (x *InsertBatchUseCase) Do(list []models.Timestamp) error {
	return x.repo.InsertBatchTS(list)
}
