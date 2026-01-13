package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRansomwareBackupPolicies_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_ransomware_backup_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRansomwareBackupPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operation_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operation_definition.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operation_definition.0.day_backups"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operation_definition.0.max_backups"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operation_definition.0.month_backups"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operation_definition.0.retention_duration_days"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operation_definition.0.timezone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operation_definition.0.week_backups"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operation_definition.0.year_backups"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.trigger.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.trigger.0.properties.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.trigger.0.properties.0.pattern.#"),
				),
			},
		},
	})
}

func testDataSourceRansomwareBackupPolicies_basic() string {
	return `
data "huaweicloud_hss_ransomware_backup_policies" "test" {
  enterprise_project_id = "all_granted_eps"
}
`
}
