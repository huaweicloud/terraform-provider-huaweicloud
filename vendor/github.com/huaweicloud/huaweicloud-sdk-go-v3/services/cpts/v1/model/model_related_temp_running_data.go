package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RelatedTempRunningData struct {
	// task_run_info_id

	TaskRunInfoId *int32 `json:"task_run_info_id,omitempty"`
	// related_temp_running_id

	RelatedTempRunningId *int32 `json:"related_temp_running_id,omitempty"`
	// temp_id

	TempId *int32 `json:"temp_id,omitempty"`
	// temp_name

	TempName *string `json:"temp_name,omitempty"`
	// related_temp_running_data

	RelatedTempRunningData *[]TempRunningData `json:"related_temp_running_data,omitempty"`
}

func (o RelatedTempRunningData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RelatedTempRunningData struct{}"
	}

	return strings.Join([]string{"RelatedTempRunningData", string(data)}, " ")
}
