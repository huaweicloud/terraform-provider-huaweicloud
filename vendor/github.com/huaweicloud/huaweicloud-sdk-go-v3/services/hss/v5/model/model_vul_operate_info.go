package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 漏洞列表
type VulOperateInfo struct {

	// 漏洞ID
	VulId *string `json:"vul_id,omitempty"`

	// 主机列表
	HostIdList *[]string `json:"host_id_list,omitempty"`
}

func (o VulOperateInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulOperateInfo struct{}"
	}

	return strings.Join([]string{"VulOperateInfo", string(data)}, " ")
}
