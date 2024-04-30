package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/nat/v2/gateways"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
)

func getPublicGatewayResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatGatewayClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v2 client: %s", err)
	}

	return gateways.Get(client, state.Primary.ID)
}

func TestAccPublicGateway_basic(t *testing.T) {
	var (
		obj gateways.Gateway

		rName         = "huaweicloud_nat_gateway.test"
		name          = acceptance.RandomAccResourceNameWithDash()
		updateName    = acceptance.RandomAccResourceNameWithDash()
		relatedConfig = common.TestBaseNetwork(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicGateway_basic_step_1(name, relatedConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "spec", string(nat.PublicSpecTypeSmall)),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "ngport_ip_address", "192.168.0.101"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
				),
			},
			{
				Config: testAccPublicGateway_basic_step_2(updateName, relatedConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "spec", string(nat.PublicSpecTypeMedium)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "ngport_ip_address", "192.168.0.101"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baaar"),
					resource.TestCheckResourceAttr(rName, "tags.newkey", "value"),
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

func testAccPublicGateway_basic_step_1(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "1"
  description           = "Created by acc test"
  ngport_ip_address     = "192.168.0.101"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, relatedConfig, name)
}

func testAccPublicGateway_basic_step_2(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "2"
  ngport_ip_address     = "192.168.0.101"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"

  tags = {
    foo    = "baaar"
    newkey = "value"
  }
}
`, relatedConfig, name)
}
