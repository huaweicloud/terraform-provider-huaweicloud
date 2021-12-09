package cci

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cci/v1/networks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getNetworkResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CciV1BetaClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud CCI Beta v1 client: %s", err)
	}
	return networks.Get(c, state.Primary.Attributes["namespace"], state.Primary.ID).Extract()
}

func TestAccCciNetwork_basic(t *testing.T) {
	var network networks.Network
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cci_network.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&network,
		getNetworkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCciNetwork_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", acceptance.HW_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "namespace",
						"${huaweicloud_cci_namespace.test.name}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "network_id",
						"${huaweicloud_vpc_subnet.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "subnet_id",
						"${huaweicloud_vpc_subnet.test.subnet_id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "security_group_id",
						"${huaweicloud_networking_secgroup.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "vpc_id",
						"${huaweicloud_vpc.test.id}"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCciNetworkImportStateFunc(resourceName),
			},
		},
	})
}

func testAccCciNetworkImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["namespace"] == "" {
			return "", fmt.Errorf("The namespace name (%s) or network ID (%s) is nil.",
				rs.Primary.Attributes["namespace"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["namespace"], rs.Primary.ID), nil
	}
}

func testAccCciNetwork_base(rName string) string {
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()

	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "%s"
  gateway_ip = "%s"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_cci_namespace" "test" {
  name = "%s"
  type = "gpu-accelerated"
}
`, rName, randCidr, rName, randCidr, randGatewayIp, rName, rName)
}

func testAccCciNetwork_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cci_network" "test" {
  name              = "%s"
  availability_zone = "%s"
  namespace         = huaweicloud_cci_namespace.test.name
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`, testAccCciNetwork_base(rName), rName, acceptance.HW_AVAILABILITY_ZONE)
}
