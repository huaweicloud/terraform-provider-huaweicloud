package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResultId 病毒查杀结果ID
type ResultId struct {
}

func (o ResultId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResultId struct{}"
	}

	return strings.Join([]string{"ResultId", string(data)}, " ")
}
