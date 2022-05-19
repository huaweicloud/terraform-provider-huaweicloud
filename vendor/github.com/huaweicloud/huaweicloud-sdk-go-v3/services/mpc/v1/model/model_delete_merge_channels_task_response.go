package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteMergeChannelsTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteMergeChannelsTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMergeChannelsTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteMergeChannelsTaskResponse", string(data)}, " ")
}
