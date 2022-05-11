package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type AssociateKeypairResponse struct {

	// 任务下发成功返回的ID
	TaskId         *string `json:"task_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AssociateKeypairResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateKeypairResponse struct{}"
	}

	return strings.Join([]string{"AssociateKeypairResponse", string(data)}, " ")
}
