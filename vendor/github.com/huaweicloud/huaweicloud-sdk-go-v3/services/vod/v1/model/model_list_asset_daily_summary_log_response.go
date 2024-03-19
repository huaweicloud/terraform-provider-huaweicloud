package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAssetDailySummaryLogResponse Response Object
type ListAssetDailySummaryLogResponse struct {

	// 记录总数
	Total *int32 `json:"total,omitempty"`

	// 日志文件列表
	SummaryResults *[]AssetDailySummaryResult `json:"summary_results,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ListAssetDailySummaryLogResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAssetDailySummaryLogResponse struct{}"
	}

	return strings.Join([]string{"ListAssetDailySummaryLogResponse", string(data)}, " ")
}
