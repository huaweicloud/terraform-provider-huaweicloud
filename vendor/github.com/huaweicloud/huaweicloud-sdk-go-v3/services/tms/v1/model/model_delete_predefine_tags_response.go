package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeletePredefineTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeletePredefineTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePredefineTagsResponse struct{}"
	}

	return strings.Join([]string{"DeletePredefineTagsResponse", string(data)}, " ")
}
