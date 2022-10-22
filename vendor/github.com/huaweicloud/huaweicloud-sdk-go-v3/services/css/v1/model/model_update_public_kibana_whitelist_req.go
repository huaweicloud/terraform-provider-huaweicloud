package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdatePublicKibanaWhitelistReq struct {

	// 允许kibana公网访问的白名单ip或网段，以逗号隔开，不可重复。
	WhiteList string `json:"whiteList"`
}

func (o UpdatePublicKibanaWhitelistReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePublicKibanaWhitelistReq struct{}"
	}

	return strings.Join([]string{"UpdatePublicKibanaWhitelistReq", string(data)}, " ")
}
