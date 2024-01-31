package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/nat/v3/dnats"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getPrivateDnatRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 client: %s", err)
	}

	return dnats.Get(client, state.Primary.ID)
}

// The backend forwarding object is the ECS instance.
func TestAccPrivateDnatRule_basic(t *testing.T) {
	var (
		obj dnats.Rule

		rName = "huaweicloud_nat_private_dnat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateDnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnatRule_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gateway_id",
						"huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "tcp"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id",
						"huaweicloud_nat_private_transit_ip.test", "id"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "1000"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrPair(rName, "backend_interface_id",
						"huaweicloud_compute_instance.test", "network.0.port"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "2000"),
					resource.TestCheckResourceAttrSet(rName, "backend_type"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
				),
			},
			{
				Config: testAccPrivateDnatRule_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "udp"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "3000"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "4000"),
				),
			},
			{
				// Check the ports of internal service and transit service.
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPrivateDnatRule_basic_step_3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "any"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "0"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// The ports of internal service and transit service are both empty, ignore import check for them.
				ImportStateVerifyIgnore: []string{
					"internal_service_port",
					"transit_service_port",
				},
			},
		},
	})
}

func testAccPrivateDnatRule_transitIpConfig(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "transit_ip_used" {
  name = "%[1]s-transit-ip"
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "transit_ip_used" {
  vpc_id     = huaweicloud_vpc.transit_ip_used.id
  name       = "%[1]s-transit-ip"
  cidr       = cidrsubnet(huaweicloud_vpc.transit_ip_used.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.transit_ip_used.cidr, 4, 1), 1)
}

resource "huaweicloud_nat_private_transit_ip" "test" {
  subnet_id             = huaweicloud_vpc_subnet.transit_ip_used.id
  enterprise_project_id = "0"
}
`, name)
}

func testAccPrivateDnatRule_ecsPart(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  image_id          = data.huaweicloud_images_image.test.id
  security_groups   = [huaweicloud_networking_secgroup.test.name]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  admin_pass        = "%[3]s"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_nat_private_gateway" "test" {
  subnet_id             = huaweicloud_vpc_subnet.test.id
  name                  = "%[2]s"
  enterprise_project_id = "0"
}
`, common.TestBaseComputeResources(name), name, acceptance.RandomPassword("!@%-_=+[]:./?"))
}

func testAccPrivateDnatRule_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  protocol              = "tcp"
  description           = "Created by acc test"
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  transit_service_port  = 1000
  backend_interface_id  = huaweicloud_compute_instance.test.network[0].port
  internal_service_port = 2000
}
`, testAccPrivateDnatRule_ecsPart(name), testAccPrivateDnatRule_transitIpConfig(name))
}

func testAccPrivateDnatRule_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  protocol              = "udp"
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  transit_service_port  = 3000
  backend_interface_id  = huaweicloud_compute_instance.test.network[0].port
  internal_service_port = 4000
}
`, testAccPrivateDnatRule_ecsPart(name), testAccPrivateDnatRule_transitIpConfig(name))
}

func testAccPrivateDnatRule_basic_step_3(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id           = huaweicloud_nat_private_gateway.test.id
  protocol             = "any"
  transit_ip_id        = huaweicloud_nat_private_transit_ip.test.id
  backend_interface_id = huaweicloud_compute_instance.test.network[0].port
}
`, testAccPrivateDnatRule_ecsPart(name), testAccPrivateDnatRule_transitIpConfig(name))
}

// The backend forwarding object is the ELB loadbalancer.
func TestAccPrivateDnatRule_elbBackend(t *testing.T) {
	var (
		obj dnats.Rule

		rName = "huaweicloud_nat_private_dnat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateDnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnatRule_elbBackend_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gateway_id",
						"huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "tcp"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id",
						"huaweicloud_nat_private_transit_ip.test", "id"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "1000"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrPair(rName, "backend_interface_id",
						"data.huaweicloud_networking_port.test", "id"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "2000"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
				),
			},
			{
				Config: testAccPrivateDnatRule_elbBackend_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "udp"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "3000"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "4000"),
				),
			},
			{
				// Check the ports of internal service and transit service.
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPrivateDnatRule_elbBackend_step_3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "any"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "0"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// The ports of internal service and transit service are both empty, ignore import check for them.
				ImportStateVerifyIgnore: []string{
					"internal_service_port",
					"transit_service_port",
				},
			},
		},
	})
}

func testAccPrivateDnatRule_elbBackend_base(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_network_acl" "test" {
  name = "%[3]s"

  subnets = [
    huaweicloud_vpc_subnet.test.id
  ]

  inbound_rules = [
    huaweicloud_network_acl_rule.test.id
  ]
}

resource "huaweicloud_network_acl_rule" "test" {
  name                   = "%[3]s"
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = huaweicloud_vpc.test.cidr
  source_port            = "8080"
  destination_ip_address = "0.0.0.0/0"
  destination_port       = "8081"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_icmp_all" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_elb_member" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  ports             = "80,8081"
  remote_ip_prefix  = huaweicloud_vpc.test.cidr
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_all_group" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  remote_group_id   = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_networking_secgroup_rule" "out_v4_all" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[3]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_compute_eip_associate" "test" {
  public_ip   = huaweicloud_vpc_eip.test.address
  instance_id = huaweicloud_compute_instance.test.id
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[3]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_elb_listener" "test" {
  name            = "%[3]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id

  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "huaweicloud_elb_pool" "test" {
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_elb_listener.test.id

  persistence {
    type = "HTTP_COOKIE"
  }
}

resource "huaweicloud_elb_monitor" "test" {
  protocol    = "HTTP"
  interval    = 20
  timeout     = 15
  max_retries = 10
  url_path    = "/"
  port        = 8080
  pool_id     = huaweicloud_elb_pool.test.id
}

resource "huaweicloud_elb_member" "test" {
  address       = huaweicloud_compute_instance.test.access_ip_v4
  protocol_port = 8080
  pool_id       = huaweicloud_elb_pool.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.subnet_id
}

data "huaweicloud_networking_port" "test" {
  network_id = huaweicloud_vpc_subnet.test.id
  fixed_ip   = huaweicloud_elb_loadbalancer.test.ipv4_address
}
`, testAccPrivateDnatRule_ecsPart(name), testAccPrivateDnatRule_transitIpConfig(name), name)
}

func testAccPrivateDnatRule_elbBackend_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  protocol              = "tcp"
  description           = "Created by acc test"
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  transit_service_port  = 1000
  backend_interface_id  = data.huaweicloud_networking_port.test.id
  internal_service_port = 2000
}
`, testAccPrivateDnatRule_elbBackend_base(name))
}

func testAccPrivateDnatRule_elbBackend_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  protocol              = "udp"
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  transit_service_port  = 3000
  backend_interface_id  = data.huaweicloud_networking_port.test.id
  internal_service_port = 4000
}
`, testAccPrivateDnatRule_elbBackend_base(name))
}

func testAccPrivateDnatRule_elbBackend_step_3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id           = huaweicloud_nat_private_gateway.test.id
  protocol             = "any"
  transit_ip_id        = huaweicloud_nat_private_transit_ip.test.id
  backend_interface_id = data.huaweicloud_networking_port.test.id
}
`, testAccPrivateDnatRule_elbBackend_base(name))
}

// The backend forwarding object is the VIP.
func TestAccPrivateDnatRule_vipBackend(t *testing.T) {
	var (
		obj dnats.Rule

		rName = "huaweicloud_nat_private_dnat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateDnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnatRule_vipBackend_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gateway_id",
						"huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "tcp"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id",
						"huaweicloud_nat_private_transit_ip.test", "id"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "1000"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrPair(rName, "backend_interface_id",
						"huaweicloud_networking_vip.test", "id"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "2000"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
				),
			},
			{
				Config: testAccPrivateDnatRule_vipBackend_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "udp"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "3000"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "4000"),
				),
			},
			{
				// Check the ports of internal service and transit service.
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPrivateDnatRule_vipBackend_step_3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "any"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "0"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// The ports of internal service and transit service are both empty, ignore import check for them.
				ImportStateVerifyIgnore: []string{
					"internal_service_port",
					"transit_service_port",
				},
			},
		},
	})
}

func testAccPrivateDnatRule_vipBackend_base(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_nat_private_gateway" "test" {
  subnet_id             = huaweicloud_vpc_subnet.test.id
  name                  = "%[3]s"
  enterprise_project_id = "0"
}

resource "huaweicloud_networking_vip" "test" {
  network_id = huaweicloud_vpc_subnet.test.id
}
`, common.TestBaseNetwork(name), testAccPrivateDnatRule_transitIpConfig(name), name)
}

func testAccPrivateDnatRule_vipBackend_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  protocol              = "tcp"
  description           = "Created by acc test"
  transit_service_port  = 1000
  backend_interface_id  = huaweicloud_networking_vip.test.id
  internal_service_port = 2000
}

`, testAccPrivateDnatRule_vipBackend_base(name))
}

func testAccPrivateDnatRule_vipBackend_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  protocol              = "udp"
  transit_service_port  = 3000
  backend_interface_id  = huaweicloud_networking_vip.test.id
  internal_service_port = 4000
}
`, testAccPrivateDnatRule_vipBackend_base(name))
}

func testAccPrivateDnatRule_vipBackend_step_3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id           = huaweicloud_nat_private_gateway.test.id
  transit_ip_id        = huaweicloud_nat_private_transit_ip.test.id
  protocol             = "any"
  backend_interface_id = huaweicloud_networking_vip.test.id
}
`, testAccPrivateDnatRule_vipBackend_base(name))
}

func TestAccPrivateDnatRule_customIpAddress(t *testing.T) {
	var (
		obj dnats.Rule

		rName = "huaweicloud_nat_private_dnat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateDnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnatRule_customIpAddress_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gateway_id",
						"huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "any"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id",
						"huaweicloud_nat_private_transit_ip.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
				),
			},
			{
				Config: testAccPrivateDnatRule_customIpAddress_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "1000"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "backend_private_ip", "172.168.0.69"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "2000"),
				),
			},
			{
				Config: testAccPrivateDnatRule_customIpAddress_step_3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "udp"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "3000"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "backend_private_ip", "172.168.0.79"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "4000"),
				),
			},
			{
				// Check the ports of internal service and transit service.
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPrivateDnatRule_customIpAddress_step_4(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "any"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "0"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// The ports of internal service and transit service are both empty, ignore import check for them.
				ImportStateVerifyIgnore: []string{
					"internal_service_port",
					"transit_service_port",
				},
			},
		},
	})
}

func testAccPrivateDnatRule_customIpAddress_base(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_nat_private_gateway" "test" {
  subnet_id             = huaweicloud_vpc_subnet.test.id
  name                  = "%[3]s"
  enterprise_project_id = "0"
}

`, common.TestBaseNetwork(name), testAccPrivateDnatRule_transitIpConfig(name), name)
}

// Default protocol 'any'
func testAccPrivateDnatRule_customIpAddress_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id         = huaweicloud_nat_private_gateway.test.id
  transit_ip_id      = huaweicloud_nat_private_transit_ip.test.id
  backend_private_ip = "172.168.0.69"
}
`, testAccPrivateDnatRule_customIpAddress_base(name))
}

func testAccPrivateDnatRule_customIpAddress_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  protocol              = "tcp"
  description           = "Created by acc test"
  transit_service_port  = 1000
  backend_private_ip    = "172.168.0.69"
  internal_service_port = 2000
}
`, testAccPrivateDnatRule_customIpAddress_base(name))
}

func testAccPrivateDnatRule_customIpAddress_step_3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  protocol              = "udp"
  transit_service_port  = 3000
  backend_private_ip    = "172.168.0.79"
  internal_service_port = 4000
}
`, testAccPrivateDnatRule_customIpAddress_base(name))
}

func testAccPrivateDnatRule_customIpAddress_step_4(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id         = huaweicloud_nat_private_gateway.test.id
  transit_ip_id      = huaweicloud_nat_private_transit_ip.test.id
  protocol           = "any"
  backend_private_ip = "172.168.0.79"
}
`, testAccPrivateDnatRule_customIpAddress_base(name))
}
