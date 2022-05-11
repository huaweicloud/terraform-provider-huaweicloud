package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneShowCatalogResponse struct {

	// 服务目录信息列表。
	Catalog *[]Catalog `json:"catalog,omitempty"`

	Links          *LinksSelf `json:"links,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o KeystoneShowCatalogResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowCatalogResponse struct{}"
	}

	return strings.Join([]string{"KeystoneShowCatalogResponse", string(data)}, " ")
}
