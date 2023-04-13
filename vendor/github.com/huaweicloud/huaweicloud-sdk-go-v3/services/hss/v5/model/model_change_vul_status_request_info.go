package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 漏洞操作列表
type ChangeVulStatusRequestInfo struct {

	// 操作类型 - ignore : 忽略 - not_ignore : 取消忽略 - immediate_repair : 修复 - verify : 验证
	OperateType string `json:"operate_type"`

	// 漏洞列表
	DataList []VulOperateInfo `json:"data_list"`
}

func (o ChangeVulStatusRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeVulStatusRequestInfo struct{}"
	}

	return strings.Join([]string{"ChangeVulStatusRequestInfo", string(data)}, " ")
}
