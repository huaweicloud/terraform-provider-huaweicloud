package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneShowCatalogRequest struct {
}

func (o KeystoneShowCatalogRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowCatalogRequest struct{}"
	}

	return strings.Join([]string{"KeystoneShowCatalogRequest", string(data)}, " ")
}
