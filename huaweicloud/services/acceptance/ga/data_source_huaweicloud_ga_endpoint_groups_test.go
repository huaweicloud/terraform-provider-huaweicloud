package ga

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceEndpointGroups_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_ga_endpoint_groups.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byEndpointGroupId   = "data.huaweicloud_ga_endpoint_groups.filter_by_endpoint_group_id"
		dcByEndpointGroupId = acceptance.InitDataSourceCheck(byEndpointGroupId)

		byName   = "data.huaweicloud_ga_endpoint_groups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_ga_endpoint_groups.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byListenerId   = "data.huaweicloud_ga_endpoint_groups.filter_by_listener_id"
		dcByListenerId = acceptance.InitDataSourceCheck(byListenerId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEndpointGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint_groups.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint_groups.0.traffic_dial_percentage"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint_groups.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint_groups.0.listener_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint_groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint_groups.0.status"),

					dcByEndpointGroupId.CheckResourceExists(),
					resource.TestCheckOutput("endpoint_group_id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByListenerId.CheckResourceExists(),
					resource.TestCheckOutput("listener_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceEndpointGroups_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_accelerator" "test" {
  name        = "%[1]s"
  description = "Terraform test"

  ip_sets {
    area = "CM"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_ga_listener" "test" {
  accelerator_id = huaweicloud_ga_accelerator.test.id
  name           = "%[1]s"
  protocol       = "TCP"
  description    = "Terraform test"

  port_ranges {
    from_port = 90
    to_port   = 99
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_ga_endpoint_group" "test" {
  name        = "%[1]s"
  description = "created by terraform"
  region_id   = "cn-south-1"

  listeners {
    id = huaweicloud_ga_listener.test.id
  }
}
`, name)
}

func testAccDataSourceEndpointGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ga_endpoint_groups" "test" {
  depends_on = [
    huaweicloud_ga_endpoint_group.test
  ]
}

# Filter by endpoint_group_id
locals {
  endpoint_group_id = data.huaweicloud_ga_endpoint_groups.test.endpoint_groups[0].id
}

data "huaweicloud_ga_endpoint_groups" "filter_by_endpoint_group_id" {
  endpoint_group_id = local.endpoint_group_id
}

locals {
  endpoint_group_id_filter_result = [
    for v in data.huaweicloud_ga_endpoint_groups.filter_by_endpoint_group_id.endpoint_groups[*].id : 
    v == local.endpoint_group_id
  ]
}

output "endpoint_group_id_filter_is_useful" {
  value = alltrue(local.endpoint_group_id_filter_result) && length(local.endpoint_group_id_filter_result) > 0
}

# Filter by name
locals {
  name = data.huaweicloud_ga_endpoint_groups.test.endpoint_groups[0].name
}

data "huaweicloud_ga_endpoint_groups" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_ga_endpoint_groups.filter_by_name.endpoint_groups[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by status
locals {
  status = data.huaweicloud_ga_endpoint_groups.test.endpoint_groups[0].status
}

data "huaweicloud_ga_endpoint_groups" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ga_endpoint_groups.filter_by_status.endpoint_groups[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

# Filter by listener_id
locals {
  listener_id = data.huaweicloud_ga_endpoint_groups.test.endpoint_groups[0].listener_id
}

data "huaweicloud_ga_endpoint_groups" "filter_by_listener_id" {
  listener_id = local.listener_id
}

locals {
  listener_id_filter_result = [
    for v in data.huaweicloud_ga_endpoint_groups.filter_by_listener_id.endpoint_groups[*].listener_id : 
    v == local.listener_id
  ]
}

output "listener_id_filter_is_useful" {
  value = alltrue(local.listener_id_filter_result) && length(local.listener_id_filter_result) > 0
}
`, testAccDataSourceEndpointGroups_base(name))
}
