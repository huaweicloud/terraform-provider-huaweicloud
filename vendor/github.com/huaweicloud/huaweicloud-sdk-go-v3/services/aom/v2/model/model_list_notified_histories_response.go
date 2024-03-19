package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListNotifiedHistoriesResponse Response Object
type ListNotifiedHistoriesResponse struct {

	// 通知历史列表。
	NotifiedHistories *[]NotifiedHistoriesResult `json:"notified_histories,omitempty"`
	HttpStatusCode    int                        `json:"-"`
}

func (o ListNotifiedHistoriesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListNotifiedHistoriesResponse struct{}"
	}

	return strings.Join([]string{"ListNotifiedHistoriesResponse", string(data)}, " ")
}
