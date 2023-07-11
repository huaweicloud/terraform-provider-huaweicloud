package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartAutoSettingRequest Request Object
type StartAutoSettingRequest struct {

	// 指定要备份的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o StartAutoSettingRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartAutoSettingRequest struct{}"
	}

	return strings.Join([]string{"StartAutoSettingRequest", string(data)}, " ")
}
