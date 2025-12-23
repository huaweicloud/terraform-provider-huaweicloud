package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSocMappings_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_secmaster_soc_mappings.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSocMappings_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.dataclass_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.dataclass_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.classifier_id"),

					resource.TestCheckOutput("name_filter_useful", "true"),
					resource.TestCheckOutput("status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSocMappings_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_soc_mappings" "test" {
  workspace_id = "%[1]s"
}

locals {
  name   = data.huaweicloud_secmaster_soc_mappings.test.data[0].name
  status = data.huaweicloud_secmaster_soc_mappings.test.data[0].status
}

data "huaweicloud_secmaster_soc_mappings" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.name
}

data "huaweicloud_secmaster_soc_mappings" "filter_by_status" {
  workspace_id  = "%[1]s"
  status        = local.status
}

output "name_filter_useful" {
  value = length(data.huaweicloud_secmaster_soc_mappings.filter_by_name.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_soc_mappings.filter_by_name.data[*].name : v == local.name]
  )
}

output "status_filter_useful" {
  value = length(data.huaweicloud_secmaster_soc_mappings.filter_by_status.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_soc_mappings.filter_by_status.data[*].status : v == local.status]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
