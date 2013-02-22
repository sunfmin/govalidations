package govalidations

import (
	"encoding/json"
)

type Error struct {
	Name    string
	Message string
}

type Errors []*Error

type Validated struct {
	Object interface{}
	Errors Errors
}

func (vd *Validated) HasError() bool {
	return len(vd.Errors) > 0
}

func (vd *Validated) AddError(name string, message string) {
	vd.Errors = append(vd.Errors, &Error{
		Name:    name,
		Message: message,
	})
	return
}

func (vd *Validated) Code() string {
	return "405"
}

func (vd *Validated) Error() string {
	return "validation error"
}

func (vd *Validated) ToError() (err error) {
	if vd.HasError() {
		return vd
	}
	return nil
}

func (vd *Validated) ToJson() (r []byte) {
	r, err := json.Marshal(vd)
	if err != nil {
		panic(err)
	}
	return
}

func (es Errors) Has(name string) bool {
	for _, e := range es {
		if e.Name == name {
			return true
		}
	}
	return false
}

func (es Errors) IfHasThen(name string, result string) (r string) {
	if es.Has(name) {
		return result
	}
	return
}

func (es Errors) On(name string) (r string) {
	for _, e := range es {
		if e.Name == name {
			r = e.Message
			return
		}
	}
	return
}
