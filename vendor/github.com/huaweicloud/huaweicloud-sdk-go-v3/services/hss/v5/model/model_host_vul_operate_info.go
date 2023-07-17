package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HostVulOperateInfo 主机维度漏洞列表
type HostVulOperateInfo struct {

	// 主机ID
	HostId string `json:"host_id"`

	// 漏洞列表
	VulIdList []string `json:"vul_id_list"`
}

func (o HostVulOperateInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostVulOperateInfo struct{}"
	}

	return strings.Join([]string{"HostVulOperateInfo", string(data)}, " ")
}
