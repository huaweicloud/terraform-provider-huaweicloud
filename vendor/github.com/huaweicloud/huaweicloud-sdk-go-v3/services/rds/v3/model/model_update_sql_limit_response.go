package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateSqlLimitResponse Response Object
type UpdateSqlLimitResponse struct {

	// 调用正常时，返回“successful”。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateSqlLimitResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSqlLimitResponse struct{}"
	}

	return strings.Join([]string{"UpdateSqlLimitResponse", string(data)}, " ")
}
