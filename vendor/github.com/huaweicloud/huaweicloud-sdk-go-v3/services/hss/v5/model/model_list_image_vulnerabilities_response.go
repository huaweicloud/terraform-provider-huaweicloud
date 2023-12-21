package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListImageVulnerabilitiesResponse Response Object
type ListImageVulnerabilitiesResponse struct {

	// 总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 镜像的漏洞列表
	DataList       *[]ImageVulInfo `json:"data_list,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ListImageVulnerabilitiesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListImageVulnerabilitiesResponse struct{}"
	}

	return strings.Join([]string{"ListImageVulnerabilitiesResponse", string(data)}, " ")
}
