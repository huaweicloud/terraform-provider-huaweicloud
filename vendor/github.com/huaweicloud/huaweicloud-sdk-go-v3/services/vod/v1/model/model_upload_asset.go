package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UploadAsset struct {

	// 媒资所在url
	Url *string `json:"url,omitempty"`

	// 新创建媒资的媒资id
	AssetId *string `json:"asset_id,omitempty"`

	// 错误码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 错误描述。
	ErrorMsg *string `json:"error_msg,omitempty"`
}

func (o UploadAsset) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UploadAsset struct{}"
	}

	return strings.Join([]string{"UploadAsset", string(data)}, " ")
}
