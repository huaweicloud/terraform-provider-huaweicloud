package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddFavoriteReq struct {

	// 配置文件名称。4～32个字符，只能包含数字、字母、中划线和下划线，且必须以字母开头。
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
