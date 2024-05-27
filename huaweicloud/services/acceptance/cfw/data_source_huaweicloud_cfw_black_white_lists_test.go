package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwBlackWhiteLists_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_black_white_lists.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwBlackWhiteLists_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.list_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.address_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.direction"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.address"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_protocol_filter_useful", "true"),
					resource.TestCheckOutput("is_address_filter_useful", "true"),
					resource.TestCheckOutput("is_direction_filter_useful", "true"),
					resource.TestCheckOutput("is_port_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCfwBlackWhiteLists_basic() string {
	return fmt.Sprintf(`
%[1]s

locals {
  id        = huaweicloud_cfw_black_white_list.l1.id
  protocol  = huaweicloud_cfw_black_white_list.l2.protocol
  address   = huaweicloud_cfw_black_white_list.l3.address
  direction = tostring(huaweicloud_cfw_black_white_list.l3.direction)
  port      = huaweicloud_cfw_black_white_list.l3.port
}

data "huaweicloud_cfw_black_white_lists" "test" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type = 4

  depends_on = [
    huaweicloud_cfw_black_white_list.l1,
    huaweicloud_cfw_black_white_list.l2,
    huaweicloud_cfw_black_white_list.l3,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_cfw_black_white_lists.test.records) >= 2
}

data "huaweicloud_cfw_black_white_lists" "filter_by_id" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type = 4
  list_id   = local.id

  depends_on = [
    huaweicloud_cfw_black_white_list.l1,
    huaweicloud_cfw_black_white_list.l2,
    huaweicloud_cfw_black_white_list.l3,
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cfw_black_white_lists.filter_by_id.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_black_white_lists.filter_by_id.records[*] : v.list_id == local.id]
  )
}

data "huaweicloud_cfw_black_white_lists" "filter_by_protocol" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type = 4
  protocol  = local.protocol

  depends_on = [
    huaweicloud_cfw_black_white_list.l1,
    huaweicloud_cfw_black_white_list.l2,
    huaweicloud_cfw_black_white_list.l3,
  ]
}

output "is_protocol_filter_useful" {
  value = length(data.huaweicloud_cfw_black_white_lists.filter_by_protocol.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_black_white_lists.filter_by_protocol.records[*] : v.protocol == local.protocol]
  )
}


data "huaweicloud_cfw_black_white_lists" "filter_by_address" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type = 5
  address   = local.address

  depends_on = [
    huaweicloud_cfw_black_white_list.l1,
    huaweicloud_cfw_black_white_list.l2,
    huaweicloud_cfw_black_white_list.l3,
  ]
}

output "is_address_filter_useful" {
  value = length(data.huaweicloud_cfw_black_white_lists.filter_by_address.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_black_white_lists.filter_by_address.records[*] : v.address == local.address]
  )
}

data "huaweicloud_cfw_black_white_lists" "filter_by_direction" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type = 5
  direction = local.direction

  depends_on = [
    huaweicloud_cfw_black_white_list.l1,
    huaweicloud_cfw_black_white_list.l2,
    huaweicloud_cfw_black_white_list.l3,
  ]
}

output "is_direction_filter_useful" {
  value = length(data.huaweicloud_cfw_black_white_lists.filter_by_direction.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_black_white_lists.filter_by_direction.records[*] : v.direction == local.direction]
  )
}

data "huaweicloud_cfw_black_white_lists" "filter_by_port" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type = 5
  port      = local.port

  depends_on = [
    huaweicloud_cfw_black_white_list.l1,
    huaweicloud_cfw_black_white_list.l2,
    huaweicloud_cfw_black_white_list.l3,
  ]
}

output "is_port_filter_useful" {
  value = length(data.huaweicloud_cfw_black_white_lists.filter_by_port.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_black_white_lists.filter_by_port.records[*] : v.port == local.port]
  )
}
`, testDataSourceCfwBlackWhiteLists_base())
}

func testDataSourceCfwBlackWhiteLists_base() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_black_white_list" "l1" {
  object_id    = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type    = 4
  direction    = 1
  protocol     = 6
  port         = "22"
  address_type = 0
  address      = "1.3.1.3"
}

resource "huaweicloud_cfw_black_white_list" "l2" {
  object_id    = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type    = 4
  direction    = 1
  protocol     = 17
  port         = "80"
  address_type = 0
  address      = "1.2.1.1"
}

resource "huaweicloud_cfw_black_white_list" "l3" {
  object_id    = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type    = 5
  direction    = 0
  protocol     = 6
  port         = "80"
  address_type = 0
  address      = "2.2.1.1"
}
`, testAccDatasourceFirewalls_basic())
}
