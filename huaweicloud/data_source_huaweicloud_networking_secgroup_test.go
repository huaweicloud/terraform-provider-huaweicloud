package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_networking_secgroup.secgroup_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV2DataSource_group(rName),
			},
			{
				Config: testAccNetworkingSecGroupV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "rules.#"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupV2DataSource_byID(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_networking_secgroup.secgroup_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV2DataSource_group(rName),
			},
			{
				Config: testAccNetworkingSecGroupV2DataSource_secGroupID(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "rules.#"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

func testAccNetworkingSecGroupV2DataSource_group(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "My neutron security group"
}
`, rName)
}

func testAccNetworkingSecGroupV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup" "secgroup_1" {
  name = huaweicloud_networking_secgroup.secgroup_1.name
}
`, testAccNetworkingSecGroupV2DataSource_group(rName))
}

func testAccNetworkingSecGroupV2DataSource_secGroupID(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup" "secgroup_1" {
  secgroup_id = huaweicloud_networking_secgroup.secgroup_1.id
}
`, testAccNetworkingSecGroupV2DataSource_group(rName))
}
