package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PreheatingTaskRequestBody struct {

	// 输入URL必须带有“http://”或“https://”，多个URL用逗号分隔，目前不支持对目录的预热，单个url的长度限制为4096字符,单次最多输入1000个url。
	Urls []string `json:"urls"`
}

func (o PreheatingTaskRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PreheatingTaskRequestBody struct{}"
	}

	return strings.Join([]string{"PreheatingTaskRequestBody", string(data)}, " ")
}
