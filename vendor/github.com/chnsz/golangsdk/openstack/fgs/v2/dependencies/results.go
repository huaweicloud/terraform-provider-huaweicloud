package dependencies

import (
	"strconv"

	"github.com/chnsz/golangsdk/pagination"
)

// Dependency is an object struct that represents the elements of the dependencies parameter.
type Dependency struct {
	// Dependency ID.
	ID string `json:"id"`
	// Dependency owner.
	Owner string `json:"owner"`
	// URL of the dependency in the OBS console.
	Link string `json:"link"`
	// Runtime.
	Runtime string `json:"runtime"`
	// Unique ID of the dependency.
	Etag string `json:"etag"`
	// Size of the dependency.
	Size int `json:"size"`
	// Name of the dependency.
	Name string `json:"name"`
	// Description of the dependency.
	Description string `json:"description"`
	// File name of the dependency.
	FileName string `json:"file_name"`
	// Dependency package version number.
	Version int `json:"version"`
}

// ListResp is an object struct that represents the result of each page.
type ListResp struct {
	// Next read location.
	Next int `json:"next_marker"`
	// Total number of dependencies.
	Count int `json:"count"`
	// Dependency list.
	Dependencies []Dependency `json:"dependencies"`
}

// DependencyPage represents the response pages of the List method.
type DependencyPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult no dependent package.
func (r DependencyPage) IsEmpty() (bool, error) {
	resp, err := ExtractDependencies(r)
	return len(resp.Dependencies) == 0, err
}

// LastMarker returns the last marker index in a ListResult.
func (r DependencyPage) LastMarker() (string, error) {
	resp, err := ExtractDependencies(r)
	if err != nil {
		return "", err
	}
	if resp.Next >= resp.Count {
		return "", nil
	}
	return strconv.Itoa(resp.Next), nil
}

// NextPageURL generates the URL for the page of results after this one.
func (r DependencyPage) NextPageURL() (string, error) {
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

// ExtractDependencies is a method which to extract the response to a dependent package list, next marker index and
// count number.
func ExtractDependencies(r pagination.Page) (ListResp, error) {
	var s ListResp
	err := r.(DependencyPage).Result.ExtractInto(&s)
	return s, err
}

// DependencyVersion is an object struct that represents the dependency (with version info) detail.
type DependencyVersion struct {
	// Dependency ID.
	ID string `json:"id"`
	// Dependency owner.
	Owner string `json:"owner"`
	// URL of the dependency in the OBS console.
	Link string `json:"link"`
	// Runtime.
	Runtime string `json:"runtime"`
	// Unique ID of the dependency.
	Etag string `json:"etag"`
	// Size of the dependency.
	Size int `json:"size"`
	// Name of the dependency.
	Name string `json:"name"`
	// Description of the dependency.
	Description string `json:"description"`
	// File name of the dependency.
	FileName string `json:"file_name"`
	// Version of the dependency.
	Version int `json:"version"`
	// The ID of the dependency.
	DepId string `json:"dep_id"`
	// The last modified time.
	LastModified int `json:"last_modified"`
}

type DependencyVersionPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult no dependent package version.
func (r DependencyVersionPage) IsEmpty() (bool, error) {
	versions, err := extractDependencieVersions(r)
	return len(versions) == 0, err
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (ListResp, error) {
	var s ListResp
	err := r.(DependencyVersionPage).Result.ExtractInto(&s)
	return s, err
}

// LastMarker returns the last marker index in a ListResult.
func (r DependencyVersionPage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(resp.Next), nil
}

func extractDependencieVersions(r pagination.Page) ([]DependencyVersion, error) {
	var s []DependencyVersion
	err := r.(DependencyVersionPage).Result.ExtractIntoSlicePtr(&s, "dependencies")
	return s, err
}
