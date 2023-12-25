package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeBlockedIpRequestInfo 解除拦截的IP详情
type ChangeBlockedIpRequestInfo struct {

	// 需要解除拦截的IP列表
	DataList *[]BlockedIpRequestInfo `json:"data_list,omitempty"`
}

func (o ChangeBlockedIpRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeBlockedIpRequestInfo struct{}"
	}

	return strings.Join([]string{"ChangeBlockedIpRequestInfo", string(data)}, " ")
}
