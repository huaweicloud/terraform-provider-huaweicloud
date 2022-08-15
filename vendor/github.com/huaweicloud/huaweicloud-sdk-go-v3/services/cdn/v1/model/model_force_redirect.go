package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ForceRedirect struct {

	// 强制跳转开关。1打开。0关闭。
	Switch int32 `json:"switch"`

	// 强制跳转类型。http：强制跳转HTTP。https：强制跳转HTTPS。
	RedirectType *string `json:"redirect_type,omitempty"`
}

func (o ForceRedirect) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ForceRedirect struct{}"
	}

	return strings.Join([]string{"ForceRedirect", string(data)}, " ")
}
