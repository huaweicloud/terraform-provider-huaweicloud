package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstances_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstances_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.updated"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.volume.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.availability_zone"),

					resource.TestCheckOutput("instance_id_filter_useful", "true"),
					resource.TestCheckOutput("name_filter_useful", "true"),
					resource.TestCheckOutput("datastore_type_filter_useful", "true"),
					resource.TestCheckOutput("mode_filter_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_useful", "true"),
					resource.TestCheckOutput("subnet_id_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceInstances_basic = `
data "huaweicloud_geminidb_instances" "test" {}

locals {
  instance_id    = data.huaweicloud_geminidb_instances.test.instances[0].id
  name           = data.huaweicloud_geminidb_instances.test.instances[0].name
  datastore_type = data.huaweicloud_geminidb_instances.test.instances[0].datastore[0].type
  mode           = data.huaweicloud_geminidb_instances.test.instances[0].mode
  vpc_id         = data.huaweicloud_geminidb_instances.test.instances[0].vpc_id
  subnet_id      = data.huaweicloud_geminidb_instances.test.instances[0].subnet_id
}

data "huaweicloud_geminidb_instances" "instance_id_filter" {
  instance_id = local.instance_id
}

output "instance_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_instances.instance_id_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_instances.instance_id_filter.instances[*].id : v == local.instance_id]
  )
}

data "huaweicloud_geminidb_instances" "name_filter" {	
  name = local.name
}

output "name_filter_useful" {
  value = length(data.huaweicloud_geminidb_instances.name_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_instances.name_filter.instances[*].name : v == local.name]
  )
}

data "huaweicloud_geminidb_instances" "datastore_type_filter" {
  datastore_type = local.datastore_type
}

output "datastore_type_filter_useful" {
  value = length(data.huaweicloud_geminidb_instances.datastore_type_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_instances.datastore_type_filter.instances[*].datastore[0].type : v == local.datastore_type]
  )
}

data "huaweicloud_geminidb_instances" "mode_filter" {
  datastore_type = local.datastore_type
  mode           = local.mode
}

output "mode_filter_useful" {
  value = length(data.huaweicloud_geminidb_instances.mode_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_instances.mode_filter.instances[*].mode : v == local.mode]
  )
}

data "huaweicloud_geminidb_instances" "vpc_id_filter" {
  vpc_id = local.vpc_id
}

output "vpc_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_instances.vpc_id_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_instances.vpc_id_filter.instances[*].vpc_id : v == local.vpc_id]
  )
}

data "huaweicloud_geminidb_instances" "subnet_id_filter" {
  subnet_id = local.subnet_id
}

output "subnet_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_instances.subnet_id_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_instances.subnet_id_filter.instances[*].subnet_id : v == local.subnet_id]
  )
}
`
