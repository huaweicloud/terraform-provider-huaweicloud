package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteRdSforMySqlProxyResponse Response Object
type DeleteRdSforMySqlProxyResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteRdSforMySqlProxyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteRdSforMySqlProxyResponse struct{}"
	}

	return strings.Join([]string{"DeleteRdSforMySqlProxyResponse", string(data)}, " ")
}
