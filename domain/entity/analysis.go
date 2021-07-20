package entity

import (
	"fmt"
	"time"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
)

// Analysis
type Analysis struct {
	freedom.Entity
	po.Analysis
	XMLBytes []byte
}

func (entity *Analysis) Identity() string {
	return fmt.Sprint(entity.ID)
}

func (entity *Analysis) GetBeginTime() time.Time {
	if entity.DateConservation != 0 {
		return time.Now().AddDate(0, 0, -entity.DateRange)
	}

	date := time.Now().Format("2006-01-02")
	currentDateTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return currentDateTime.AddDate(0, 0, -entity.DateRange)
}

func (entity *Analysis) GetEndTime() time.Time {
	if entity.DateConservation != 0 {
		return time.Now()
	}

	date := time.Now().Format("2006-01-02")
	currentDateTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return currentDateTime.Add(-1 * time.Second)
}
