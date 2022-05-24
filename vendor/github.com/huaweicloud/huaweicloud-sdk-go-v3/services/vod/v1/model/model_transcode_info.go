package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 转码生成文件信息。  > 仅当转码成功后才能查询到此信息，未转码、正在转码以及转码失败时，无此字段信息。
type TranscodeInfo struct {

	// 转码模板组名称。
	TemplateGroupName string `json:"template_group_name"`

	// 转码输出数组。 - HLS或DASH格式：此数组的成员个数为n+1，n为转码输出路数。 - MP4格式：此数组的成员个数为n，n为转码输出路数。
	Output []Output `json:"output"`

	// 执行情况描述。
	ExecDesc *string `json:"exec_desc,omitempty"`

	// 转码状态。  取值如下： - UN_TRANSCODE：未转码 - WAITING_TRANSCODE：待转码 - TRANSCODING：转码中 - TRANSCODE_SUCCEED：转码成功 - TRANSCODE_FAILED：转码失败
	TranscodeStatus *string `json:"transcode_status,omitempty"`
}

func (o TranscodeInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TranscodeInfo struct{}"
	}

	return strings.Join([]string{"TranscodeInfo", string(data)}, " ")
}
