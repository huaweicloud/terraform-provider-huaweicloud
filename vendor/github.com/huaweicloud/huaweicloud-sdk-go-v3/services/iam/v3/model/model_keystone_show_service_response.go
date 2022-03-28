package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneShowServiceResponse struct {
	Service        *Service `json:"service,omitempty"`
	HttpStatusCode int      `json:"-"`
}

func (o KeystoneShowServiceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowServiceResponse struct{}"
	}

	return strings.Join([]string{"KeystoneShowServiceResponse", string(data)}, " ")
}
