package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteGaussMySqlReadonlyNodeResponse struct {
	// 任务ID。

	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteGaussMySqlReadonlyNodeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteGaussMySqlReadonlyNodeResponse struct{}"
	}

	return strings.Join([]string{"DeleteGaussMySqlReadonlyNodeResponse", string(data)}, " ")
}
