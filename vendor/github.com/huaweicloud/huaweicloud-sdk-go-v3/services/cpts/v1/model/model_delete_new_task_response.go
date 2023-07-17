package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteNewTaskResponse Response Object
type DeleteNewTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteNewTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteNewTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteNewTaskResponse", string(data)}, " ")
}
