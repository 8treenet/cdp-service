package entity

import (
	"fmt"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
)

type Report struct {
	po.AnalysisReport
	freedom.Entity
}

func (entity *Report) Identity() string {
	return fmt.Sprint(entity.ID)
}
