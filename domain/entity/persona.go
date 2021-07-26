package entity

import (
	"fmt"
	"time"

	"github.com/8treenet/cdp-service/domain/po"
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
