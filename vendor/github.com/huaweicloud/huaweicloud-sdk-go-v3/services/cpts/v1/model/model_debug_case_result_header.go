package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DebugCaseResultHeader struct {
	// Connection

	Connection *string `json:"Connection,omitempty"`
	// Content-Type

	ContentType *string `json:"Content-Type,omitempty"`
	// Host

	Host *string `json:"Host,omitempty"`
}

func (o DebugCaseResultHeader) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DebugCaseResultHeader struct{}"
	}

	return strings.Join([]string{"DebugCaseResultHeader", string(data)}, " ")
}
