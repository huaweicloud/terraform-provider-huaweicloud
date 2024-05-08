package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListVulnerabilitiesResponse Response Object
type ListVulnerabilitiesResponse struct {

	// 漏洞总数
	TotalNum *int64 `json:"total_num,omitempty"`

	// 软件漏洞列表
	DataList       *[]VulInfo `json:"data_list,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o ListVulnerabilitiesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVulnerabilitiesResponse struct{}"
	}

	return strings.Join([]string{"ListVulnerabilitiesResponse", string(data)}, " ")
}
