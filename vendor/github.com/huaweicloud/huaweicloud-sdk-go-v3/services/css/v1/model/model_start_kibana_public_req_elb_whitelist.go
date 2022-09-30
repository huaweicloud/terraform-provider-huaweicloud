package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartKibanaPublicReqElbWhitelist struct {

	// 是否开启白名单。 - true: 开启白名单。 - false: 关闭白名单。
	EnableWhiteList bool `json:"enableWhiteList"`

	// 白名单。
	WhiteList string `json:"whiteList"`
}

func (o StartKibanaPublicReqElbWhitelist) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartKibanaPublicReqElbWhitelist struct{}"
	}

	return strings.Join([]string{"StartKibanaPublicReqElbWhitelist", string(data)}, " ")
}
