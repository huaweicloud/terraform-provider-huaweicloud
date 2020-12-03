package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHuaweiCloudNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudNetworkingSecGroupV2DataSource_group(rName),
			},
			{
				Config: testAccHuaweiCloudNetworkingSecGroupV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.huaweicloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_secgroup_v2.secgroup_1", "name", rName),
				),
			},
		},
	})
}

func TestAccHuaweiCloudNetworkingSecGroupV2DataSource_secGroupID(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudNetworkingSecGroupV2DataSource_group(rName),
			},
			{
				Config: testAccHuaweiCloudNetworkingSecGroupV2DataSource_secGroupID(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.huaweicloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_secgroup_v2.secgroup_1", "name", rName),
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

func testAccHuaweiCloudNetworkingSecGroupV2DataSource_group(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
  name        = "%s"
  description = "My neutron security group"
}
`, rName)
}

func testAccHuaweiCloudNetworkingSecGroupV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup_v2" "secgroup_1" {
  name = "${huaweicloud_networking_secgroup_v2.secgroup_1.name}"
}
`, testAccHuaweiCloudNetworkingSecGroupV2DataSource_group(rName))
}

func testAccHuaweiCloudNetworkingSecGroupV2DataSource_secGroupID(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup_v2" "secgroup_1" {
  name = "${huaweicloud_networking_secgroup_v2.secgroup_1.name}"
}
`, testAccHuaweiCloudNetworkingSecGroupV2DataSource_group(rName))
}
