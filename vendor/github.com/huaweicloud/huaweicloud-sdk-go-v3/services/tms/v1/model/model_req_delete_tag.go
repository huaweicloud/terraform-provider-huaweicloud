package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ReqDeleteTag 删除标签请求
type ReqDeleteTag struct {

	// 项目ID，resource_type为region级别服务时为必选项。
	ProjectId *string `json:"project_id,omitempty"`

	// 资源列表
	Resources []ResourceTagBody `json:"resources"`

	// 标签列表
	Tags []DeleteTagRequest `json:"tags"`
}

func (o ReqDeleteTag) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReqDeleteTag struct{}"
	}

	return strings.Join([]string{"ReqDeleteTag", string(data)}, " ")
}
