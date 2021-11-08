package sms

import(
	"github.com/chnsz/golangsdk"
)


type Template struct{
	ID				string `json:"id"`
}


type commonResult struct {
	golangsdk.Result
}

func (r commonResult) Extract() (*Template, error) {
	var s Template
	err := r.ExtractInto(&s)
	return &s, err
}


type CreateResult struct {
	commonResult
}