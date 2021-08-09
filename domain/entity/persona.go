package entity

import (
	"fmt"
	"time"

	"cdp-service/domain/po"

	"github.com/8treenet/freedom"
)

// Persona
type Persona struct {
	freedom.Entity
	po.Persona
	XMLBytes []byte
}

func (entity *Persona) Identity() string {
	return fmt.Sprint(entity.ID)
}

func (entity *Persona) GetBeginTime(now time.Time) time.Time {
	return now.AddDate(0, 0, -entity.DateRange)
}

func (entity *Persona) GetDeadHour() int {
	if entity.Deleted == 0 {
		return 0
	}
	return int(time.Now().Sub(entity.Updated).Hours())
}
