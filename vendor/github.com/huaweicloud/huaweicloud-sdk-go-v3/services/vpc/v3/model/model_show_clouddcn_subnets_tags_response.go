package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClouddcnSubnetsTagsResponse Response Object
type ShowClouddcnSubnetsTagsResponse struct {

	// 本次请求的编号
	RequestId *string `json:"request_id,omitempty"`

	// 单个资源的租户标签列表。
	Tags *[]TagEntity `json:"tags,omitempty"`

	// 单个资源的系统标签列表。
	SysTags        *[]SysTag `json:"sys_tags,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ShowClouddcnSubnetsTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClouddcnSubnetsTagsResponse struct{}"
	}

	return strings.Join([]string{"ShowClouddcnSubnetsTagsResponse", string(data)}, " ")
}
