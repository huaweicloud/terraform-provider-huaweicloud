package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTagKeysRequest struct {

	// 区域ID
	RegionId *string `json:"region_id,omitempty"`

	// 查询记录数。 最小为1，最大为200，未输入时默认为200。
	Limit *int32 `json:"limit,omitempty"`

	// 分页位置标识（索引）。 从marker指定索引的下一条数据开始查询。 说明： 查询第一页数据时，不需要传入此参数，查询后续页码数据时，将查询前一页数据响应体中marker值配入此参数，当返回的next_marker为空时表示查询到最后一页。
	Marker *string `json:"marker,omitempty"`
}

func (o ListTagKeysRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTagKeysRequest struct{}"
	}

	return strings.Join([]string{"ListTagKeysRequest", string(data)}, " ")
}
