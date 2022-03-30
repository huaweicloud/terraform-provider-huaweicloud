package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteGaussMySqlInstanceResponse struct {
	// 任务ID。

	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteGaussMySqlInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteGaussMySqlInstanceResponse struct{}"
	}

	return strings.Join([]string{"DeleteGaussMySqlInstanceResponse", string(data)}, " ")
}
