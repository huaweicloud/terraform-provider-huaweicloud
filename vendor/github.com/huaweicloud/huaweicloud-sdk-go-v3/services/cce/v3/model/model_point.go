package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Point struct {
}

func (o Point) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Point struct{}"
	}

	return strings.Join([]string{"Point", string(data)}, " ")
}
