package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAccessCodeRequest Request Object
type ListAccessCodeRequest struct {
}

func (o ListAccessCodeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAccessCodeRequest struct{}"
	}

	return strings.Join([]string{"ListAccessCodeRequest", string(data)}, " ")
}
