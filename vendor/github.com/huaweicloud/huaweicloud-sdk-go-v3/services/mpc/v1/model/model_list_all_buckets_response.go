package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAllBucketsResponse Response Object
type ListAllBucketsResponse struct {

	// 桶列表
	Buckets        *[]ObsBucket `json:"buckets,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListAllBucketsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAllBucketsResponse struct{}"
	}

	return strings.Join([]string{"ListAllBucketsResponse", string(data)}, " ")
}
