package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
)

func getPrivateTransitSubnetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 client: %s", err)
	}

	return nat.GetTransitSubnet(client, state.Primary.ID)
}

func TestAccPrivateTransitSubnet_basic(t *testing.T) {
	var (
		obj        interface{}
		rName      = "huaweicloud_nat_private_transit_subnet.test"
		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateTransitSubnetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateTransitSubnet_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "virsubnet_id", "data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(rName, "virsubnet_project_id", acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccPrivateTransitSubnet_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "virsubnet_id", "data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(rName, "virsubnet_project_id", acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test modified."),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baaar"),
					resource.TestCheckResourceAttr(rName, "tags.newkey", "value"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
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

func testAccPrivateTransitSubnet_base() string {
	return `
data "huaweicloud_vpc" "myvpc" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}
`
}

func testAccPrivateTransitSubnet_basic(name string) string {
	return fmt.Sprintf(`
%[3]s

resource "huaweicloud_nat_private_transit_subnet" "test" {
  name                  = "%[1]s"
  description           = "Created by acc test"
  virsubnet_id          = data.huaweicloud_vpc_subnet.test.id
  virsubnet_project_id  = "%[2]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_PROJECT_ID, testAccPrivateTransitSubnet_base())
}

func testAccPrivateTransitSubnet_update(name string) string {
	return fmt.Sprintf(`
%[3]s

resource "huaweicloud_nat_private_transit_subnet" "test" {
  name                 = "%[1]s"
  description          = "Created by acc test modified."
  virsubnet_id         = data.huaweicloud_vpc_subnet.test.id
  virsubnet_project_id = "%[2]s"

  tags = {
    foo    = "baaar"
    newkey = "value"
  }
}
`, name, acceptance.HW_PROJECT_ID, testAccPrivateTransitSubnet_base())
}
