package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataspaces_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_secmaster_dataspaces.test"
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
				Config: testDataSourceDataspaces_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.dataspace_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.dataspace_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.dataspace_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.create_by"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.update_time"),

					resource.TestCheckOutput("dataspace_id_filter_useful", "true"),
					resource.TestCheckOutput("dataspace_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataspaces_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_dataspaces" "test" {
  workspace_id = "%[1]s"
}

locals {
  dataspace_id   = data.huaweicloud_secmaster_dataspaces.test.records[0].dataspace_id
  dataspace_name = data.huaweicloud_secmaster_dataspaces.test.records[0].dataspace_name
}

data "huaweicloud_secmaster_dataspaces" "filter_by_dataspace_id" {
  workspace_id = "%[1]s"
  dataspace_id = local.dataspace_id
}

data "huaweicloud_secmaster_dataspaces" "filter_by_dataspace_name" {
  workspace_id   = "%[1]s"
  dataspace_name = local.dataspace_name
}

output "dataspace_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_dataspaces.filter_by_dataspace_id.records) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_dataspaces.filter_by_dataspace_id.records[*].dataspace_id : v == local.dataspace_id]
  )
}

output "dataspace_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_dataspaces.filter_by_dataspace_name.records) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_dataspaces.filter_by_dataspace_name.records[*].dataspace_name : v == local.dataspace_name]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
