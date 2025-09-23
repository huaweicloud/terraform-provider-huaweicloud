package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
)

func getPrivateDnatRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 client: %s", err)
	}

	return nat.GetPrivateDnatRule(client, state.Primary.ID)
}

// The backend forwarding object is the ECS instance.
func TestAccPrivateDnatRule_basic(t *testing.T) {
	var (
		obj   interface{}
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
					resource.TestCheckResourceAttrPair(rName, "gateway_id", "huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "tcp"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id", "huaweicloud_nat_private_transit_ip.test", "id"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "1000"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrPair(rName, "backend_interface_id", "huaweicloud_compute_instance.test", "network.0.port"),
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
		obj   interface{}
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
					resource.TestCheckResourceAttrPair(rName, "gateway_id", "huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "tcp"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id", "huaweicloud_nat_private_transit_ip.test", "id"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "1000"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrPair(rName, "backend_interface_id", "huaweicloud_elb_loadbalancer.test", "ipv4_port_id"),
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

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_nat_private_gateway" "test" {
  subnet_id             = huaweicloud_vpc_subnet.test.id
  name                  = "%[2]s"
  enterprise_project_id = "0"
}
`, common.TestBaseNetwork(name), name)
}

func testAccPrivateDnatRule_elbBackend_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  protocol              = "tcp"
  description           = "Created by acc test"
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  transit_service_port  = 1000
  backend_interface_id  = huaweicloud_elb_loadbalancer.test.ipv4_port_id
  internal_service_port = 2000
}
`, testAccPrivateDnatRule_elbBackend_base(name), testAccPrivateDnatRule_transitIpConfig(name))
}

func testAccPrivateDnatRule_elbBackend_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  protocol              = "udp"
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  transit_service_port  = 3000
  backend_interface_id  = huaweicloud_elb_loadbalancer.test.ipv4_port_id
  internal_service_port = 4000
}
`, testAccPrivateDnatRule_elbBackend_base(name), testAccPrivateDnatRule_transitIpConfig(name))
}

func testAccPrivateDnatRule_elbBackend_step_3(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id           = huaweicloud_nat_private_gateway.test.id
  protocol             = "any"
  transit_ip_id        = huaweicloud_nat_private_transit_ip.test.id
  backend_interface_id = huaweicloud_elb_loadbalancer.test.ipv4_port_id
}
`, testAccPrivateDnatRule_elbBackend_base(name), testAccPrivateDnatRule_transitIpConfig(name))
}

// The backend forwarding object is the VIP.
func TestAccPrivateDnatRule_vipBackend(t *testing.T) {
	var (
		obj   interface{}
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
					resource.TestCheckResourceAttrPair(rName, "gateway_id", "huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "tcp"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id", "huaweicloud_nat_private_transit_ip.test", "id"),
					resource.TestCheckResourceAttr(rName, "transit_service_port", "1000"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrPair(rName, "backend_interface_id", "huaweicloud_networking_vip.test", "id"),
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
		obj   interface{}
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
					resource.TestCheckResourceAttrPair(rName, "gateway_id", "huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "any"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id", "huaweicloud_nat_private_transit_ip.test", "id"),
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
