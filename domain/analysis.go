package domain

import (
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindService(func() *AnalysisService {
			return &AnalysisService{}
		})
		initiator.InjectController(func(ctx freedom.Context) (service *AnalysisService) {
			initiator.FetchService(ctx, &service)
			return
		})
	})
}

// AnalysisService .
type AnalysisService struct {
	Worker freedom.Worker
}
