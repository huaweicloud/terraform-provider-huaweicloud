package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DisassociateKeypairResponse struct {

	// 任务下发成功返回的ID
	TaskId         *string `json:"task_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DisassociateKeypairResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DisassociateKeypairResponse struct{}"
	}

	return strings.Join([]string{"DisassociateKeypairResponse", string(data)}, " ")
}
