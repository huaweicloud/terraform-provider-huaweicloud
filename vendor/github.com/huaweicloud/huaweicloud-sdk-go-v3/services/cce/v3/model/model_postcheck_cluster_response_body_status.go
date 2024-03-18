package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostcheckClusterResponseBodyStatus 集群升级后确认的状态信息
type PostcheckClusterResponseBodyStatus struct {

	// 状态，取值如下 - Success 成功 - Failed 失败 - Error 错误
	Phase *string `json:"phase,omitempty"`
}

func (o PostcheckClusterResponseBodyStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostcheckClusterResponseBodyStatus struct{}"
	}

	return strings.Join([]string{"PostcheckClusterResponseBodyStatus", string(data)}, " ")
}
