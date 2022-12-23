package usecase

import (
	"testcase/server/internal/domain"
	"time"
)

type GetAverageValueUseCase struct {
	repo domain.Repository
}

func NewGetAverageValueUseCase(repo domain.Repository) *GetAverageValueUseCase {
	return &GetAverageValueUseCase{repo: repo}
}

func (x *GetAverageValueUseCase) Do(last time.Duration) (float32, error) {
	return x.repo.GetAverageValueTSForLast(last)
}
