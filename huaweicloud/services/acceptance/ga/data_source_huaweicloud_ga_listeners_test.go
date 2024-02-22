package ga

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceListeners_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_ga_listeners.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byListenerId   = "data.huaweicloud_ga_listeners.filter_by_listener_id"
		dcByListenerId = acceptance.InitDataSourceCheck(byListenerId)

		byName   = "data.huaweicloud_ga_listeners.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byAcceleratorId   = "data.huaweicloud_ga_listeners.filter_by_accelerator_id"
		dcByAcceleratorId = acceptance.InitDataSourceCheck(byAcceleratorId)

		byStatus   = "data.huaweicloud_ga_listeners.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byProtocol   = "data.huaweicloud_ga_listeners.filter_by_protocol"
		dcByProtocol = acceptance.InitDataSourceCheck(byProtocol)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceListeners_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "listeners.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "listeners.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "listeners.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "listeners.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSourceName, "listeners.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "listeners.0.port_ranges.0.from_port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "listeners.0.port_ranges.0.to_port"),
					dcByListenerId.CheckResourceExists(),
					resource.TestCheckOutput("listener_id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByAcceleratorId.CheckResourceExists(),
					resource.TestCheckOutput("accelerator_id_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByProtocol.CheckResourceExists(),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceListeners_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_accelerator" "test" {
  name        = "%[1]s"
  description = "terraform test"

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
`, name)
}

func testAccDataSourceListeners_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ga_listeners" "test" {
  depends_on = [
    huaweicloud_ga_listener.test
  ]
}

locals {
  listener_id = data.huaweicloud_ga_listeners.test.listeners[0].id
}

data "huaweicloud_ga_listeners" "filter_by_listener_id" {
  listener_id = local.listener_id
}

locals {
  listener_id_filter_result = [
    for v in data.huaweicloud_ga_listeners.filter_by_listener_id.listeners[*].id : v == local.listener_id
  ]
}

output "listener_id_filter_is_useful" {
  value = alltrue(local.listener_id_filter_result) && length(local.listener_id_filter_result) > 0
}

locals {
  name = data.huaweicloud_ga_listeners.test.listeners[0].name
}

data "huaweicloud_ga_listeners" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_ga_listeners.filter_by_name.listeners[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

locals {
  accelerator_id = data.huaweicloud_ga_listeners.test.listeners[0].accelerator_id
}

data "huaweicloud_ga_listeners" "filter_by_accelerator_id" {
  accelerator_id = local.accelerator_id
}

locals {
  accelerator_id_filter_result = [
    for v in data.huaweicloud_ga_listeners.filter_by_accelerator_id.listeners[*].accelerator_id : 
    v == local.accelerator_id
  ]
}

output "accelerator_id_filter_is_useful" {
  value = alltrue(local.accelerator_id_filter_result) && length(local.accelerator_id_filter_result) > 0
}

locals {
  status = data.huaweicloud_ga_listeners.test.listeners[0].status
}

data "huaweicloud_ga_listeners" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ga_listeners.filter_by_status.listeners[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

locals {
  protocol = data.huaweicloud_ga_listeners.test.listeners[0].protocol
}

data "huaweicloud_ga_listeners" "filter_by_protocol" {
  protocol = local.protocol
}

locals {
  protocol_filter_result = [
    for v in data.huaweicloud_ga_listeners.filter_by_protocol.listeners[*].protocol : v == local.protocol
  ]
}

output "protocol_filter_is_useful" {
  value = alltrue(local.protocol_filter_result) && length(local.protocol_filter_result) > 0
}
`, testAccDataSourceListeners_base(name))
}
