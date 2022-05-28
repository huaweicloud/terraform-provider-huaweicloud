package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateCoverByThumbnailResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateCoverByThumbnailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCoverByThumbnailResponse struct{}"
	}

	return strings.Join([]string{"UpdateCoverByThumbnailResponse", string(data)}, " ")
}
