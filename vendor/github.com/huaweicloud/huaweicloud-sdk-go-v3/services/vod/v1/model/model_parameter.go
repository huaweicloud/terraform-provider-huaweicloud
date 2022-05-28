package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 需要提取音频的参数信息。
type Parameter struct {

	// 封装格式。  取值如下： - MP3 - AAC
	Format *string `json:"format,omitempty"`
}

func (o Parameter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Parameter struct{}"
	}

	return strings.Join([]string{"Parameter", string(data)}, " ")
}
