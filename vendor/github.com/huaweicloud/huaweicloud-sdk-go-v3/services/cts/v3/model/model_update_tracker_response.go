package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateTrackerResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateTrackerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTrackerResponse struct{}"
	}

	return strings.Join([]string{"UpdateTrackerResponse", string(data)}, " ")
}
