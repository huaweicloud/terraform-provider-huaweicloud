package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneShowVersionResponse struct {
	Version        *Version `json:"version,omitempty"`
	HttpStatusCode int      `json:"-"`
}

func (o KeystoneShowVersionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowVersionResponse struct{}"
	}

	return strings.Join([]string{"KeystoneShowVersionResponse", string(data)}, " ")
}
