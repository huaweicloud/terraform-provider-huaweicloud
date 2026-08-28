package dds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdsInstance_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dds_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdsInstance_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.role"),

					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_vpc_id_filter_useful", "true"),
					resource.TestCheckOutput("is_subnet_id_filter_useful", "true"),
					resource.TestCheckOutput("is_mode_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDdsInstance_basic() string {
	return `
data "huaweicloud_dds_instances" "test" {}

// filter by instance_id
locals {
  instance_id = data.huaweicloud_dds_instances.test.instances[0].id
}

data "huaweicloud_dds_instances" "filter_by_instance_id" {
  instance_id = local.instance_id
}

output "is_instance_id_filter_useful" {
  value = length(data.huaweicloud_dds_instances.filter_by_instance_id.instances) > 0 && alltrue(
    [for v in data.huaweicloud_dds_instances.filter_by_instance_id.instances[*].id : v == local.instance_id]
  )
}

// filter by name
locals {
  name = data.huaweicloud_dds_instances.test.instances[0].name
}

data "huaweicloud_dds_instances" "filter_by_name" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_dds_instances.filter_by_name.instances) > 0 && alltrue(
    [for v in data.huaweicloud_dds_instances.filter_by_name.instances[*].name : v == local.name]
  )
}

// filter by vpc_id
locals {
  vpc_id = data.huaweicloud_dds_instances.test.instances[0].vpc_id
}

data "huaweicloud_dds_instances" "filter_by_vpc_id" {
  vpc_id = local.vpc_id
}

output "is_vpc_id_filter_useful" {
  value = length(data.huaweicloud_dds_instances.filter_by_vpc_id.instances) > 0 && alltrue(
    [for v in data.huaweicloud_dds_instances.filter_by_vpc_id.instances[*].vpc_id : v == local.vpc_id]
  )
}

// filter by subnet_id
locals {
  subnet_id = data.huaweicloud_dds_instances.test.instances[0].subnet_id
}

data "huaweicloud_dds_instances" "filter_by_subnet_id" {
  subnet_id = local.subnet_id
}

output "is_subnet_id_filter_useful" {
  value = length(data.huaweicloud_dds_instances.filter_by_subnet_id.instances) > 0 && alltrue(
    [for v in data.huaweicloud_dds_instances.filter_by_subnet_id.instances[*].subnet_id : v == local.subnet_id]
  )
}

// filter by mode
locals {
  mode = data.huaweicloud_dds_instances.test.instances[0].mode
}

data "huaweicloud_dds_instances" "filter_by_mode" {
  mode = local.mode
}

output "is_mode_filter_useful" {
  value = length(data.huaweicloud_dds_instances.filter_by_mode.instances) > 0 && alltrue(
    [for v in data.huaweicloud_dds_instances.filter_by_mode.instances[*].mode : v == local.mode]
  )
}
`
}
