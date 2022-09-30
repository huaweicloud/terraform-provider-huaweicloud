package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartPublicWhitelistReq struct {

	// 开启白名单的用户IP。
	WhiteList string `json:"whiteList"`
}

func (o StartPublicWhitelistReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartPublicWhitelistReq struct{}"
	}

	return strings.Join([]string{"StartPublicWhitelistReq", string(data)}, " ")
}
