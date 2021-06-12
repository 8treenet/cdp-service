package infra

import (
	"errors"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindInfra(false, func() *Termination {
			return &Termination{}
		})
		initiator.InjectController(func(ctx freedom.Context) (com *Termination) {
			initiator.FetchInfra(ctx, &com)
			return
		})
	})
}

// Termination .
type Termination struct {
	freedom.Infra
}

// Fuck code:9999.
func (response *Termination) Fuck() error {
	setErrorCode(response.Worker(), 9999)
	return errors.New("fuck")
}

// Custom  透传上游
func (response *Termination) Custom(code int, msg string) error {
	setErrorCode(response.Worker(), code)
	return errors.New(msg)
}
