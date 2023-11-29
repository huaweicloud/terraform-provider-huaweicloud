package vpn

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGatewayResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGateway: Query the VPN gateway detail
	var (
		getGatewayHttpUrl = "v5/{project_id}/vpn-gateways/{id}"
		getGatewayProduct = "vpn"
	)
	getGatewayClient, err := config.NewServiceClient(getGatewayProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Gateway Client: %s", err)
	}

	getGatewayPath := getGatewayClient.Endpoint + getGatewayHttpUrl
	getGatewayPath = strings.ReplaceAll(getGatewayPath, "{project_id}", getGatewayClient.ProjectID)
	getGatewayPath = strings.ReplaceAll(getGatewayPath, "{id}", state.Primary.ID)

	getGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getGatewayResp, err := getGatewayClient.Request("GET", getGatewayPath, &getGatewayOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Gateway: %s", err)
	}
	return utils.FlattenResponse(getGatewayResp)
}

func TestAccGateway_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ha_mode", "active-active"),
					resource.TestCheckResourceAttrPair(rName, "connect_subnet", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(rName, "eip1.0.id", "huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "eip2.0.id", "huaweicloud_vpc_eip.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
				),
			},
			{
				Config: testGateway_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttr(rName, "local_subnets.1", "192.168.2.0/24"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGateway_activeStandbyHAMode(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_activeStandbyHAMode(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ha_mode", "active-standby"),
					resource.TestCheckResourceAttrPair(rName, "connect_subnet", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(rName, "eip1.0.id", "huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "eip2.0.id", "huaweicloud_vpc_eip.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGateway_deprecated(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_deprecated(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ha_mode", "active-standby"),
					resource.TestCheckResourceAttrPair(rName, "connect_subnet", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(rName, "master_eip.0.id", "huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "slave_eip.0.id", "huaweicloud_vpc_eip.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
				),
			},
			{
				Config: testGateway_deprecated_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttr(rName, "local_subnets.1", "192.168.2.0/24"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGateway_withER(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_withER(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "network_type", "private"),
					resource.TestCheckResourceAttr(rName, "attachment_type", "er"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "access_private_ip_1", "172.16.0.99"),
					resource.TestCheckResourceAttr(rName, "access_private_ip_2", "172.16.0.100"),
					resource.TestCheckResourceAttrPair(rName, "er_id", "huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "access_vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "access_subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGateway_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = "professional1"
  attachment_type = "vpc"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_vpc_eip" "test1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-1"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip" "test2" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-2"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, name)
}

func testGateway_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s"
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test1.id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name)
}

func testGateway_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s-update"
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr, "192.168.2.0/24"]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test1.id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name)
}

func testGateway_activeStandbyHAMode(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s"
  vpc_id             = huaweicloud_vpc.test.id
  ha_mode            = "active-standby"
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test1.id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name)
}

func testGateway_withER(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "172.16.0.0/24"
  gateway_ip = "172.16.0.1"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[3]
  ]

  name = "%[1]s"
  asn  = "65000"
}

data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = "professional1"
  attachment_type = "er"
}

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%[1]s"
  network_type       = "private"
  attachment_type    = "er"
  er_id              = huaweicloud_er_instance.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  access_vpc_id    = huaweicloud_vpc.test.id
  access_subnet_id = huaweicloud_vpc_subnet.test.id
  
  access_private_ip_1 = "172.16.0.99"
  access_private_ip_2 = "172.16.0.100"
}
`, name)
}

func testGateway_deprecated(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s"
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  master_eip {
    id = huaweicloud_vpc_eip.test1.id
  }

  slave_eip {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name)
}

func testGateway_deprecated_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s-update"
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr, "192.168.2.0/24"]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  master_eip {
    id = huaweicloud_vpc_eip.test1.id
  }

  slave_eip {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name)
}
