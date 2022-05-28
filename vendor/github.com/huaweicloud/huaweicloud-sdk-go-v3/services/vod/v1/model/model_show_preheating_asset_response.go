package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowPreheatingAssetResponse struct {

	// 预热任务数组
	PreheatingResults *[]PreheatingResult `json:"preheating_results,omitempty"`
	HttpStatusCode    int                 `json:"-"`
}

func (o ShowPreheatingAssetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPreheatingAssetResponse struct{}"
	}

	return strings.Join([]string{"ShowPreheatingAssetResponse", string(data)}, " ")
}
