package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneShowMappingResponse struct {
	Mapping        *MappingResult `json:"mapping,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o KeystoneShowMappingResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowMappingResponse struct{}"
	}

	return strings.Join([]string{"KeystoneShowMappingResponse", string(data)}, " ")
}
