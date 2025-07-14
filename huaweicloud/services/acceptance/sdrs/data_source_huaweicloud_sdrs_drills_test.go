package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSdrsDrills_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sdrs_drills.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSdrsDrills_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "disaster_recovery_drills.#"),
					resource.TestCheckResourceAttrSet(dataSource, "disaster_recovery_drills.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "disaster_recovery_drills.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "disaster_recovery_drills.0.drill_servers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "disaster_recovery_drills.0.drill_vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "disaster_recovery_drills.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "disaster_recovery_drills.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "disaster_recovery_drills.0.server_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "disaster_recovery_drills.0.status"),

					resource.TestCheckOutput("is_drill_vpc_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_server_group_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSdrsDrills_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_sdrs_drills" "test" {
  depends_on = [huaweicloud_sdrs_drill.test]
}

# Filter by drill_vpc_id
locals {
  drill_vpc_id = data.huaweicloud_sdrs_drills.test.disaster_recovery_drills[0].drill_vpc_id
}

data "huaweicloud_sdrs_drills" "filter_by_drill_vpc_id" {
  drill_vpc_id = local.drill_vpc_id
}

locals {
  drill_vpc_id_filter_result = [
    for v in data.huaweicloud_sdrs_drills.filter_by_drill_vpc_id.disaster_recovery_drills[*].drill_vpc_id : v == local.drill_vpc_id
  ]
}

output "is_drill_vpc_id_filter_useful" {
  value = length(local.drill_vpc_id_filter_result) > 0 && alltrue(local.drill_vpc_id_filter_result)
}

# Filter by name
locals {
  name = data.huaweicloud_sdrs_drills.test.disaster_recovery_drills[0].name
}

data "huaweicloud_sdrs_drills" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_sdrs_drills.filter_by_name.disaster_recovery_drills[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by server_group_id
locals {
  server_group_id = data.huaweicloud_sdrs_drills.test.disaster_recovery_drills[0].server_group_id
}

data "huaweicloud_sdrs_drills" "filter_by_server_group_id" {
  server_group_id = local.server_group_id
}

locals {
  server_group_id_filter_result = [
    for v in data.huaweicloud_sdrs_drills.filter_by_server_group_id.disaster_recovery_drills[*].server_group_id :
    v == local.server_group_id
  ]
}

output "is_server_group_id_filter_useful" {
  value = length(local.server_group_id_filter_result) > 0 && alltrue(local.server_group_id_filter_result)
}

# Filter by status
locals {
  status = data.huaweicloud_sdrs_drills.test.disaster_recovery_drills[0].status
}

data "huaweicloud_sdrs_drills" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_sdrs_drills.filter_by_status.disaster_recovery_drills[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}
`, testDrill_basic(name))
}
