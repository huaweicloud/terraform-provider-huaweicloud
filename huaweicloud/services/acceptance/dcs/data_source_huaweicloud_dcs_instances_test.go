package dcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDcsInstance_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dcs_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDCSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDcsInstance_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.vpc_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.maintain_begin"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.maintain_end"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.max_memory"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.used_memory"),

					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_private_ip_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDcsInstance_basic() string {
	return `
data "huaweicloud_dcs_instances" "test" {}

// filter by instance_id
locals {
  instance_id = data.huaweicloud_dcs_instances.test.instances[0].id
}

data "huaweicloud_dcs_instances" "filter_by_instance_id" {
  instance_id = local.instance_id
}

output "is_instance_id_filter_useful" {
  value = length(data.huaweicloud_dcs_instances.filter_by_instance_id.instances) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_instances.filter_by_instance_id.instances[*].id : v == local.instance_id]
  )
}

// filter by name
locals {
  name = data.huaweicloud_dcs_instances.test.instances[0].name
}

data "huaweicloud_dcs_instances" "filter_by_name" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_dcs_instances.filter_by_name.instances) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_instances.filter_by_name.instances[*].name : v == local.name]
  )
}

// filter by status
locals {
  status = data.huaweicloud_dcs_instances.test.instances[0].status
}

data "huaweicloud_dcs_instances" "filter_by_status" {
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_dcs_instances.filter_by_status.instances) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_instances.filter_by_status.instances[*].status : v == local.status]
  )
}

// filter by private_ip
locals {
  private_ip = data.huaweicloud_dcs_instances.test.instances[0].private_ip
}

data "huaweicloud_dcs_instances" "filter_by_private_ip" {
  private_ip = local.private_ip
}

output "is_private_ip_filter_useful" {
  value = length(data.huaweicloud_dcs_instances.filter_by_private_ip.instances) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_instances.filter_by_private_ip.instances[*].private_ip : v == local.private_ip]
  )
}
`
}
