package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTranscodeProductRequest Request Object
type DeleteTranscodeProductRequest struct {
	Body *DeleteTranscodeProductReq `json:"body,omitempty"`
}

func (o DeleteTranscodeProductRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTranscodeProductRequest struct{}"
	}

	return strings.Join([]string{"DeleteTranscodeProductRequest", string(data)}, " ")
}
