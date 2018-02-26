package flavors

import "github.com/huaweicloud/golangsdk"

type Flavor struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Ram      int    `json:"ram"`
	SpecCode string `json:"specCode"`
}

type ListResult struct {
	golangsdk.Result
}

func (lr ListResult) Extract() ([]Flavor, error) {
	var a struct {
		Flavors []Flavor `json:"flavors"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Flavors, err
}
