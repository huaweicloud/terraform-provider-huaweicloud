package ga

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceAddressGroups_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_ga_address_groups.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byAddressGroupId   = "data.huaweicloud_ga_address_groups.filter_by_address_group_id"
		dcByAddressGroupId = acceptance.InitDataSourceCheck(byAddressGroupId)

		byName   = "data.huaweicloud_ga_address_groups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_ga_address_groups.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAddressGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.0.ip_addresses.0.cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.0.associated_listeners.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.0.associated_listeners.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_groups.0.updated_at"),

					dcByAddressGroupId.CheckResourceExists(),
					resource.TestCheckOutput("address_group_id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAddressGroups_base(name string) string {
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

resource "huaweicloud_ga_address_group" "test" {
  name        = "%[1]s"
  description = "terraform create"

  ip_addresses {
    cidr        = "192.168.1.0/24"
    description = "The IP addresses included in the address group"
  }

  listeners {
    id   = huaweicloud_ga_listener.test.id
    type = "WHITE"
  }
}
`, name)
}

func testAccDataSourceAddressGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ga_address_groups" "test" {
  depends_on = [
    huaweicloud_ga_address_group.test
  ]
}

locals {
  address_group_id = data.huaweicloud_ga_address_groups.test.address_groups[0].id
}

data "huaweicloud_ga_address_groups" "filter_by_address_group_id" {
  address_group_id = local.address_group_id
}

locals {
  address_group_id_filter_result = [
    for v in data.huaweicloud_ga_address_groups.filter_by_address_group_id.address_groups[*].id : 
    v == local.address_group_id
  ]
}

output "address_group_id_filter_is_useful" {
  value = alltrue(local.address_group_id_filter_result) && length(local.address_group_id_filter_result) > 0
}

locals {
  name = data.huaweicloud_ga_address_groups.test.address_groups[0].name
}

data "huaweicloud_ga_address_groups" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_ga_address_groups.filter_by_name.address_groups[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

locals {
  status = data.huaweicloud_ga_address_groups.test.address_groups[0].status
}

data "huaweicloud_ga_address_groups" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ga_address_groups.filter_by_status.address_groups[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testAccDataSourceAddressGroups_base(name))
}
