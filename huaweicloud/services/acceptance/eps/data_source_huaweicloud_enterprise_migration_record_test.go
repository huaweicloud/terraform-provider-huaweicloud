package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectMigrationRecord_basic(t *testing.T) {
	all := "data.huaweicloud_enterprise_project_migrate_record.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectMigrationRecord_basic_test(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "records.#"),
				),
			},
		},
	})
}

func testAccDataEnterpriseProjectMigrationRecord_basic_test() string {
	return `
data "huaweicloud_enterprise_project_migrate_record" "test" {
  start_time  = "2020-01-01 00:00:00"
  end_time    = "2025-11-19 15:30:30"
  resource_id = "3dde353d-0117-4e1a-a09c-4750f61a3c5d"
}

output "huaweicloud_enterprise_project_migrate_record" {
  value = data.huaweicloud_enterprise_project_migrate_record.test
}
`
}
