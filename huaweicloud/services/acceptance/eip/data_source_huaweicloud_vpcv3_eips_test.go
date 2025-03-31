package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceEipVpcv3Eips_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpcv3_eips.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceEipVpcv3Eips_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.#"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.alias"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.ip_version"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.publicip_pool_name"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.billing_info"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.vnic.#"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.vnic.0.mac"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.vnic.0.instance_type"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.vnic.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.vnic.0.port_id"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.vnic.0.device_owner"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.vnic.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.vnic.0.private_ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.bandwidth.#"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.bandwidth.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.bandwidth.0.share_type"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.bandwidth.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.bandwidth.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.bandwidth.0.billing_info"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.publicip_pool_id"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.associate_instance_type"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.associate_instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.public_ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.public_ipv6_address"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.allow_share_bandwidth_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "publicips.0.updated_at"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("alias_filter_is_useful", "true"),
					resource.TestCheckOutput("alias_like_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_version_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("public_ip_address_filter_is_useful", "true"),
					resource.TestCheckOutput("public_ip_address_like_filter_is_useful", "true"),
					resource.TestCheckOutput("public_ipv6_address_filter_is_useful", "true"),
					resource.TestCheckOutput("public_ipv6_address_like_filter_is_useful", "true"),
					resource.TestCheckOutput("publicip_pool_name_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_key_and_dir_filter_is_useful", "true"),
					resource.TestCheckOutput("vnic_private_ip_address_filter_is_useful", "true"),
					resource.TestCheckOutput("vnic_private_ip_address_like_filter_is_useful", "true"),
					resource.TestCheckOutput("vnic_device_owner_filter_is_useful", "true"),
					resource.TestCheckOutput("vnic_device_owner_prefixlike_filter_is_useful", "true"),
					resource.TestCheckOutput("vnic_vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("vnic_port_id_filter_is_useful", "true"),
					resource.TestCheckOutput("vnic_instance_type_filter_is_useful", "true"),
					resource.TestCheckOutput("vnic_instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("public_ipv6_address_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_id_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_name_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_name_like_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_size_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_share_type_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_charge_mode_filter_is_useful", "true"),
					resource.TestCheckOutput("billing_info_filter_is_useful", "true"),
					resource.TestCheckOutput("billing_mode_filter_is_useful", "true"),
					resource.TestCheckOutput("associate_instance_type_filter_is_useful", "true"),
					resource.TestCheckOutput("associate_instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("public_border_group_filter_is_useful", "true"),
					resource.TestCheckOutput("allow_share_bandwidth_type_any_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceEipVpcv3Eips_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc_eip" "test" {
  name = "%[2]s"

  publicip {
    type       = "5_bgp"
    ip_version = 6
  }

  bandwidth {
    name        = "%[2]s"
    share_type  = "PER"
    size        = 5
    charge_mode = "bandwidth"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "ReplicaSet"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "replica"
    storage   = "ULTRAHIGH"
    num       = 1
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.repset"
  }
}

resource "huaweicloud_dds_instance_eip_associate" "test" { 
  instance_id = huaweicloud_dds_instance.instance.id
  node_id     = [for v in huaweicloud_dds_instance.instance.nodes : v.id if v.role == "Primary"][0]
  public_ip   = huaweicloud_vpc_eip.test.address
}`, common.TestBaseNetwork(rName), rName)
}

func testDataSourceDataSourceEipVpcv3Eips_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcv3_eips" "test" {
  depends_on = [huaweicloud_dds_instance_eip_associate.test]
}

locals {
  type = data.huaweicloud_vpcv3_eips.test.publicips[0].type
}

data "huaweicloud_vpcv3_eips" "type_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  type = [data.huaweicloud_vpcv3_eips.test.publicips[0].type]
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.type_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.type_filter.publicips[*].type : v == local.type]
  )
}

locals {
  alias = "%[2]s"
}

data "huaweicloud_vpcv3_eips" "alias_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test
  ]

  alias = ["%[2]s"]
}

output "alias_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.alias_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.alias_filter.publicips[*].alias : v == local.alias]
  )
}

data "huaweicloud_vpcv3_eips" "alias_like_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
  ]

  alias_like = split("_", "%[2]s")[0]
}

output "alias_like_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.alias_like_filter.publicips) > 0
}

locals {
  ip_version = data.huaweicloud_vpcv3_eips.test.publicips[0].ip_version
}

data "huaweicloud_vpcv3_eips" "ip_version_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  ip_version = [data.huaweicloud_vpcv3_eips.test.publicips[0].ip_version]
}

output "ip_version_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.ip_version_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.ip_version_filter.publicips[*].ip_version : v == local.ip_version]
  )
}

locals {
  status = data.huaweicloud_vpcv3_eips.test.publicips[0].status
}

data "huaweicloud_vpcv3_eips" "status_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  status = [data.huaweicloud_vpcv3_eips.test.publicips[0].status]
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.status_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.status_filter.publicips[*].status : v == local.status]
  )
}

locals {
  public_ip_address = data.huaweicloud_vpcv3_eips.test.publicips[0].public_ip_address
}

data "huaweicloud_vpcv3_eips" "public_ip_address_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  public_ip_address = [data.huaweicloud_vpcv3_eips.test.publicips[0].public_ip_address]
}

output "public_ip_address_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.public_ip_address_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.public_ip_address_filter.publicips[*].public_ip_address : v == local.public_ip_address]
  )
}

data "huaweicloud_vpcv3_eips" "public_ip_address_like_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  public_ip_address_like = split(".", data.huaweicloud_vpcv3_eips.test.publicips[0].public_ip_address)[0]
}

output "public_ip_address_like_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.public_ip_address_like_filter.publicips) > 0
}

locals {
  public_ipv6_address = data.huaweicloud_vpcv3_eips.test.publicips[0].public_ipv6_address
}

data "huaweicloud_vpcv3_eips" "public_ipv6_address_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  public_ipv6_address = [data.huaweicloud_vpcv3_eips.test.publicips[0].public_ipv6_address]
}

output "public_ipv6_address_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.public_ipv6_address_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.public_ipv6_address_filter.publicips[*].public_ipv6_address :
  v == local.public_ipv6_address]
  )
}

data "huaweicloud_vpcv3_eips" "public_ipv6_address_like_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  public_ipv6_address_like = split(":", data.huaweicloud_vpcv3_eips.test.publicips[0].public_ipv6_address)[0]
}

output "public_ipv6_address_like_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.public_ipv6_address_like_filter.publicips) > 0
}

locals {
  publicip_pool_name = data.huaweicloud_vpcv3_eips.test.publicips[0].publicip_pool_name
}

data "huaweicloud_vpcv3_eips" "publicip_pool_name_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  publicip_pool_name = [data.huaweicloud_vpcv3_eips.test.publicips[0].publicip_pool_name]
}

output "publicip_pool_name_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.publicip_pool_name_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.publicip_pool_name_filter.publicips[*].publicip_pool_name :
  v == local.publicip_pool_name]
  )
}

data "huaweicloud_vpcv3_eips" "sort_key_and_dir_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test
  ]

  sort_key = "id"
  sort_dir = "asc"
}

output "sort_key_and_dir_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.sort_key_and_dir_filter.publicips) > 0
}

locals {
  vnic_private_ip_address = data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].private_ip_address
}

data "huaweicloud_vpcv3_eips" "vnic_private_ip_address_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  vnic_private_ip_address = [data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].private_ip_address]
}

output "vnic_private_ip_address_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.vnic_private_ip_address_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.vnic_private_ip_address_filter.publicips[*].vnic.0.private_ip_address :
  v == local.vnic_private_ip_address]
  )
}

data "huaweicloud_vpcv3_eips" "vnic_private_ip_address_like_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  vnic_private_ip_address_like = split(".", data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].private_ip_address)[0]
}

output "vnic_private_ip_address_like_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.vnic_private_ip_address_like_filter.publicips) > 0
}

locals {
  vnic_device_owner = data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].device_owner
}

data "huaweicloud_vpcv3_eips" "vnic_device_owner_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  vnic_device_owner = [data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].device_owner]
}

output "vnic_device_owner_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.vnic_device_owner_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.vnic_device_owner_filter.publicips[*].vnic.0.device_owner : v == local.vnic_device_owner]
  )
}

data "huaweicloud_vpcv3_eips" "vnic_device_owner_prefixlike_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  vnic_device_owner_prefixlike = split(":", data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].device_owner)[0]
}

output "vnic_device_owner_prefixlike_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.vnic_device_owner_prefixlike_filter.publicips) > 0
}

locals {
  vnic_vpc_id = data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].vpc_id
}

data "huaweicloud_vpcv3_eips" "vnic_vpc_id_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  vnic_vpc_id = [data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].vpc_id]
}

output "vnic_vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.vnic_vpc_id_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.vnic_vpc_id_filter.publicips[*].vnic.0.vpc_id : v == local.vnic_vpc_id]
  )
}

locals {
  vnic_port_id = data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].port_id
}

data "huaweicloud_vpcv3_eips" "vnic_port_id_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  vnic_port_id = [data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].port_id]
}

output "vnic_port_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.vnic_port_id_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.vnic_port_id_filter.publicips[*].vnic.0.port_id : v == local.vnic_port_id]
  )
}

locals {
  vnic_instance_type = data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].instance_type
}

data "huaweicloud_vpcv3_eips" "vnic_instance_type_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  vnic_instance_type = [data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].instance_type]
}

output "vnic_instance_type_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.vnic_instance_type_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.vnic_instance_type_filter.publicips[*].vnic.0.instance_type :
  v == local.vnic_instance_type]
  )
}

locals {
  vnic_instance_id = data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].instance_id 
}

data "huaweicloud_vpcv3_eips" "vnic_instance_id_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  vnic_instance_id = [data.huaweicloud_vpcv3_eips.test.publicips[0].vnic[0].instance_id ]
}

output "vnic_instance_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.vnic_instance_id_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.vnic_instance_id_filter.publicips[*].vnic.0.instance_id : v == local.vnic_instance_id]
  )
}

locals {
  bandwidth_id = data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].id
}

data "huaweicloud_vpcv3_eips" "bandwidth_id_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  bandwidth_id = [data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].id]
}

output "bandwidth_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.bandwidth_id_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.bandwidth_id_filter.publicips[*].bandwidth.0.id : v == local.bandwidth_id]
  )
}

locals {
  bandwidth_name = data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].name
}

data "huaweicloud_vpcv3_eips" "bandwidth_name_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  bandwidth_name = [data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].name]
}

output "bandwidth_name_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.bandwidth_name_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.bandwidth_name_filter.publicips[*].bandwidth.0.name : v == local.bandwidth_name]
  )
}

data "huaweicloud_vpcv3_eips" "bandwidth_name_like_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  bandwidth_name_like = [split("-", data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].name)[0]]
}

output "bandwidth_name_like_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.bandwidth_name_like_filter.publicips) > 0
}

locals {
  bandwidth_size = data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].size
}

data "huaweicloud_vpcv3_eips" "bandwidth_size_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  bandwidth_size = [data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].size]
}

output "bandwidth_size_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.bandwidth_size_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.bandwidth_size_filter.publicips[*].bandwidth.0.size : v == local.bandwidth_size]
  )
}

locals {
  bandwidth_share_type = data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].share_type
}

data "huaweicloud_vpcv3_eips" "bandwidth_share_type_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  bandwidth_share_type = [data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].share_type]
}

output "bandwidth_share_type_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.bandwidth_share_type_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.bandwidth_share_type_filter.publicips[*].bandwidth.0.share_type :
  v == local.bandwidth_share_type]
  )
}

locals {
  bandwidth_charge_mode = data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].charge_mode
}

data "huaweicloud_vpcv3_eips" "bandwidth_charge_mode_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  bandwidth_charge_mode = [data.huaweicloud_vpcv3_eips.test.publicips[0].bandwidth[0].charge_mode]
}

output "bandwidth_charge_mode_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.bandwidth_charge_mode_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.bandwidth_charge_mode_filter.publicips[*].bandwidth.0.charge_mode :
  v == local.bandwidth_charge_mode]
  )
}

locals {
  billing_info = data.huaweicloud_vpcv3_eips.test.publicips[0].billing_info
}

data "huaweicloud_vpcv3_eips" "billing_info_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  billing_info = [data.huaweicloud_vpcv3_eips.test.publicips[0].billing_info]
}

output "billing_info_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.billing_info_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.billing_info_filter.publicips[*].billing_info : v == local.billing_info]
  )
}

locals {
  billing_mode = "YEARLY_MONTHLY"
}

data "huaweicloud_vpcv3_eips" "billing_mode_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test
  ]

  billing_mode = "YEARLY_MONTHLY"
}

output "billing_mode_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.billing_mode_filter.publicips) > 0
}

locals {
  associate_instance_type = data.huaweicloud_vpcv3_eips.test.publicips[0].associate_instance_type
}

data "huaweicloud_vpcv3_eips" "associate_instance_type_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  associate_instance_type = [data.huaweicloud_vpcv3_eips.test.publicips[0].associate_instance_type]
}

output "associate_instance_type_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.associate_instance_type_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.associate_instance_type_filter.publicips[*].associate_instance_type :
  v == local.associate_instance_type]
  )
}

locals {
  associate_instance_id = data.huaweicloud_vpcv3_eips.test.publicips[0].associate_instance_id
}

data "huaweicloud_vpcv3_eips" "associate_instance_id_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  associate_instance_id = [data.huaweicloud_vpcv3_eips.test.publicips[0].associate_instance_id]
}

output "associate_instance_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.associate_instance_id_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.associate_instance_id_filter.publicips[*].associate_instance_id :
  v == local.associate_instance_id]
  )
}

locals {
  enterprise_project_id = data.huaweicloud_vpcv3_eips.test.publicips[0].enterprise_project_id
}

data "huaweicloud_vpcv3_eips" "enterprise_project_id_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  enterprise_project_id = [data.huaweicloud_vpcv3_eips.test.publicips[0].enterprise_project_id]
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.enterprise_project_id_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.enterprise_project_id_filter.publicips[*].enterprise_project_id :
  v == local.enterprise_project_id]
  )
}

locals {
  public_border_group = data.huaweicloud_vpcv3_eips.test.publicips[0].public_border_group
}

data "huaweicloud_vpcv3_eips" "public_border_group_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  public_border_group = [data.huaweicloud_vpcv3_eips.test.publicips[0].public_border_group]
}

output "public_border_group_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.public_border_group_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.public_border_group_filter.publicips[*].public_border_group :
  v == local.public_border_group]
  )
}

locals {
  allow_share_bandwidth_type_any = data.huaweicloud_vpcv3_eips.test.publicips[0].allow_share_bandwidth_types[0]
}

data "huaweicloud_vpcv3_eips" "allow_share_bandwidth_type_any_filter" {
  depends_on = [
    huaweicloud_dds_instance_eip_associate.test,
    data.huaweicloud_vpcv3_eips.test
  ]

  allow_share_bandwidth_type_any = [data.huaweicloud_vpcv3_eips.test.publicips[0].allow_share_bandwidth_types[0]]
}

output "allow_share_bandwidth_type_any_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_eips.allow_share_bandwidth_type_any_filter.publicips) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_eips.allow_share_bandwidth_type_any_filter.publicips[*].allow_share_bandwidth_types :
  contains(v, local.allow_share_bandwidth_type_any)]
  )
}
`, testDataSourceDataSourceEipVpcv3Eips_base(name), name)
}
