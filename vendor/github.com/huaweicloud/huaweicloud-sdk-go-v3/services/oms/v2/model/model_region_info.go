package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RegionInfo Region信息
type RegionInfo struct {

	// 云服务名称
	CloudType *string `json:"cloud_type,omitempty"`

	// region名称
	Value *string `json:"value,omitempty"`

	// region的描述信息
	Description *string `json:"description,omitempty"`
}

func (o RegionInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RegionInfo struct{}"
	}

	return strings.Join([]string{"RegionInfo", string(data)}, " ")
}
