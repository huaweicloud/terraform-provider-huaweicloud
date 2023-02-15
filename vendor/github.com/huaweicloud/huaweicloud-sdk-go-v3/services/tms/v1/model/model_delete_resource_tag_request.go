package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteResourceTagRequest struct {
	Body *ReqDeleteTag `json:"body,omitempty"`
}

func (o DeleteResourceTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteResourceTagRequest struct{}"
	}

	return strings.Join([]string{"DeleteResourceTagRequest", string(data)}, " ")
}
