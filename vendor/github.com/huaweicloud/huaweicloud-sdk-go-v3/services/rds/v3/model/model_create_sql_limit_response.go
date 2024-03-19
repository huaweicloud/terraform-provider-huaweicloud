package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSqlLimitResponse Response Object
type CreateSqlLimitResponse struct {

	// 调用正常时，返回“successful”。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateSqlLimitResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSqlLimitResponse struct{}"
	}

	return strings.Join([]string{"CreateSqlLimitResponse", string(data)}, " ")
}
