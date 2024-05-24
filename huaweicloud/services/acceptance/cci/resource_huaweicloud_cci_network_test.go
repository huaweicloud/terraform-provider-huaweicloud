package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cci/v1/networks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr"),
					resource.TestCheckResourceAttrPair(resourceName, "namespace",
						"huaweicloud_cci_namespace.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
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
			return "", fmt.Errorf("the namespace name (%s) or network ID (%s) is nil",
				rs.Primary.Attributes["namespace"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["namespace"], rs.Primary.ID), nil
	}
}

func testAccCciNetwork_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cci_namespace" "test" {
  name = "%s"
  type = "general-computing"
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccCciNetwork_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cci_network" "test" {
  name              = "%s"
  namespace         = huaweicloud_cci_namespace.test.name
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`, testAccCciNetwork_base(rName), rName)
}
