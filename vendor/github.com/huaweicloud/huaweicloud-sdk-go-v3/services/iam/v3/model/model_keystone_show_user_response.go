package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneShowUserResponse struct {
	User           *KeystoneShowUserResult `json:"user,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o KeystoneShowUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowUserResponse struct{}"
	}

	return strings.Join([]string{"KeystoneShowUserResponse", string(data)}, " ")
}
