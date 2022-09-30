package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type StartKibanaPublicResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StartKibanaPublicResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartKibanaPublicResponse struct{}"
	}

	return strings.Join([]string{"StartKibanaPublicResponse", string(data)}, " ")
}
