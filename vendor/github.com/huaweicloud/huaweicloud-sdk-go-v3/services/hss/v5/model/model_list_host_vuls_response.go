package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHostVulsResponse Response Object
type ListHostVulsResponse struct {

	// 总数
	TotalNum *int64 `json:"total_num,omitempty"`

	// 服务器上的漏洞列表
	DataList       *[]HostVulInfo `json:"data_list,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListHostVulsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHostVulsResponse struct{}"
	}

	return strings.Join([]string{"ListHostVulsResponse", string(data)}, " ")
}
