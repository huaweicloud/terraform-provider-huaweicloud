package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartTargetClusterConnectivityTestResponse Response Object
type StartTargetClusterConnectivityTestResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartTargetClusterConnectivityTestResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartTargetClusterConnectivityTestResponse struct{}"
	}

	return strings.Join([]string{"StartTargetClusterConnectivityTestResponse", string(data)}, " ")
}
