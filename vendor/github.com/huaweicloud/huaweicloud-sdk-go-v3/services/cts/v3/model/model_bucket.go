package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Bucket struct {

	// 标识OBS桶名称。由数字或字母开头，支持小写字母、数字、“-”、“.”，长度为3～63个字符。
	BucketName string `json:"bucket_name"`

	// 标识桶位置。
	BucketLocation *string `json:"bucket_location,omitempty"`

	// 事件文件转储加密所采用的秘钥id。 如果is_support_trace_files_encryption\"参数值为“是”时，此参数为必选项。
	KmsId *string `json:"kms_id,omitempty"`

	// 事件文件转储加密功能开关。 该参数必须与kms_id参数同时使用。
	IsSupportTraceFilesEncryption bool `json:"is_support_trace_files_encryption"`

	CheckBucketResponse *CheckBucketResponse `json:"check_bucket_response,omitempty"`
}

func (o Bucket) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Bucket struct{}"
	}

	return strings.Join([]string{"Bucket", string(data)}, " ")
}
