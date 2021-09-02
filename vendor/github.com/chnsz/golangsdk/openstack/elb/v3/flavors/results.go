package flavors

import (
	"github.com/chnsz/golangsdk/pagination"
)

type Flavor struct {
	// Specifies the ID of the flavor.
	ID string `json:"id"`

	// Specifies the info of the flavor.
	Info FlavorInfo `json:"info"`

	// Specifies the name of the flavor.
	Name string `json:"name"`

	// Specifies whether shared.
	Shared bool `json:"shared"`

	// Specifies the type of the flavor.
	Type string `json:"type"`

	// Specifies whether sold out.
	SoldOut bool `json:"flavor_sold_out"`
}

type FlavorInfo struct {
	// Specifies the connection
	Connection int `json:"connection"`

	// Specifies the cps.
	Cps int `json:"cps"`

	// Specifies the qps
	Qps int `json:"qps"`

	// Specifies the bandwidth
	Bandwidth int `json:"bandwidth"`
}

// FlavorsPage is the page returned by a pager when traversing over a
// collection of flavor.
type FlavorsPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a FlavorsPage struct is empty.
func (r FlavorsPage) IsEmpty() (bool, error) {
	is, err := ExtractFlavors(r)
	return len(is) == 0, err
}

// ExtractFlavors accepts a Page struct, specifically a FlavorsPage struct,
// and extracts the elements into a slice of flavor structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractFlavors(r pagination.Page) ([]Flavor, error) {
	var s struct {
		Flavors []Flavor `json:"flavors"`
	}
	err := (r.(FlavorsPage)).ExtractInto(&s)
	return s.Flavors, err
}
