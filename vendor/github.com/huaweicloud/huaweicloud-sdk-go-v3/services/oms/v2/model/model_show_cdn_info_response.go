package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCdnInfoResponse Response Object
type ShowCdnInfoResponse struct {

	// 是否相同云类型
	IsSameCloudType *bool `json:"is_same_cloud_type,omitempty"`

	// 是否下载可用
	IsDownloadAvailable *bool `json:"is_download_available,omitempty"`

	// 返回的已检查的对象数组
	CheckedKeys    *[]CheckedKey `json:"checked_keys,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o ShowCdnInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCdnInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowCdnInfoResponse", string(data)}, " ")
}
