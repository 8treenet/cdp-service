//Package infra generated by 'freedom new-project github.com/8treenet/crm-service'
package infra

import (
	"io/ioutil"
	"reflect"

	"encoding/json"

	"github.com/8treenet/freedom"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindInfra(false, func() *Request {
			return &Request{}
		})
		initiator.InjectController(func(ctx freedom.Context) (com *Request) {
			initiator.FetchInfra(ctx, &com)
			return
		})
	})
}

// Request .
type Request struct {
	freedom.Infra
}

// BeginRequest .
func (req *Request) BeginRequest(worker freedom.Worker) {
	req.Infra.BeginRequest(worker)
}

// ReadJSON .
func (req *Request) ReadJSON(obj interface{}, validate ...bool) error {
	rawData, err := ioutil.ReadAll(req.Worker().IrisContext().Request().Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(rawData, obj); err != nil {
		return err
	}
	if len(validate) == 0 || !validate[0] {
		return nil
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			if err = req.Validate(val.Index(i).Interface()); err != nil {
				return err
			}
		}
		return nil
	}

	return req.Validate(obj)
}

// Validate .
func (req *Request) Validate(obj interface{}) error {
	return validate.Struct(obj)
}

// ReadQuery .
func (req *Request) ReadQuery(obj interface{}) error {
	if err := req.Worker().IrisContext().ReadQuery(obj); err != nil {
		return err
	}
	return validate.Struct(obj)
}

// ReadForm .
func (req *Request) ReadForm(obj interface{}) error {
	if err := req.Worker().IrisContext().ReadForm(obj); err != nil {
		return err
	}
	return validate.Struct(obj)
}
