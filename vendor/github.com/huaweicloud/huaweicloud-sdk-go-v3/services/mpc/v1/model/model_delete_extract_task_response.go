package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteExtractTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteExtractTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteExtractTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteExtractTaskResponse", string(data)}, " ")
}
