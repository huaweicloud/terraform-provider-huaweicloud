package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowGetLogSettingRequest Request Object
type ShowGetLogSettingRequest struct {

	// 指定查询集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ShowGetLogSettingRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGetLogSettingRequest struct{}"
	}

	return strings.Join([]string{"ShowGetLogSettingRequest", string(data)}, " ")
}
