package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListBatchTaskFilesResponse struct {

	// 批量任务文件列表。
	Files          *[]BatchTaskFile `json:"files,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListBatchTaskFilesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListBatchTaskFilesResponse struct{}"
	}

	return strings.Join([]string{"ListBatchTaskFilesResponse", string(data)}, " ")
}
