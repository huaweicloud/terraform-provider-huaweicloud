package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRansomwareBackupVaults_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_ransomware_backup_vaults.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRansomwareBackupVaults_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vault_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vault_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vault_size"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vault_used"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vault_allocated"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vault_charging_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vault_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.backup_policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.backup_policy_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.backup_policy_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.resources_num"),

					resource.TestCheckOutput("is_vault_name_filter_useful", "true"),
					resource.TestCheckOutput("is_vault_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRansomwareBackupVaults_basic() string {
	return `
data "huaweicloud_hss_ransomware_backup_vaults" "test" {}

# Filter using vault_name.
locals {
  vault_name = data.huaweicloud_hss_ransomware_backup_vaults.test.data_list[0].vault_name
}

data "huaweicloud_hss_ransomware_backup_vaults" "vault_name_filter" {
  vault_name = local.vault_name
}

output "is_vault_name_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_backup_vaults.vault_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_backup_vaults.vault_name_filter.data_list[*].vault_name : v == local.vault_name]
  )
}

# Filter using vault_id.
locals {
  vault_id = data.huaweicloud_hss_ransomware_backup_vaults.test.data_list[0].vault_id
}

data "huaweicloud_hss_ransomware_backup_vaults" "vault_id_filter" {
  vault_id = local.vault_id
}

output "is_vault_id_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_backup_vaults.vault_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_backup_vaults.vault_id_filter.data_list[*].vault_id : v == local.vault_id]
  )
}
`
}
