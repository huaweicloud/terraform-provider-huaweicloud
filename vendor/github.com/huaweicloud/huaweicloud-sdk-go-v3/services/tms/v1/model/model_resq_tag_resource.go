package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResqTagResource 获取标签下资源请求
type ResqTagResource struct {

	// 项目ID，resource_type为region级别服务时为必选项。
	ProjectId *string `json:"project_id,omitempty"`

	// 资源类型， 此参数为可输入的值（区分大小写）。例如：ecs,scaling_group, images, disk,vpcs,security-groups, shared_bandwidth,eip, cdn等，具体请参见“附录-标签支持的资源类型”章节。
	ResourceTypes []string `json:"resource_types"`

	// 标签列表
	Tags []Tag `json:"tags"`

	// 是否仅查询未带标签的资源。该字段为true时查询不带标签的资源。
	WithoutAnyTag *bool `json:"without_any_tag,omitempty"`

	// 索引位置， 从offset指定的下一条数据开始查询，必须为数字，不能为负数，默认为0。
	Offset *int32 `json:"offset,omitempty"`

	// 查询记录数，不传默认为200，limit最多为200, 最小值为1。
	Limit *int32 `json:"limit,omitempty"`
}

func (o ResqTagResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResqTagResource struct{}"
	}

	return strings.Join([]string{"ResqTagResource", string(data)}, " ")
}
