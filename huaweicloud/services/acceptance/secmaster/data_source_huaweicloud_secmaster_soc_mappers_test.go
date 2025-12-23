package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSocMapper_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_secmaster_soc_mappers.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterMappingId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSocMapper_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.dataclass_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.dataclass_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.mapper_type_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.update_time"),

					resource.TestCheckOutput("name_filter_useful", "true"),
					resource.TestCheckOutput("start_time_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSocMapper_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_soc_mappers" "test" {
  workspace_id = "%[1]s"
  mapping_id   = "%[2]s"
}

locals {
  name       = data.huaweicloud_secmaster_soc_mappers.test.data[0].name
  start_time = data.huaweicloud_secmaster_soc_mappers.test.data[0].create_time
}

data "huaweicloud_secmaster_soc_mappers" "filter_by_name" {
  workspace_id = "%[1]s"
  mapping_id   = "%[2]s"
  name         = local.name
}

data "huaweicloud_secmaster_soc_mappers" "filter_by_start_time" {
  workspace_id = "%[1]s"
  mapping_id   = "%[2]s"
  start_time   = local.start_time
}

output "name_filter_useful" {
  value = length(data.huaweicloud_secmaster_soc_mappers.filter_by_name.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_soc_mappers.filter_by_name.data[*].name : v == local.name]
  )
}

output "start_time_filter_useful" {
  value = length(data.huaweicloud_secmaster_soc_mappers.filter_by_start_time.data) > 0
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_MAPPING_ID)
}
