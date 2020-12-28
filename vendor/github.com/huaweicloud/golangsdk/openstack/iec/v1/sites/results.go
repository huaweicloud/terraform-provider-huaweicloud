package sites

import (
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
	"github.com/huaweicloud/golangsdk/pagination"
)

// Sites 站点集合
type Sites struct {
	// 数量
	Count int32 `json:"count"`

	Sites []common.Site `json:"sites"`
}

// SitePage is the page returned by a pager when traversing over a collection
// of sites.
type SitePage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no site results.
func (r SitePage) IsEmpty() (bool, error) {
	s, err := ExtractSites(r)
	return s.Count == 0, err
}

// ExtractSites interprets the results of a single page from a List() call,
// producing a slice of site entities.
func ExtractSites(r pagination.Page) (Sites, error) {
	var s Sites
	err := r.(SitePage).ExtractInto(&s)
	return s, err
}
