package entity

import (
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
)

// Analysis
type Analysis struct {
	freedom.Entity
	po.Analysis
	XMLBytes []byte
}
