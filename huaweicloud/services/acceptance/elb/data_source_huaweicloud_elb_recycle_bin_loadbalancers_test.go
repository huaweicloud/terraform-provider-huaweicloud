package elb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceElbRecycleBinInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_recycle_bin_loadbalancers.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbRecycleBinInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.availability_zone_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.l4_flavor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.l7_flavor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.ip_target_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.pools.#"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.provider"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.protection_status"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.vip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.vip_port_id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.operating_status"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.deletion_protection_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.elb_virsubnet_ids.#"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.waf_failure_action"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.guaranteed"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.elb_virsubnet_type"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.listeners.#"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.vip_subnet_cidr_id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.auto_terminate_time"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.updated_at"),
					resource.TestCheckOutput("loadbalancer_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("operating_status_filter_is_useful", "true"),
					resource.TestCheckOutput("guaranteed_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("vip_port_id_filter_is_useful", "true"),
					resource.TestCheckOutput("vip_address_filter_is_useful", "true"),
					resource.TestCheckOutput("vip_subnet_cidr_id_filter_is_useful", "true"),
					resource.TestCheckOutput("availability_zone_list_filter_is_useful", "true"),
					resource.TestCheckOutput("l4_flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("l7_flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("deletion_protection_enable_filter_is_useful", "true"),
					resource.TestCheckOutput("elb_virsubnet_type_filter_is_useful", "true"),
					resource.TestCheckOutput("protection_status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceElbRecycleBinInstances_basic() string {
	return `
data "huaweicloud_elb_recycle_bin_loadbalancers" "test" {}

locals {
  loadbalancer = data.huaweicloud_elb_recycle_bin_loadbalancers.test.loadbalancers[0]
}

locals {
  loadbalancer_id = local.loadbalancer.id
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "loadbalancer_id_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  loadbalancer_id = [local.loadbalancer.id]
}
output "loadbalancer_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.loadbalancer_id_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.loadbalancer_id_filter.loadbalancers[*].id :
    v == local.loadbalancer_id]
  )
}

locals {
  name = local.loadbalancer.name
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "name_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  name = [local.loadbalancer.name]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.name_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.name_filter.loadbalancers[*].name : v == local.name]
  )
}

locals {
  description = local.loadbalancer.description
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "description_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  description = [local.loadbalancer.description]
}
output "description_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.description_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.description_filter.loadbalancers[*].description :
    v == local.description]
  )
}

locals {
  operating_status = local.loadbalancer.operating_status
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "operating_status_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  operating_status = [local.loadbalancer.operating_status]
}
output "operating_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.operating_status_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.operating_status_filter.loadbalancers[*].operating_status :
    v == local.operating_status]
  )
}

locals {
  guaranteed = true
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "guaranteed_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  guaranteed = true
}
output "guaranteed_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.guaranteed_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.guaranteed_filter.loadbalancers[*].guaranteed :
    v == local.guaranteed]
  )
}

locals {
  vpc_id = local.loadbalancer.vpc_id
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "vpc_id_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  vpc_id = [local.loadbalancer.vpc_id]
}
output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.vpc_id_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.vpc_id_filter.loadbalancers[*].vpc_id : v == local.vpc_id]
  )
}

locals {
  vip_port_id = local.loadbalancer.vip_port_id
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "vip_port_id_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  vip_port_id = [local.loadbalancer.vip_port_id]
}
output "vip_port_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.vip_port_id_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.vip_port_id_filter.loadbalancers[*].vip_port_id :
    v == local.vip_port_id]
  )
}

locals {
  vip_address = local.loadbalancer.vip_address
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "vip_address_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  vip_address = [local.loadbalancer.vip_address]
}
output "vip_address_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.vip_address_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.vip_address_filter.loadbalancers[*].vip_address :
    v == local.vip_address]
  )
}

locals {
  vip_subnet_cidr_id = local.loadbalancer.vip_subnet_cidr_id
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "vip_subnet_cidr_id_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  vip_subnet_cidr_id = [local.loadbalancer.vip_subnet_cidr_id]
}
output "vip_subnet_cidr_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.vip_subnet_cidr_id_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.vip_subnet_cidr_id_filter.loadbalancers[*].vip_subnet_cidr_id :
    v == local.vip_subnet_cidr_id]
  )
}

locals {
  availability_zone_list = local.loadbalancer.availability_zone_list[0]
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "availability_zone_list_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  availability_zone_list = [local.loadbalancer.availability_zone_list[0]]
}
output "availability_zone_list_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.availability_zone_list_filter.loadbalancers) > 0
}

locals {
  l4_flavor_id = local.loadbalancer.l4_flavor_id
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "l4_flavor_id_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  l4_flavor_id = [local.loadbalancer.l4_flavor_id]
}
output "l4_flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.l4_flavor_id_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.l4_flavor_id_filter.loadbalancers[*].l4_flavor_id :
    v == local.l4_flavor_id]
  )
}

locals {
  l7_flavor_id = local.loadbalancer.l7_flavor_id
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "l7_flavor_id_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  l7_flavor_id = [local.loadbalancer.l7_flavor_id]
}
output "l7_flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.l7_flavor_id_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.l7_flavor_id_filter.loadbalancers[*].l7_flavor_id :
    v == local.l7_flavor_id]
  )
}

locals {
  enterprise_project_id = local.loadbalancer.enterprise_project_id
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "enterprise_project_id_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  enterprise_project_id = [local.loadbalancer.enterprise_project_id]
}
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.enterprise_project_id_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.enterprise_project_id_filter.loadbalancers[*].enterprise_project_id :
    v == local.enterprise_project_id]
  )
}

locals {
  deletion_protection_enable = false
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "deletion_protection_enable_filter" {
  deletion_protection_enable = false
}
output "deletion_protection_enable_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.deletion_protection_enable_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.deletion_protection_enable_filter.loadbalancers[*].deletion_protection_enable :
    v == local.deletion_protection_enable]
  )
}

locals {
  elb_virsubnet_type = local.loadbalancer.elb_virsubnet_type
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "elb_virsubnet_type_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  elb_virsubnet_type = [local.loadbalancer.elb_virsubnet_type]
}
output "elb_virsubnet_type_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.elb_virsubnet_type_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.elb_virsubnet_type_filter.loadbalancers[*].elb_virsubnet_type :
    v == local.elb_virsubnet_type]
  )
}

locals {
  protection_status = local.loadbalancer.protection_status
}
data "huaweicloud_elb_recycle_bin_loadbalancers" "protection_status_filter" {
  depends_on = [data.huaweicloud_elb_recycle_bin_loadbalancers.test]

  protection_status = [local.loadbalancer.protection_status]
}
output "protection_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_recycle_bin_loadbalancers.protection_status_filter.loadbalancers) > 0 && alltrue(
    [for v in data.huaweicloud_elb_recycle_bin_loadbalancers.protection_status_filter.loadbalancers[*].protection_status :
    v == local.protection_status]
  )
}
`
}
