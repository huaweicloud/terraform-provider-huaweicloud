package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ChangeHostsGroupRequestInfo struct {

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 服务器组ID
	GroupId string `json:"group_id"`

	// 服务器ID列表
	HostIdList *[]string `json:"host_id_list,omitempty"`
}

func (o ChangeHostsGroupRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeHostsGroupRequestInfo struct{}"
	}

	return strings.Join([]string{"ChangeHostsGroupRequestInfo", string(data)}, " ")
}
