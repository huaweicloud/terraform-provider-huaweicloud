package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListReleasesResponse Response Object
type ListReleasesResponse struct {
	Body           *[]ReleaseResp `json:"body,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListReleasesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListReleasesResponse struct{}"
	}

	return strings.Join([]string{"ListReleasesResponse", string(data)}, " ")
}
