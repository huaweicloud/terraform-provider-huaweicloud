package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListVulScanTaskHostResponse Response Object
type ListVulScanTaskHostResponse struct {

	// 总数
	TotalNum *int64 `json:"total_num,omitempty"`

	// 漏洞扫描任务对应的主机列表
	DataList       *[]VulScanTaskHostInfo `json:"data_list,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ListVulScanTaskHostResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVulScanTaskHostResponse struct{}"
	}

	return strings.Join([]string{"ListVulScanTaskHostResponse", string(data)}, " ")
}
