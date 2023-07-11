package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListJarPackageStatisticsResponse Response Object
type ListJarPackageStatisticsResponse struct {

	// Jar包统计信息总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// Jar包统计信息列表
	DataList       *[]JarPackageStatisticsResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                                 `json:"-"`
}

func (o ListJarPackageStatisticsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListJarPackageStatisticsResponse struct{}"
	}

	return strings.Join([]string{"ListJarPackageStatisticsResponse", string(data)}, " ")
}
