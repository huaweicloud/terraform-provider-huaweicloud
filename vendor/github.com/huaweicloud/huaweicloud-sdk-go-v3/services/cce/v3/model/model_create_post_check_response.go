package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePostCheckResponse Response Object
type CreatePostCheckResponse struct {

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	Metadata *PostcheckCluserResponseMetadata `json:"metadata,omitempty"`

	Spec *PostcheckSpec `json:"spec,omitempty"`

	Status         *PostcheckClusterResponseBodyStatus `json:"status,omitempty"`
	HttpStatusCode int                                 `json:"-"`
}

func (o CreatePostCheckResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePostCheckResponse struct{}"
	}

	return strings.Join([]string{"CreatePostCheckResponse", string(data)}, " ")
}
