package eps

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectMigrationRecord_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_enterprise_project_migrate_record.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectMigrationRecord_base(name),
			},
			{
				Config: testAccDataEnterpriseProjectMigrationRecord_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.event_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.operate_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.user_name"),

					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataEnterpriseProjectMigrationRecord_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name                  = "%s"
  enterprise_project_id = "%s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataEnterpriseProjectMigrationRecord_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name                  = "%s"
  enterprise_project_id = "%s"
}

data "huaweicloud_enterprise_project_migrate_record" "test" {
  depends_on = [huaweicloud_cc_connection.test]
}

# Filter by resource_id
locals {
  resource_id = data.huaweicloud_enterprise_project_migrate_record.test.records[0].resource_id
}

data "huaweicloud_enterprise_project_migrate_record" "filter_by_resource_id" {
  resource_id = local.resource_id
}

locals {
  resource_id_filter_result = [
    for v in data.huaweicloud_enterprise_project_migrate_record.filter_by_resource_id.records[*].resource_id : v == local.resource_id
  ]
}

output "resource_id_filter_is_useful" {
  value = alltrue(local.resource_id_filter_result) && length(local.resource_id_filter_result) > 0
}

`, name, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}
