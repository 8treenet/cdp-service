package infra

import "github.com/8treenet/freedom"

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindInfra(false, func() *ResponseCode {
			return &ResponseCode{}
		})
		initiator.InjectController(func(ctx freedom.Context) (com *ResponseCode) {
			initiator.FetchInfra(ctx, &com)
			return
		})
	})
}

// ResponseCode .
type ResponseCode struct {
	freedom.Infra
}

// SetFuckCode .
func (response *ResponseCode) SetFuckCode() {
	setErrorCode(response.Worker(), 9999)
}
