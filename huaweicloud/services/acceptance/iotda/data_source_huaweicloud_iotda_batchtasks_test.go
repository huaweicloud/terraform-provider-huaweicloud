package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBatchTasks_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_batchtasks.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBatchTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "batchtasks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "batchtasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "batchtasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "batchtasks.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "batchtasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "batchtasks.0.task_progress.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "batchtasks.0.created_at"),

					resource.TestCheckOutput("is_space_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceBatchTasks_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_batchtasks" "test" {
  type = huaweicloud_iotda_batchtask.test_freeze.type
}

# Filter using space ID.
data "huaweicloud_iotda_batchtasks" "space_id_filter" {
  depends_on = [huaweicloud_iotda_batchtask.test_freeze]

  type     = data.huaweicloud_iotda_batchtasks.test.batchtasks[0].type
  space_id = huaweicloud_iotda_space.test.id
}

output "is_space_id_filter_useful" {
  value = length(data.huaweicloud_iotda_batchtasks.space_id_filter.batchtasks) > 0
}

# Filter using status.
locals {
  status = data.huaweicloud_iotda_batchtasks.test.batchtasks[0].status
}

data "huaweicloud_iotda_batchtasks" "status_filter" {
  type   = data.huaweicloud_iotda_batchtasks.test.batchtasks[0].type
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_iotda_batchtasks.status_filter.batchtasks) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_batchtasks.status_filter.batchtasks[*].status : v == local.status]
  )
}
`, testBatchTask_basic(name, name))
}
