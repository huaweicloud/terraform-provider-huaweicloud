package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateEsListenerResponse Response Object
type UpdateEsListenerResponse struct {
	Listener       *EsListenerResponse `json:"listener,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o UpdateEsListenerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateEsListenerResponse struct{}"
	}

	return strings.Join([]string{"UpdateEsListenerResponse", string(data)}, " ")
}
