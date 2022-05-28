package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateTakeOverTaskReq struct {

	// 源桶名。
	Bucket string `json:"bucket"`

	// 源目录名或源文件名。
	Object string `json:"object"`

	// 批量托管时的文件后缀名列表。不传或传空值时，表示托管所有音视频文件，不进行后缀名过滤。
	Suffix *[]string `json:"suffix,omitempty"`

	// 转码模板组名称。  若不为空，则使用指定的转码模板对上传的音视频进行转码，您可以在视频点播控制台配置转码模板，具体请参见转码设置。  > 若同时设置了“**template_group_name**”和“**workflow_name**”字段，则“**template_group_name**”字段生效。
	TemplateGroupName *string `json:"template_group_name,omitempty"`

	// 工作流名称。  若不为空，则使用指定的工作流对上传的音视频进行处理，您可以在视频点播控制台配置工作流，具体请参见[工作流设置](https://support.huaweicloud.com/usermanual-vod/vod010041.html)。
	WorkflowName *string `json:"workflow_name,omitempty"`

	// 表示音视频处理后生成的媒资文件所存储的位置类型。  取值如下所示： - 0：表示存储到点播桶。 - 1：表示存储在租户桶。 - 2：表示存储到租户桶，并且存储路径与源文件一致。
	HostType *int32 `json:"host_type,omitempty"`

	// 输出桶名，host_type为1时必选
	OutputBucket *string `json:"output_bucket,omitempty"`

	// 输出路径名，host_type为1时必选
	OutputPath *string `json:"output_path,omitempty"`
}

func (o CreateTakeOverTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTakeOverTaskReq struct{}"
	}

	return strings.Join([]string{"CreateTakeOverTaskReq", string(data)}, " ")
}
