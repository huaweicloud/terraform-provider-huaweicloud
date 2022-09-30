package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteIkThesaurusResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteIkThesaurusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteIkThesaurusResponse struct{}"
	}

	return strings.Join([]string{"DeleteIkThesaurusResponse", string(data)}, " ")
}
