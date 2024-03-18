package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSqlLimitResponse Response Object
type DeleteSqlLimitResponse struct {

	// 调用正常时，返回“successful”。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteSqlLimitResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSqlLimitResponse struct{}"
	}

	return strings.Join([]string{"DeleteSqlLimitResponse", string(data)}, " ")
}
