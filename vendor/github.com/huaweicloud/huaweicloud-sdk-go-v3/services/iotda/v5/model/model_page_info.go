package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PageInfo 查询批量分页结构体，定义了分页页码、每页记录数、记录总数、该页记录的最大Id。
type PageInfo struct {

	// 满足查询条件的记录总数。
	Count *int64 `json:"count,omitempty"`

	// 本次分页查询结果中最后一条记录的ID，可在下一次分页查询时使用。
	Marker *string `json:"marker,omitempty"`
}

func (o PageInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PageInfo struct{}"
	}

	return strings.Join([]string{"PageInfo", string(data)}, " ")
}
