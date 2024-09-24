package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RevokeKubernetesClusterCertResponse Response Object
type RevokeKubernetesClusterCertResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RevokeKubernetesClusterCertResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokeKubernetesClusterCertResponse struct{}"
	}

	return strings.Join([]string{"RevokeKubernetesClusterCertResponse", string(data)}, " ")
}
