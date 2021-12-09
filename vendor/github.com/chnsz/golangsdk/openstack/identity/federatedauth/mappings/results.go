package mappings

import (
	"github.com/chnsz/golangsdk/pagination"
)

type IdentityMapping struct {
	ID    string        `json:"id"`
	Rules []MappingRule `json:"rules"`
	Links LinksSelf     `json:"links"`
}

type LinksSelf struct {
	Self string `json:"self"`
}

type Links struct {
	Next     string `json:"next"`
	Self     string `json:"self"`
	Previous string `json:"previous"`
}

type IdentityMappingPage struct {
	pagination.LinkedPageBase
}

func (r IdentityMappingPage) IsEmpty() (bool, error) {
	mappings, err := ExtractMappings(r)
	return len(mappings) == 0, err
}

func (r IdentityMappingPage) NextPageURL() (string, error) {
	var s struct {
		Links Links `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	url := s.Links.Next
	if url == "" {
		return "", nil
	}
	return url, nil
}

func ExtractMappings(r pagination.Page) ([]IdentityMapping, error) {
	var s struct {
		Links    Links             `json:"links"`
		Mappings []IdentityMapping `json:"mappings"`
	}
	err := (r.(IdentityMappingPage)).ExtractInto(&s)
	return s.Mappings, err
}
