package application

import (
	"strconv"

	"github.com/chnsz/golangsdk/pagination"
)

// Template is the structure that represents the aplication template details.
type Template struct {
	// Template ID.
	ID string `json:"id"`
	// The name of template.
	Name string `json:"name"`
	// Template execution runtime.
	Runtime string `json:"runtime"`
	// The category of template.
	Category string `json:"category"`
	// Template image file(Base64-encoded).
	Image string `json:"image"`
	// Template description.
	Description string `json:"description"`
	// The type of the function application.
	Type string `json:"type"`
}

// pageInfo is the structure that represents response of the ListTemplate method.
type pageInfo struct {
	// The list of templates.
	Templates []Template `json:"templates"`
	// Next record location.
	NextMarker int `json:"next_marker"`
	// Total number of application template.
	Count int `json:"count"`
}

// TemplatePage represents the response pages of the ListTemplate method.
type TemplatePage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult no application template.
func (r TemplatePage) IsEmpty() (bool, error) {
	resp, err := extractPageInfo(r)
	return resp.Count == 0, err
}

// LastMarker returns the last marker index in a pageInfo.
func (r TemplatePage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", nil
	}
	return strconv.Itoa(resp.NextMarker), nil
}

// NextPageURL generates the URL for the page of results after this one.
func (r TemplatePage) NextPageURL() (string, error) {
	currentURL := r.URL
	mark, err := r.Owner.LastMarker()
	if err != nil {
		return "", err
	}
	if mark == "" {
		return "", nil
	}

	q := currentURL.Query()
	q.Set("marker", mark)
	currentURL.RawQuery = q.Encode()

	return currentURL.String(), nil
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*pageInfo, error) {
	var s pageInfo
	err := r.(TemplatePage).Result.ExtractInto(&s)
	return &s, err
}
