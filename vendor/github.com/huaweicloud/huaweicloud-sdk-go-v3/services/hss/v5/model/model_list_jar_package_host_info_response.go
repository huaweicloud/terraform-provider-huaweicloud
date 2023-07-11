package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListJarPackageHostInfoResponse Response Object
type ListJarPackageHostInfoResponse struct {

	// 总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 服务器列表
	DataList       *[]JarPackageHostInfo `json:"data_list,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o ListJarPackageHostInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListJarPackageHostInfoResponse struct{}"
	}

	return strings.Join([]string{"ListJarPackageHostInfoResponse", string(data)}, " ")
}
