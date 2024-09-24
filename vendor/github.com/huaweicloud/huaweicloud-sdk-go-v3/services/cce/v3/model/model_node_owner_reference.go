package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeOwnerReference 节点的属主对象信息
type NodeOwnerReference struct {

	// 节点池名称
	NodepoolName *string `json:"nodepoolName,omitempty"`

	// 节点池UID
	NodepoolID *string `json:"nodepoolID,omitempty"`
}

func (o NodeOwnerReference) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeOwnerReference struct{}"
	}

	return strings.Join([]string{"NodeOwnerReference", string(data)}, " ")
}
