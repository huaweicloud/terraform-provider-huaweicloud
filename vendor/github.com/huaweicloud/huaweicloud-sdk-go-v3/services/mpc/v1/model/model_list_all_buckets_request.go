package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListAllBucketsRequest struct {
}

func (o ListAllBucketsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAllBucketsRequest struct{}"
	}

	return strings.Join([]string{"ListAllBucketsRequest", string(data)}, " ")
}
