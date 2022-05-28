package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 封面信息。
type CoverInfo struct {

	// 封面文件的下载地址。
	CoverUrl *string `json:"cover_url,omitempty"`
}

func (o CoverInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CoverInfo struct{}"
	}

	return strings.Join([]string{"CoverInfo", string(data)}, " ")
}
