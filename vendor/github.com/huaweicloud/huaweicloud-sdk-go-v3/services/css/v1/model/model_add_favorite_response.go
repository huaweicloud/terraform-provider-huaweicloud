package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddFavoriteResponse Response Object
type AddFavoriteResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o AddFavoriteResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddFavoriteResponse struct{}"
	}

	return strings.Join([]string{"AddFavoriteResponse", string(data)}, " ")
}
