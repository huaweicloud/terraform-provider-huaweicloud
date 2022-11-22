package vpn

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

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
					resource.TestCheckResourceAttr(rName, "connect_subnet", "192.168.1.0/24"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "availability_zones.0", "cn-north-4a"),
					resource.TestCheckResourceAttr(rName, "availability_zones.1", "cn-north-4b"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(rName, "master_eip.0.id", "huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "slave_eip.0.id", "huaweicloud_vpc_eip.test2", "id"),
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

func testGateway_base(name string) string {
	return fmt.Sprintf(`
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
  connect_subnet     = "192.168.1.0/24"
  availability_zones = ["cn-north-4a", "cn-north-4b"]

  master_eip {
    id = huaweicloud_vpc_eip.test1.id
  }

  slave_eip {
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
  connect_subnet     = "192.168.1.0/24"
  availability_zones = ["cn-north-4a", "cn-north-4b"]

  master_eip {
    id = huaweicloud_vpc_eip.test1.id
  }

  slave_eip {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name)
}
