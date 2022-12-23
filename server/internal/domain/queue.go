package domain

import (
	"testcase/server/internal/entity/models"
)

type Queue interface {
	Add(v models.Timestamp)
}
