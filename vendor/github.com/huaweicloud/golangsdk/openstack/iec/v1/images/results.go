package images

import (
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
	"github.com/huaweicloud/golangsdk/pagination"
)

// EdgeImages 边缘镜像列表信息
type EdgeImages struct {
	Total  int                    `json:"total"`
	Images []common.EdgeImageInfo `json:"images"`
}

// ImagePage represents the results of a List request.
type ImagePage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no image results.
func (r ImagePage) IsEmpty() (bool, error) {
	s, err := ExtractImages(r)
	return s.Total == 0, err
}

// ExtractImages 输出边缘镜像列表
func ExtractImages(r pagination.Page) (EdgeImages, error) {
	var s EdgeImages
	err := (r.(ImagePage)).ExtractInto(&s)
	return s, err
}
