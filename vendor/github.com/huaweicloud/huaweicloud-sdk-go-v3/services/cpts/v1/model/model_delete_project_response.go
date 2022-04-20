package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteProjectResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteProjectResponse struct{}"
	}

	return strings.Join([]string{"DeleteProjectResponse", string(data)}, " ")
}
