package deh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceDehInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_deh_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDehInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.dedicated_host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.auto_placement"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.available_memory"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.instance_total"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.instance_uuids.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.available_vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.host_properties.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.host_properties.0.host_type"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.host_properties.0.host_type_name"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.host_properties.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.host_properties.0.cores"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.host_properties.0.sockets"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.host_properties.0.memory"),
					resource.TestCheckResourceAttrSet(dataSource,
						"dedicated_hosts.0.host_properties.0.available_instance_capacities.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"dedicated_hosts.0.host_properties.0.available_instance_capacities.0.flavor"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.tags.%"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.sys_tags.%"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_hosts.0.allocated_at"),
					resource.TestCheckOutput("dedicated_host_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("host_type_filter_is_useful", "true"),
					resource.TestCheckOutput("host_type_name_filter_is_useful", "true"),
					resource.TestCheckOutput("flavor_filter_is_useful", "true"),
					resource.TestCheckOutput("state_filter_is_useful", "true"),
					resource.TestCheckOutput("availability_zone_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDehInstances_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_deh_instance" "test" {
  availability_zone = "cn-southwest-242a"
  name              = "%[2]s"
  host_type         = "s3"
  auto_placement    = "on"

  metadata = {
    "ha_enabled" = "true"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"

  tags = {
    "key"   = "value"
    "owner" = "terraform"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceDataSourceDehInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_deh_instances" "test" {
  depends_on = [huaweicloud_deh_instance.test]
}

locals {
  dedicated_host_id = huaweicloud_deh_instance.test.id
}

data "huaweicloud_deh_instances" "dedicated_host_id_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  dedicated_host_id = huaweicloud_deh_instance.test.id
}

output "dedicated_host_id_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances.dedicated_host_id_filter.dedicated_hosts) > 0 && alltrue(
  [for v in data.huaweicloud_deh_instances.dedicated_host_id_filter.dedicated_hosts[*].dedicated_host_id :
  v == local.dedicated_host_id]
  )
}

locals {
  name = "%[2]s"
}

data "huaweicloud_deh_instances" "name_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  name = "%[2]s"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances.name_filter.dedicated_hosts) > 0 && alltrue(
  [for v in data.huaweicloud_deh_instances.name_filter.dedicated_hosts[*].name : v == local.name]
  )
}

locals {
  host_type = huaweicloud_deh_instance.test.host_type
}

data "huaweicloud_deh_instances" "host_type_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  host_type = huaweicloud_deh_instance.test.host_type
}

output "host_type_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances.host_type_filter.dedicated_hosts) > 0 && alltrue(
  [for v in data.huaweicloud_deh_instances.host_type_filter.dedicated_hosts[*].host_properties[0].host_type :
  v == local.host_type]
  )
}

locals {
  host_type_name = huaweicloud_deh_instance.test.host_properties[0].host_type_name
}

data "huaweicloud_deh_instances" "host_type_name_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  host_type_name = huaweicloud_deh_instance.test.host_properties[0].host_type_name
}

output "host_type_name_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances.host_type_name_filter.dedicated_hosts) > 0 && alltrue(
  [for v in data.huaweicloud_deh_instances.host_type_name_filter.dedicated_hosts[*].host_properties[0].host_type_name :
  v == local.host_type_name]
  )
}

locals {
  flavor = huaweicloud_deh_instance.test.host_properties[0].available_instance_capacities[0].flavor
  dedicated_hosts = data.huaweicloud_deh_instances.flavor_filter.dedicated_hosts
}

data "huaweicloud_deh_instances" "flavor_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  flavor = huaweicloud_deh_instance.test.host_properties[0].available_instance_capacities[0].flavor
}

output "flavor_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances.flavor_filter.dedicated_hosts) > 0 && alltrue(
  [for v in local.dedicated_hosts[*].host_properties[0].available_instance_capacities[*].flavor : contains(v, local.flavor)]
  )
}

locals {
  state = huaweicloud_deh_instance.test.state
}

data "huaweicloud_deh_instances" "state_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  state = huaweicloud_deh_instance.test.state
}

output "state_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances.state_filter.dedicated_hosts) > 0 && alltrue(
  [for v in data.huaweicloud_deh_instances.state_filter.dedicated_hosts[*].state : v == local.state]
  )
}

locals {
  availability_zone = huaweicloud_deh_instance.test.availability_zone
}

data "huaweicloud_deh_instances" "availability_zone_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  availability_zone = huaweicloud_deh_instance.test.availability_zone
}

output "availability_zone_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances.availability_zone_filter.dedicated_hosts) > 0 && alltrue(
  [for v in data.huaweicloud_deh_instances.availability_zone_filter.dedicated_hosts[*].availability_zone :
  v == local.availability_zone]
  )
}

locals {
  tags = "key"
}

data "huaweicloud_deh_instances" "tags_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  tags = "key"
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances.tags_filter.dedicated_hosts) > 0 && alltrue(
  [for v in data.huaweicloud_deh_instances.tags_filter.dedicated_hosts[*].tags : contains(keys(v), local.tags)]
  )
}
`, testDataSourceDataSourceDehInstances_base(name), name)
}
