package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRansomwareBackups_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_ransomware_backups.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test requires setting a host ID with ransomware protection enabled and generating a backup.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRansomwareBackups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.backup_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.backup_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.backup_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.create_time"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRansomwareBackups_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_ransomware_backups" "test" {
  host_id = "%[1]s"
}

# Filter using name.
locals {
  name = data.huaweicloud_hss_ransomware_backups.test.data_list[0].backup_name
}

data "huaweicloud_hss_ransomware_backups" "name_filter" {
  host_id = "%[1]s"
  name    = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_backups.name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_backups.name_filter.data_list[*].backup_name : v == local.name]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_ransomware_backups" "enterprise_project_id_filter" {
  host_id               = "%[1]s"
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_backups.enterprise_project_id_filter.data_list) > 0
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
