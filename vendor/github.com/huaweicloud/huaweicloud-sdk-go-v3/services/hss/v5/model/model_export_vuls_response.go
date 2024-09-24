package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExportVulsResponse Response Object
type ExportVulsResponse struct {

	// 任务ID
	TaskId         *string `json:"task_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ExportVulsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExportVulsResponse struct{}"
	}

	return strings.Join([]string{"ExportVulsResponse", string(data)}, " ")
}
