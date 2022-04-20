package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DebugCaseReturnHeader struct {
	// Connection

	Connection *string `json:"Connection,omitempty"`
	// Content-Length

	ContentLength *string `json:"Content-Length,omitempty"`
	// Content-Type

	ContentType *string `json:"Content-Type,omitempty"`
	// Date

	Date *string `json:"Date,omitempty"`
	// Vary

	Vary *string `json:"Vary,omitempty"`
}

func (o DebugCaseReturnHeader) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DebugCaseReturnHeader struct{}"
	}

	return strings.Join([]string{"DebugCaseReturnHeader", string(data)}, " ")
}
