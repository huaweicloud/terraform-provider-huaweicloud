package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddFavoriteReq struct {

	// 自定义模板名称。
	Name string `json:"name"`

	Template *AddFavoriteReqTemplate `json:"template"`
}

func (o AddFavoriteReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddFavoriteReq struct{}"
	}

	return strings.Join([]string{"AddFavoriteReq", string(data)}, " ")
}
