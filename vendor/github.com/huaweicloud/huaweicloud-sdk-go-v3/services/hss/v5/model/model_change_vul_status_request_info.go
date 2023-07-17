package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeVulStatusRequestInfo 漏洞操作列表
type ChangeVulStatusRequestInfo struct {

	// 操作类型 - ignore : 忽略 - not_ignore : 取消忽略 - immediate_repair : 修复 - manual_repair: 人工修复 - verify : 验证 - add_to_whitelist : 加入白名单
	OperateType string `json:"operate_type"`

	// 备注
	Remark *string `json:"remark,omitempty"`

	// 选择全部漏洞类型 - all_vul : 选择全部漏洞 - all_host : 选择全部主机漏洞
	SelectType *string `json:"select_type,omitempty"`

	// 漏洞类型，默认为linux_vul，包括如下：   - linux_vul : 漏洞类型-linux漏洞   - windows_vul : 漏洞类型-windows漏洞   - web_cms : Web-CMS漏洞   - app_vul : 应用漏洞
	Type *string `json:"type,omitempty"`

	// 漏洞列表
	DataList []VulOperateInfo `json:"data_list"`

	// 主机维度漏洞列表
	HostDataList *[]HostVulOperateInfo `json:"host_data_list,omitempty"`
}

func (o ChangeVulStatusRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeVulStatusRequestInfo struct{}"
	}

	return strings.Join([]string{"ChangeVulStatusRequestInfo", string(data)}, " ")
}
