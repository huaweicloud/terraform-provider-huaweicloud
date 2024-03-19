package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateRdSforMySqlProxyResponse Response Object
type CreateRdSforMySqlProxyResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateRdSforMySqlProxyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRdSforMySqlProxyResponse struct{}"
	}

	return strings.Join([]string{"CreateRdSforMySqlProxyResponse", string(data)}, " ")
}
