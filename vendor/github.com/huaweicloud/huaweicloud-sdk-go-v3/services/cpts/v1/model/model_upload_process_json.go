package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// json
type UploadProcessJson struct {
	// details

	Details *[]UploadProcessJsonDetail `json:"details,omitempty"`
	// process_status

	ProcessStatus *int32 `json:"process_status,omitempty"`
}

func (o UploadProcessJson) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UploadProcessJson struct{}"
	}

	return strings.Join([]string{"UploadProcessJson", string(data)}, " ")
}
