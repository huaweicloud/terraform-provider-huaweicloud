package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowRegionInfoResp 查询Region信息响应体
type ShowRegionInfoResp struct {

	// 服务名称
	ServiceName *string `json:"service_name,omitempty"`

	// Region列表
	RegionList *[]RegionInfo `json:"region_list,omitempty"`
}

func (o ShowRegionInfoResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRegionInfoResp struct{}"
	}

	return strings.Join([]string{"ShowRegionInfoResp", string(data)}, " ")
}
