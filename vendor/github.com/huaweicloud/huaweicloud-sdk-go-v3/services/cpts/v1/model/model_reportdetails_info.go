package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ReportdetailsInfo struct {
	// 表格数据

	Data *[]ReportdetailItemInfo `json:"data,omitempty"`
	// 页码

	PageIndex *int32 `json:"pageIndex,omitempty"`
	// 每页大小

	PageSize *int32 `json:"pageSize,omitempty"`
	// 总记录数

	Total *int32 `json:"total,omitempty"`
}

func (o ReportdetailsInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReportdetailsInfo struct{}"
	}

	return strings.Join([]string{"ReportdetailsInfo", string(data)}, " ")
}
