package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListVulScanTaskResponse Response Object
type ListVulScanTaskResponse struct {

	// 总数
	TotalNum *int64 `json:"total_num,omitempty"`

	// 漏洞扫描任务列表
	DataList       *[]VulScanTaskInfo `json:"data_list,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ListVulScanTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVulScanTaskResponse struct{}"
	}

	return strings.Join([]string{"ListVulScanTaskResponse", string(data)}, " ")
}
