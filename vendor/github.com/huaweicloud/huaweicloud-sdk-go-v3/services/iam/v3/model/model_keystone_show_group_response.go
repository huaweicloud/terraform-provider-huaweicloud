package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneShowGroupResponse struct {
	Group          *KeystoneGroupResult `json:"group,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o KeystoneShowGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowGroupResponse struct{}"
	}

	return strings.Join([]string{"KeystoneShowGroupResponse", string(data)}, " ")
}
