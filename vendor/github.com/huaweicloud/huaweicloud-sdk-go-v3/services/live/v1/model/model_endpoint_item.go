package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EndpointItem 拉流打包信息
type EndpointItem struct {

	// HLS打包信息
	HlsPackage *[]HlsPackageItem `json:"hls_package,omitempty"`

	// DASH打包信息
	DashPackage *[]DashPackageItem `json:"dash_package,omitempty"`

	// MSS打包信息
	MssPackage *[]MssPackageItem `json:"mss_package,omitempty"`
}

func (o EndpointItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EndpointItem struct{}"
	}

	return strings.Join([]string{"EndpointItem", string(data)}, " ")
}
