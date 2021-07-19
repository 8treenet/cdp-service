package repository

import (
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *AnalysisRepository {
			return &AnalysisRepository{}
		})
	})
}

// AnalysisRepository .
type AnalysisRepository struct {
	SignRepository
	Common *infra.CommonRequest
}
