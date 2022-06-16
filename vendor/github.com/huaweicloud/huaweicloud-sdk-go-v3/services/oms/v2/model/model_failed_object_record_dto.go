package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type FailedObjectRecordDto struct {

	// 是否支持失败对象重传。
	Result *bool `json:"result,omitempty"`

	// 失败对象列表文件路径。
	ListFileKey *string `json:"list_file_key,omitempty"`

	// 失败对象列表上传失败的错误码。
	ErrorCode *string `json:"error_code,omitempty"`
}

func (o FailedObjectRecordDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FailedObjectRecordDto struct{}"
	}

	return strings.Join([]string{"FailedObjectRecordDto", string(data)}, " ")
}
