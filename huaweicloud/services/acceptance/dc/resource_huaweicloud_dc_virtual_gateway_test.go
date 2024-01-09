package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dc/v3/gateways"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVirtualGatewayFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DcV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DC v3 client: %s", err)
	}

	return gateways.Get(client, state.Primary.ID)
}

func TestAccVirtualGateway_basic(t *testing.T) {
	var (
		gateway gateways.VirtualGateway

		rName      = "huaweicloud_dc_virtual_gateway.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		cidr       = acceptance.RandomCidr()
		updateCidr = acceptance.RandomCidr()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&gateway,
		getVirtualGatewayFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualGateway_basic(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "local_ep_group.0", cidr),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrSet(rName, "asn"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVirtualGateway_update(updateName, updateCidr),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "local_ep_group.0", updateCidr),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_update"),
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

func testAccVirtualGateway_basic(name, cidr string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "%[2]s"
}

resource "huaweicloud_dc_virtual_gateway" "test" {
  vpc_id      = huaweicloud_vpc.test.id
  name        = "%[1]s"
  description = "Created by acc test"

  local_ep_group = [
    huaweicloud_vpc.test.cidr,
  ]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, cidr)
}

func testAccVirtualGateway_update(name, cidr string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "%[2]s"
}

resource "huaweicloud_dc_virtual_gateway" "test" {
  vpc_id = huaweicloud_vpc.test.id
  name   = "%[1]s"

  local_ep_group = [
    huaweicloud_vpc.test.cidr,
  ]

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, name, cidr)
}
