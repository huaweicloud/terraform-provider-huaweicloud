package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateAssetReviewTaskResponse struct {

	// 媒资ID
	AssetId *string `json:"asset_id,omitempty"`

	Review         *Review `json:"review,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateAssetReviewTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAssetReviewTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateAssetReviewTaskResponse", string(data)}, " ")
}
