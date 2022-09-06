package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CheckMd5DuplicationRequest struct {

	// 文件大小
	Size int64 `json:"size"`

	// 文件MD5。  建议参考[上传校验](https://support.huaweicloud.com/api-vod/vod_04_0212.html#vod_04_0212__section575102165412)生成对应的MD5值。
	Md5 string `json:"md5"`
}

func (o CheckMd5DuplicationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckMd5DuplicationRequest struct{}"
	}

	return strings.Join([]string{"CheckMd5DuplicationRequest", string(data)}, " ")
}
