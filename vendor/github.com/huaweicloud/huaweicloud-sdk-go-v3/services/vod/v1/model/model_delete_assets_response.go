package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteAssetsResponse struct {

	// 删除媒资任务的处理结果。
	DeleteResultArray *[]DeleteResult `json:"delete_result_array,omitempty"`
	HttpStatusCode    int             `json:"-"`
}

func (o DeleteAssetsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAssetsResponse struct{}"
	}

	return strings.Join([]string{"DeleteAssetsResponse", string(data)}, " ")
}
