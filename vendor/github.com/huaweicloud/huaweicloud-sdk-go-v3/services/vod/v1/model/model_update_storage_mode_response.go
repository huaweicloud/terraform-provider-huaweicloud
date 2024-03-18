package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateStorageModeResponse Response Object
type UpdateStorageModeResponse struct {
	TaskResultArray *[]TaskResult `json:"task_result_array,omitempty"`
	HttpStatusCode  int           `json:"-"`
}

func (o UpdateStorageModeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateStorageModeResponse struct{}"
	}

	return strings.Join([]string{"UpdateStorageModeResponse", string(data)}, " ")
}
