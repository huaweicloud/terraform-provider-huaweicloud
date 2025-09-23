package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceLoadBalancers_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_loadbalancers.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceLoadBalancers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.#"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.name"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.description"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.availability_zone.#"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.billing_info"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.charge_mode"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.deletion_protection_enable"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.listeners.#"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.pools.#"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.log_group_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.log_topic_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.operating_status"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.public_border_group"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.publicips.#"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.publicips.0.publicip_address"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.publicips.0.publicip_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.publicips.0.ip_version"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.waf_failure_action"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.ipv4_address"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.ipv4_port_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.l4_flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.l7_flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("availability_zone_filter_is_useful", "true"),
					resource.TestCheckOutput("billing_info_filterr_is_useful", "true"),
					resource.TestCheckOutput("deletion_protection_enable_filter_is_useful", "true"),
					resource.TestCheckOutput("log_group_id_filter_is_useful", "true"),
					resource.TestCheckOutput("log_topic_id_filter_is_useful", "true"),
					resource.TestCheckOutput("member_address_filter_is_useful", "true"),
					resource.TestCheckOutput("member_device_id_filter_is_useful", "true"),
					resource.TestCheckOutput("operating_status_filter_is_useful", "true"),
					resource.TestCheckOutput("protection_status_filter_is_useful", "true"),
					resource.TestCheckOutput("provisioning_status_filter_is_useful", "true"),
					resource.TestCheckOutput("publicips_publicip_id_filter_is_useful", "true"),
					resource.TestCheckOutput("publicips_publicip_address_filter_is_useful", "true"),
					resource.TestCheckOutput("publicips_ip_version_filter_is_useful", "true"),
					resource.TestCheckOutput("ipv4_address_filter_is_useful", "true"),
					resource.TestCheckOutput("ipv4_port_id_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ipv4_subnet_id_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("l4_flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("l7_flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_is_useful", "true"),
				),
			},
		},
	})
}

func TestAccDatasourceLoadBalancers_gateway(t *testing.T) {
	rName := "data.huaweicloud_elb_loadbalancers.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbGatewayType(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceLoadBalancers_gateway(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.#"),
					resource.TestCheckResourceAttr(rName, "loadbalancers.0.name", name),
					resource.TestCheckResourceAttr(rName, "loadbalancers.0.loadbalancer_type", "gateway"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.gw_flavor_id"),
				),
			},
		},
	})
}

func testAccDatasourceLoadBalancers_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 22.04 server 64bit"
  most_recent = true
}

data "huaweicloud_elb_flavors" "l4flavors" {
  type            = "L4"
  max_connections = 1000000
  cps             = 20000
  bandwidth       = 100
}

data "huaweicloud_elb_flavors" "l7flavors" {
  type            = "L7"
  max_connections = 400000
  cps             = 4000
  bandwidth       = 100
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  description    = "update flavors"
  l4_flavor_id   = data.huaweicloud_elb_flavors.l4flavors.ids[0]
  l7_flavor_id   = data.huaweicloud_elb_flavors.l7flavors.ids[0]
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%[2]s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8083
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = false
}

resource "huaweicloud_elb_pool" "test" {
  name        = "%[2]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_elb_listener.test.id
}

resource "huaweicloud_elb_member" "test" {
  name          = "%[2]s"
  address       = huaweicloud_compute_instance.test.access_ip_v4
  weight        = 2
  protocol_port = 8080
  pool_id       = huaweicloud_elb_pool.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_elb_logtank" "test" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  log_group_id    = huaweicloud_lts_group.test.id
  log_topic_id    = huaweicloud_lts_stream.test.id
}

resource "huaweicloud_vpc_eipv3_associate" "test" {
  publicip_id             = huaweicloud_vpc_eip.test.id
  associate_instance_type = "ELB"
  associate_instance_id   = huaweicloud_elb_loadbalancer.test.id
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccDatasourceLoadBalancers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_loadbalancers" "test" {
  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}

data "huaweicloud_elb_loadbalancers" "name_filter" {
  name = "%[2]s"

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.name_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.name_filter.loadbalancers[*].name :v == "%[2]s"]
  )  
}

data "huaweicloud_elb_loadbalancers" "availability_zone_filter" {
  availability_zone = tolist(huaweicloud_elb_loadbalancer.test.availability_zone)[0]

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  availability_zone = tolist(huaweicloud_elb_loadbalancer.test.availability_zone)[0]
}
output "availability_zone_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.availability_zone_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.availability_zone_filter.loadbalancers[*].availability_zone :
  contains(v, local.availability_zone)]
  )  
}

data "huaweicloud_elb_loadbalancers" "billing_info_filter" {
  billing_info = [for v in data.huaweicloud_elb_loadbalancers.test.loadbalancers[*].billing_info : v if v != ""][0]

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  billing_info = [for v in data.huaweicloud_elb_loadbalancers.test.loadbalancers[*].billing_info : v if v != ""][0]
}
output "billing_info_filterr_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.billing_info_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.billing_info_filter.loadbalancers[*].billing_info :v == local.billing_info]
  )  
}

data "huaweicloud_elb_loadbalancers" "deletion_protection_enable_filter" {
  deletion_protection_enable = data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].deletion_protection_enable

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test,
    data.huaweicloud_elb_loadbalancers.test
  ]
}
locals {
  deletion_protection_enable = data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].deletion_protection_enable
}
output "deletion_protection_enable_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.deletion_protection_enable_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.deletion_protection_enable_filter.loadbalancers[*].deletion_protection_enable :
  v == local.deletion_protection_enable]
  )  
}

data "huaweicloud_elb_loadbalancers" "log_group_id_filter" {
  log_group_id = huaweicloud_lts_group.test.id

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  log_group_id = huaweicloud_lts_group.test.id
}
output "log_group_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.log_group_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.log_group_id_filter.loadbalancers[*].log_group_id : v == local.log_group_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "log_topic_id_filter" {
  log_topic_id = huaweicloud_lts_stream.test.id

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  log_topic_id = huaweicloud_lts_stream.test.id
}
output "log_topic_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.log_topic_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.log_topic_id_filter.loadbalancers[*].log_topic_id : v == local.log_topic_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "member_address_filter" {
  member_address = huaweicloud_compute_instance.test.access_ip_v4

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
output "member_address_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.member_address_filter.loadbalancers) > 0  
}

data "huaweicloud_elb_loadbalancers" "member_device_id_filter" {
  member_device_id = huaweicloud_compute_instance.test.id

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
output "member_device_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.member_device_id_filter.loadbalancers) > 0
}

data "huaweicloud_elb_loadbalancers" "operating_status_filter" {
  operating_status = huaweicloud_elb_loadbalancer.test.operating_status

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  operating_status = huaweicloud_elb_loadbalancer.test.operating_status
}
output "operating_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.operating_status_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.operating_status_filter.loadbalancers[*].operating_status :
  v == local.operating_status]
  )  
}

data "huaweicloud_elb_loadbalancers" "protection_status_filter" {
  protection_status = huaweicloud_elb_loadbalancer.test.protection_status

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  protection_status = huaweicloud_elb_loadbalancer.test.protection_status
}
output "protection_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.protection_status_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.protection_status_filter.loadbalancers[*].protection_status :
  v == local.protection_status]
  )  
}

data "huaweicloud_elb_loadbalancers" "provisioning_status_filter" {
  provisioning_status = data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].provisioning_status

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    data.huaweicloud_elb_loadbalancers.test
  ]
}
locals {
  provisioning_status = data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].provisioning_status
}
output "provisioning_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.provisioning_status_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.provisioning_status_filter.loadbalancers[*].provisioning_status :
  v == local.provisioning_status]
  )  
}

locals {
  publicip = [for v in data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].publicips : v if length(v) > 0][0]
}
data "huaweicloud_elb_loadbalancers" "publicips_publicip_id_filter" {
  publicips = [format("publicip_id=%%s", local.publicip.publicip_id)]

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    data.huaweicloud_elb_loadbalancers.test
  ]
}
output "publicips_publicip_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.publicips_publicip_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.publicips_publicip_id_filter.loadbalancers[*].publicips[0].publicip_id :
  v == local.publicip.publicip_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "publicips_publicip_address_filter" {
  publicips = [format("publicip_address=%%s", local.publicip.publicip_address)]

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    data.huaweicloud_elb_loadbalancers.test
  ]
}
output "publicips_publicip_address_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.publicips_publicip_address_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.publicips_publicip_address_filter.loadbalancers[*].publicips[0].publicip_address :
  v == local.publicip.publicip_address]
  )  
}

data "huaweicloud_elb_loadbalancers" "publicips_ip_version_filter" {
  publicips = [format("ip_version=%%s", local.publicip.ip_version)]

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    data.huaweicloud_elb_loadbalancers.test
  ]
}
output "publicips_ip_version_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.publicips_ip_version_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.publicips_ip_version_filter.loadbalancers[*].publicips[0].ip_version :
  v == local.publicip.ip_version]
  )  
}

data "huaweicloud_elb_loadbalancers" "ipv4_address_filter" {
  ipv4_address = huaweicloud_elb_loadbalancer.test.ipv4_address

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  ipv4_address = huaweicloud_elb_loadbalancer.test.ipv4_address
}
output "ipv4_address_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.ipv4_address_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.ipv4_address_filter.loadbalancers[*].ipv4_address : v == local.ipv4_address]
  )  
}

data "huaweicloud_elb_loadbalancers" "ipv4_port_id_filter" {
  ipv4_port_id = huaweicloud_elb_loadbalancer.test.ipv4_port_id

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  ipv4_port_id = huaweicloud_elb_loadbalancer.test.ipv4_port_id
}
output "ipv4_port_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.ipv4_port_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.ipv4_port_id_filter.loadbalancers[*].ipv4_port_id :
  v == local.ipv4_port_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "vpc_id_filter" {
  vpc_id = huaweicloud_elb_loadbalancer.test.vpc_id

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  vpc_id = huaweicloud_elb_loadbalancer.test.vpc_id
}
output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.vpc_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.vpc_id_filter.loadbalancers[*].vpc_id : v == local.vpc_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "ipv4_subnet_id_filter" {
  ipv4_subnet_id = huaweicloud_elb_loadbalancer.test.ipv4_subnet_id

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  ipv4_subnet_id = huaweicloud_elb_loadbalancer.test.ipv4_subnet_id
}
output "ipv4_subnet_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.vpc_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.ipv4_subnet_id_filter.loadbalancers[*].ipv4_subnet_id : 
  v == local.ipv4_subnet_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "description_filter" {
  description = huaweicloud_elb_loadbalancer.test.description

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  description = huaweicloud_elb_loadbalancer.test.description
}
output "description_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.description_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.description_filter.loadbalancers[*].description : 
  v == local.description]
  )  
}

data "huaweicloud_elb_loadbalancers" "l4_flavor_id_filter" {
  l4_flavor_id = huaweicloud_elb_loadbalancer.test.l4_flavor_id

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  l4_flavor_id = huaweicloud_elb_loadbalancer.test.l4_flavor_id
}
output "l4_flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.l4_flavor_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.l4_flavor_id_filter.loadbalancers[*].l4_flavor_id : 
  v == local.l4_flavor_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "l7_flavor_id_filter" {
  l7_flavor_id = huaweicloud_elb_loadbalancer.test.l7_flavor_id

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  l7_flavor_id = huaweicloud_elb_loadbalancer.test.l7_flavor_id
}
output "l7_flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.l7_flavor_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.l7_flavor_id_filter.loadbalancers[*].l7_flavor_id : 
  v == local.l7_flavor_id]
  ) 
}

data "huaweicloud_elb_loadbalancers" "type_filter" {
  type = "dedicated"

  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  type = "dedicated"
}
output "type_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.type_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.type_filter.loadbalancers[*].type : 
  v == local.type]
  ) 
}

data "huaweicloud_elb_loadbalancers" "enterprise_project_id_filter" {
  depends_on = [
    huaweicloud_elb_logtank.test,
    huaweicloud_elb_member.test,
    huaweicloud_vpc_eipv3_associate.test
  ]
}
locals {
  enterprise_project_id = huaweicloud_elb_loadbalancer.test.enterprise_project_id
}
output "enterprise_project_id_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.enterprise_project_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.enterprise_project_id_filter.loadbalancers[*].enterprise_project_id : 
  v == local.enterprise_project_id]
  ) 
}
`, testAccDatasourceLoadBalancers_base(name), name)
}

func testAccDatasourceLoadBalancers_gateway(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%[1]s"
  vpc_id      = huaweicloud_vpc.test.id
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  ipv6_enable = true
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[1]s"
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id   = huaweicloud_vpc_subnet.test.id
  loadbalancer_type = "gateway"
  description       = "test gateway description"
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

data "huaweicloud_elb_loadbalancers" "test" {
  depends_on = [huaweicloud_elb_loadbalancer.test]
  name       = "%[1]s"
}
`, name)
}
