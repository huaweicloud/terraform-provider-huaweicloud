package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTrafficMirrorFiltersRequest Request Object
type ListTrafficMirrorFiltersRequest struct {

	// 使用ID过滤查询或排序
	Id *string `json:"id,omitempty"`

	// 使用name过滤或排序
	Name *string `json:"name,omitempty"`

	// 使用description过滤查询
	Description *string `json:"description,omitempty"`

	// 使用创建时间戳排序
	CreatedAt *string `json:"created_at,omitempty"`

	// 使用更新时间戳排序
	UpdatedAt *string `json:"updated_at,omitempty"`
}

func (o ListTrafficMirrorFiltersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTrafficMirrorFiltersRequest struct{}"
	}

	return strings.Join([]string{"ListTrafficMirrorFiltersRequest", string(data)}, " ")
}
