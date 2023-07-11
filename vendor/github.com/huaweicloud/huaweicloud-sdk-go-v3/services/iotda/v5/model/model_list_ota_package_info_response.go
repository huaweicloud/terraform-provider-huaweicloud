package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOtaPackageInfoResponse Response Object
type ListOtaPackageInfoResponse struct {

	// 升级包列表
	Packages *[]OtaPackageInfo `json:"packages,omitempty"`

	Page           *PageInfo `json:"page,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ListOtaPackageInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOtaPackageInfoResponse struct{}"
	}

	return strings.Join([]string{"ListOtaPackageInfoResponse", string(data)}, " ")
}
