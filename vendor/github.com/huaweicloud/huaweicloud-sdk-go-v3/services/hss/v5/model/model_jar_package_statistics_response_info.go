package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// JarPackageStatisticsResponseInfo Jar包统计信息列表
type JarPackageStatisticsResponseInfo struct {

	// Jar包名称
	FileName *string `json:"file_name,omitempty"`

	// Jar包统计信息总数
	Num *int32 `json:"num,omitempty"`
}

func (o JarPackageStatisticsResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "JarPackageStatisticsResponseInfo struct{}"
	}

	return strings.Join([]string{"JarPackageStatisticsResponseInfo", string(data)}, " ")
}
