package catalog

import (
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/tokens"
	"github.com/huaweicloud/golangsdk/pagination"
)

// CatalogPage is a single page of Service results.
type CatalogPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if the CatalogPage contains no results.
func (p CatalogPage) IsEmpty() (bool, error) {
	services, err := ExtractServiceCatalog(p)
	return len(services) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r CatalogPage) NextPageURL() (string, error) {
	var s struct {
		Links struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Links.Next, err
}

// ExtractServiceCatalog extracts a slice of Catalog from a Collection acquired
// from List.
func ExtractServiceCatalog(r pagination.Page) ([]tokens.CatalogEntry, error) {
	var s struct {
		Entries []tokens.CatalogEntry `json:"catalog"`
	}
	err := (r.(CatalogPage)).ExtractInto(&s)
	return s.Entries, err
}
