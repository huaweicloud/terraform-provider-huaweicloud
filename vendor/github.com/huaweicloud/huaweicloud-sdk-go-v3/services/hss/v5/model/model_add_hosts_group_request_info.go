package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddHostsGroupRequestInfo struct {

	// 服务器组名称
	GroupName string `json:"group_name"`

	// 服务器ID列表
	HostIdList []string `json:"host_id_list"`
}

func (o AddHostsGroupRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddHostsGroupRequestInfo struct{}"
	}

	return strings.Join([]string{"AddHostsGroupRequestInfo", string(data)}, " ")
}
