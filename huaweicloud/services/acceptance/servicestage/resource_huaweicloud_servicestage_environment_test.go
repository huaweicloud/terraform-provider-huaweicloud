package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/servicestage/v2/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getEnvResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ServiceStageV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage v2 client: %s", err)
	}
	return environments.Get(c, state.Primary.ID)
}

func TestAccEnvironment_basic(t *testing.T) {
	var (
		env          environments.Environment
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_servicestage_environment.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&env,
		getEnvResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironment_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "basic_resources.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "optional_resources.#", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "deploy_mode"),
				),
			},
			{
				Config: testAccEnvironment_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform test"),
					resource.TestCheckResourceAttr(resourceName, "basic_resources.#", "8"),
					resource.TestCheckResourceAttr(resourceName, "optional_resources.#", "8"),
					resource.TestCheckResourceAttrSet(resourceName, "deploy_mode"),
				),
			},
			{
				Config: testAccEnvironment_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "basic_resources.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "optional_resources.#", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "deploy_mode"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEnvironment_withEpsId(t *testing.T) {
	var (
		env          environments.Environment
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_servicestage_environment.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&env,
		getEnvResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironment_withEpsId(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEnvironment_base(rName string) string {
	return fmt.Sprintf(`
variable "subnet_config" {
  type = list(object({
    cidr       = string
    gateway_ip = string
  }))

  default = [
    {cidr = "192.168.192.0/18", gateway_ip = "192.168.192.1"},
    {cidr = "192.168.128.0/18", gateway_ip = "192.168.128.1"},
  ]
}

variable "rds_config" {
  type = list(object({
    fixed_ip = string
    port     = string
  }))

  default = [
    {fixed_ip = "192.168.0.58", port = "8636"},
    {fixed_ip = "192.168.0.160", port = "8637"},
  ]
}

variable "dcs_config" {
  type = list(object({
    port = number
  }))

  default = [
    {port = 6388},
    {port = 6389},
  ]
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

%s

%s`, rName, testAccEnvironment_baseRes(rName), testAccEnvironment_optioanlRes(rName))
}

func testAccEnvironment_baseRes(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "cce_bind" {
  count = 2

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    size        = 5
    name        = "%[1]s_${count.index}"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cce_cluster" "test" {
  count = 2

  name                   = "%[1]s-${count.index}"
  description            = "Created by terraform script and test for ServiceStage environment."
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  flavor_id              = "cce.s2.medium"
  container_network_type = "vpc-router"
  cluster_version        = "v1.19"
  cluster_type           = "VirtualMachine"
  eip                    = huaweicloud_vpc_eip.cce_bind[count.index].address

  kube_proxy_mode = "iptables"

  dynamic "masters" {
    for_each = slice(data.huaweicloud_availability_zones.test.names, 0, 3)

    content {
      availability_zone = masters.value
    }
  }
}

resource "huaweicloud_cce_node" "test" {
  count = 2

  cluster_id        = huaweicloud_cce_cluster.test[count.index].id
  name              = "%[1]s-${count.index}"
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_kps_keypair.test.name

  root_volume {
    volumetype = "SSD"
    size       = 100
  }

  data_volumes {
    volumetype = "SSD"
    size       = 100
  }

  lifecycle {
    ignore_changes = [
      tags,
    ]
  }  
}

resource "huaweicloud_cci_namespace" "test" {
  count = 2

  name                = "%[1]s-${count.index}"
  type                = "gpu-accelerated"
  auto_expend_enabled = true
}

resource "huaweicloud_vpc_subnet" "cci_bind" {
  count = 2

  name       = "%[1]s-${count.index}"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = var.subnet_config[count.index].cidr
  gateway_ip = var.subnet_config[count.index].gateway_ip
}

resource "huaweicloud_cci_network" "test" {
  count = 2

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[1]s-${count.index}"
  namespace         = huaweicloud_cci_namespace.test[count.index].name
  network_id        = huaweicloud_vpc_subnet.cci_bind[count.index].id
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name               = "%[1]s-${count.index}"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  key_pair           = huaweicloud_kps_keypair.test.name
  security_group_ids = [huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_as_configuration" "test" {
  scaling_configuration_name = "%[1]s"

  instance_config {
    flavor   = data.huaweicloud_compute_flavors.test.ids[0]
    image    = data.huaweicloud_images_image.test.id
    key_name = huaweicloud_kps_keypair.test.name

    disk {
      disk_type   = "SYS"
      volume_type = "GPSSD"
      size        = 40
    }
  }
}

resource "huaweicloud_as_group" "test" {
  count = 2

  scaling_group_name       = "%[1]s"
  scaling_configuration_id = huaweicloud_as_configuration.test.id
  vpc_id                   = huaweicloud_vpc.test.id

  max_instance_number    = 3
  min_instance_number    = 0
  desire_instance_number = 2

  delete_instances = "yes"
  delete_publicip  = true

  cool_down_time = 86400

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }

  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
}
`, rName)
}

func testAccEnvironment_optioanlRes(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl" "test" {
  name = "%[1]s"

  subnets = [
    huaweicloud_vpc_subnet.test.id,
  ]

  inbound_rules = [
    huaweicloud_network_acl_rule.test.id
  ]
}

resource "huaweicloud_network_acl_rule" "test" {
  name                   = "%[1]s"
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = huaweicloud_vpc.test.cidr
  source_port            = "8080"
  destination_ip_address = "0.0.0.0/0"
  destination_port       = "8081"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_elb_member" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  ports             = "80,8081"
  remote_ip_prefix  = huaweicloud_vpc.test.cidr
}

resource "huaweicloud_elb_loadbalancer" "test" {
  count = 2

  name            = "%[1]s_${count.index}"
  description     = "Created by terraform."
  vpc_id          = huaweicloud_vpc.test.id
  ipv4_subnet_id  = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id = huaweicloud_vpc_subnet.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_elb_listener" "test" {
  count = 2

  name            = "%[1]s_${count.index}"
  description     = "Created by terraform."
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_elb_loadbalancer.test[count.index].id

  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "huaweicloud_elb_pool" "test" {
  count = 2

  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_elb_listener.test[count.index].id

  persistence {
    type = "HTTP_COOKIE"
  }
}

resource "huaweicloud_elb_monitor" "test" {
  count = 2

  protocol    = "HTTP"
  interval    = 20
  timeout     = 15
  max_retries = 10
  url_path    = "/"
  port        = 8080
  pool_id     = huaweicloud_elb_pool.test[count.index].id
}

resource "huaweicloud_elb_member" "test" {
  count = 2

  address       = huaweicloud_compute_instance.test[count.index].access_ip_v4
  protocol_port = 8080
  pool_id       = huaweicloud_elb_pool.test[count.index].id
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_vpc_eip" "test" {
  count = 2

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s_${count.index}"
    size        = 10
    charge_mode = "traffic"
  }
}

resource "huaweicloud_rds_instance" "test" {
  count = 2

  name              = "%[1]s_${count.index}"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         =  huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = var.rds_config[count.index].fixed_ip

  db {
    password = "Huawei##123"
    type     = "PostgreSQL"
    version  = "12"
    port     = var.rds_config[count.index].port
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "huaweicloud_dcs_instance" "test" {
  count = 2

  name               = "%[1]s_${count.index}"
  engine_version     = "5.0"
  password           = "Huawei##123"
  engine             = "Redis"
  port               = var.dcs_config[count.index].port
  capacity           = 0.125
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  maintain_begin     = "22:00:00"
  maintain_end       = "02:00:00"

  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }
}
`, rName)
}

func testAccEnvironment_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_servicestage_environment" "test" {
  name        = "%s"
  description = "Created by terraform test"
  vpc_id      = huaweicloud_vpc.test.id

  basic_resources {
    type = "cce"
    id   = huaweicloud_cce_cluster.test[0].id
  }
  basic_resources {
    type = "cci"
    id   = huaweicloud_cci_namespace.test[0].name
  }
  basic_resources {
    type = "ecs"
    id   = huaweicloud_compute_instance.test[0].id
  }
  basic_resources {
    type = "as"
    id   = huaweicloud_as_group.test[0].id
  }

  optional_resources {
    type = "elb"
    id   = huaweicloud_elb_loadbalancer.test[0].id
  }
  optional_resources {
    type = "eip"
    id   = huaweicloud_vpc_eip.test[0].id
  }
  optional_resources {
    type = "rds"
    id   = huaweicloud_rds_instance.test[0].id
  }
  optional_resources {
    type = "dcs"
    id   = huaweicloud_dcs_instance.test[0].id
  }
}
`, testAccEnvironment_base(rName), rName)
}

func testAccEnvironment_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_servicestage_environment" "test" {
  name        = "%s-update"
  description = "Updated by terraform test"
  vpc_id      = huaweicloud_vpc.test.id

  dynamic "basic_resources" {
    for_each = huaweicloud_cce_cluster.test[*].id
    content {
      type = "cce"
      id   = basic_resources.value
    }
  }
  dynamic "basic_resources" {
    for_each = huaweicloud_cci_namespace.test[*].name
    content {
      type = "cci"
      id   = basic_resources.value
    }
  }
  dynamic "basic_resources" {
    for_each = huaweicloud_compute_instance.test[*].id
    content {
      type = "ecs"
      id   = basic_resources.value
    }
  }
  dynamic "basic_resources" {
    for_each = huaweicloud_as_group.test[*].id
    content {
      type = "as"
      id   = basic_resources.value
    }
  }

  dynamic "optional_resources" {
    for_each = huaweicloud_elb_loadbalancer.test[*].id
    content {
      type = "elb"
      id   = optional_resources.value
    }
  }
  dynamic "optional_resources" {
    for_each = huaweicloud_vpc_eip.test[*].id
    content {
      type = "eip"
      id   = optional_resources.value
    }
  }
  dynamic "optional_resources" {
    for_each = huaweicloud_rds_instance.test[*].id
    content {
      type = "rds"
      id   = optional_resources.value
    }
  }
  dynamic "optional_resources" {
    for_each = huaweicloud_dcs_instance.test[*].id
    content {
      type = "dcs"
      id   = optional_resources.value
    }
  }
}
`, testAccEnvironment_base(rName), rName)
}

func testAccEnvironment_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_servicestage_environment" "test" {
  name                  = "%s"
  vpc_id                = huaweicloud_vpc.test.id
  enterprise_project_id = "%s"

  dynamic "basic_resources" {
    for_each = huaweicloud_cce_cluster.test[*].id
    content {
      type = "cce"
      id   = basic_resources.value
    }
  }
}
`, testAccEnvironment_base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
