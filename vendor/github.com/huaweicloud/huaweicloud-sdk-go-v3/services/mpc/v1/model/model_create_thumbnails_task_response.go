package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateThumbnailsTaskResponse struct {

	// 任务ID。
	TaskId *string `json:"task_id,omitempty"`

	// 任务状态
	Status *string `json:"status,omitempty"`

	// 任务创建时间
	CreateTime *string `json:"create_time,omitempty"`

	Output *ObsObjInfo `json:"output,omitempty"`

	// 截图文件名称
	OutputFileName *string `json:"output_file_name,omitempty"`

	// 指定的截图时间点
	ThumbnailTime *string `json:"thumbnail_time,omitempty"`

	// 截图任务描述，当截图出现异常时，此字段为异常的原因
	Description    *string `json:"description,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateThumbnailsTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateThumbnailsTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateThumbnailsTaskResponse", string(data)}, " ")
}
