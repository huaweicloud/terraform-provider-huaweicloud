package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SwitchSqlLimitResponse Response Object
type SwitchSqlLimitResponse struct {

	// 调用正常时，返回“successful”。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SwitchSqlLimitResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchSqlLimitResponse struct{}"
	}

	return strings.Join([]string{"SwitchSqlLimitResponse", string(data)}, " ")
}
