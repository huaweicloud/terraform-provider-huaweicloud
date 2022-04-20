package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TempRunningData struct {
	// content_method_url

	ContentMethodUrl *[]string `json:"content_method_url,omitempty"`
	// crawler_status

	CrawlerStatus *int32 `json:"crawler_status,omitempty"`
	// related_temp_running_id

	RelatedTempRunningId *int32 `json:"related_temp_running_id,omitempty"`
	// task_run_info_id

	TaskRunInfoId *int32 `json:"task_run_info_id,omitempty"`
	// temp_id

	TempId *int32 `json:"temp_id,omitempty"`
	// temp_name

	TempName *string `json:"temp_name,omitempty"`
	// temp_running_status

	TempRunningStatus *int32 `json:"temp_running_status,omitempty"`
}

func (o TempRunningData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TempRunningData struct{}"
	}

	return strings.Join([]string{"TempRunningData", string(data)}, " ")
}
