package flavors

import (
	"github.com/chnsz/golangsdk/pagination"
)

type DbFlavorsResp struct {
	Flavorslist []Flavors `json:"flavors"`
}
type Flavors struct {
	ID           string            `json:"id" `
	Vcpus        string            `json:"vcpus" `
	Ram          int               `json:"ram" `
	Speccode     string            `json:"spec_code"  `
	Instancemode string            `json:"instance_mode" `
	Azstatus     map[string]string `json:"az_status" `
	VersionName  []string          `json:"version_name" `
	GroupType    string            `json:"group_type" `
}

type DbFlavorsPage struct {
	pagination.SinglePageBase
}

func (r DbFlavorsPage) IsEmpty() (bool, error) {
	data, err := ExtractDbFlavors(r)
	if err != nil {
		return false, err
	}
	return len(data.Flavorslist) == 0, err
}

func ExtractDbFlavors(r pagination.Page) (DbFlavorsResp, error) {
	var s DbFlavorsResp
	err := (r.(DbFlavorsPage)).ExtractInto(&s)
	return s, err
}
