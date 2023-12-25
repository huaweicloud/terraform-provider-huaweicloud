package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowReleaseHistoryResponse Response Object
type ShowReleaseHistoryResponse struct {
	Body           *[]ReleaseResp `json:"body,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ShowReleaseHistoryResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowReleaseHistoryResponse struct{}"
	}

	return strings.Join([]string{"ShowReleaseHistoryResponse", string(data)}, " ")
}
