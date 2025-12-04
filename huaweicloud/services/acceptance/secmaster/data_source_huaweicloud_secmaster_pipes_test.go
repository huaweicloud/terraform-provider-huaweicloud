package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterPipes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_pipes.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterPipes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.create_by"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dataspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dataspace_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.pipe_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.pipe_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.pipe_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.shards"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.storage_period"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.update_by"),

					resource.TestCheckOutput("is_dataspace_id_filter_useful", "true"),
					resource.TestCheckOutput("is_pipe_id_filter_useful", "true"),
					resource.TestCheckOutput("is_pipe_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterPipes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_pipes" "test" {
  workspace_id = "%[1]s"
}

# Filter by dataspace_id
locals {
  dataspace_id = data.huaweicloud_secmaster_pipes.test.records[0].dataspace_id
}

data "huaweicloud_secmaster_pipes" "filter_by_dataspace_id" {
  workspace_id = "%[1]s"
  dataspace_id = local.dataspace_id
}

locals {
  list_by_dataspace_id = data.huaweicloud_secmaster_pipes.filter_by_dataspace_id.records
}

output "is_dataspace_id_filter_useful" {
  value = length(local.list_by_dataspace_id) > 0 && alltrue(
    [for v in local.list_by_dataspace_id[*].dataspace_id : v == local.dataspace_id]
  )
}

# Filter by pipe_id
locals {
  pipe_id = data.huaweicloud_secmaster_pipes.test.records[0].pipe_id
}

data "huaweicloud_secmaster_pipes" "filter_by_pipe_id" {
  workspace_id = "%[1]s"
  pipe_id      = local.pipe_id
}

locals {
  list_by_pipe_id = data.huaweicloud_secmaster_pipes.filter_by_pipe_id.records
}

output "is_pipe_id_filter_useful" {
  value = length(local.list_by_pipe_id) > 0 && alltrue(
    [for v in local.list_by_pipe_id[*].pipe_id : v == local.pipe_id]
  )
}

# Filter by pipe_name
locals {
  pipe_name = data.huaweicloud_secmaster_pipes.test.records[0].pipe_name
}

data "huaweicloud_secmaster_pipes" "filter_by_pipe_name" {
  workspace_id = "%[1]s"
  pipe_name    = local.pipe_name
}

locals {
  list_by_pipe_name = data.huaweicloud_secmaster_pipes.filter_by_pipe_name.records
}

output "is_pipe_name_filter_useful" {
  value = length(local.list_by_pipe_name) > 0 && alltrue(
    [for v in local.list_by_pipe_name[*].pipe_name : v == local.pipe_name]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
