package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListStreamForbiddenResponse struct {

	// 查询结果的总元素数量
	Total *int32 `json:"total,omitempty"`

	// 禁播黑名单列表
	Blocks         *[]StreamForbiddenList `json:"blocks,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ListStreamForbiddenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListStreamForbiddenResponse struct{}"
	}

	return strings.Join([]string{"ListStreamForbiddenResponse", string(data)}, " ")
}
