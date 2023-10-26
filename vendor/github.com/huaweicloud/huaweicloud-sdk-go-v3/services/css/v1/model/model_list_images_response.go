package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListImagesResponse Response Object
type ListImagesResponse struct {

	// 是否需要上传升级后版本的插件。
	NeedUploadUpgradePlugin *bool `json:"needUploadUpgradePlugin,omitempty"`

	ImageInfoList  *[]GetTargetImageIdDetail `json:"imageInfoList,omitempty"`
	HttpStatusCode int                       `json:"-"`
}

func (o ListImagesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListImagesResponse struct{}"
	}

	return strings.Join([]string{"ListImagesResponse", string(data)}, " ")
}
