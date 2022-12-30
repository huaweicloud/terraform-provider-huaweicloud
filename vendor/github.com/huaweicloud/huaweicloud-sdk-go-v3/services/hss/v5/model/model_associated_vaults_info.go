package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AssociatedVaultsInfo struct {

	// 关联的远端存储库ID
	DestinationVaultId *string `json:"destination_vault_id,omitempty"`

	// 存储库ID
	VaultId *string `json:"vault_id,omitempty"`
}

func (o AssociatedVaultsInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociatedVaultsInfo struct{}"
	}

	return strings.Join([]string{"AssociatedVaultsInfo", string(data)}, " ")
}
