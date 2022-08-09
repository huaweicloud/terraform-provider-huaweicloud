package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRecordContentsResponse struct {

	// 查询结果的总元素数量
	Total *int32 `json:"total,omitempty"`

	// 录制内容数组
	RecordContents *[]RecordContentInfoV2 `json:"record_contents,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListRecordContentsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRecordContentsResponse struct{}"
	}

	return strings.Join([]string{"ListRecordContentsResponse", string(data)}, " ")
}
