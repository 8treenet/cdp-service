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

func (entity *Analysis) GetBeginTime(now time.Time) time.Time {
	if entity.DateConservation != 0 {
		return now.AddDate(0, 0, -entity.DateRange)
	}

	date := now.Format("2006-01-02")
	currentDateTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return currentDateTime.AddDate(0, 0, -entity.DateRange)
}

func (entity *Analysis) GetEndTime(now time.Time) time.Time {
	if entity.DateConservation != 0 {
		return now.Add(-20 * time.Minute) //endtime固定推后20分钟
	}

	date := now.Format("2006-01-02")
	currentDateTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return currentDateTime.Add(-1 * time.Second) //昨日23:59:59
}
