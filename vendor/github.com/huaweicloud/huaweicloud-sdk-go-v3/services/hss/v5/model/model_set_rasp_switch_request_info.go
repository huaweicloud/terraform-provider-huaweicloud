package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetRaspSwitchRequestInfo struct {

	// HostId list
	HostIdList *[]string `json:"host_id_list,omitempty"`

	// 动态网页防篡改状态
	Status *bool `json:"status,omitempty"`
}

func (o SetRaspSwitchRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetRaspSwitchRequestInfo struct{}"
	}

	return strings.Join([]string{"SetRaspSwitchRequestInfo", string(data)}, " ")
}
