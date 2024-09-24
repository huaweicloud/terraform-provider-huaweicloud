package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListVulScanTaskHostRequest Request Object
type ListVulScanTaskHostRequest struct {

	// 任务ID
	TaskId string `json:"task_id"`

	// 企业租户ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 主机的扫描状态，包含如下：   -scanning : 扫描中   -success : 扫描成功   -failed : 扫描失败
	ScanStatus *string `json:"scan_status,omitempty"`
}

func (o ListVulScanTaskHostRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVulScanTaskHostRequest struct{}"
	}

	return strings.Join([]string{"ListVulScanTaskHostRequest", string(data)}, " ")
}
