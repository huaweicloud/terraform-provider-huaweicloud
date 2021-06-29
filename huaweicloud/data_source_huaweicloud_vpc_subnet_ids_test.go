package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcSubnetIdsV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_vpc_subnet_ids.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetIdV2DataSource_vpcsubnet(rName),
			},
			{
				Config: testAccSubnetIdV2DataSource_subnetids(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccSubnetIdV2DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ids.#", "1"),
				),
			},
		},
	})
}

func testAccSubnetIdV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find vpc subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Vpc Subnet data source ID not set")
		}

		return nil
	}
}

func testAccSubnetIdV2DataSource_vpcsubnet(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "172.16.8.0/24"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "172.16.8.0/24"
  gateway_ip = "172.16.8.1"
  vpc_id     = huaweicloud_vpc.test.id
}
`, rName, rName)
}

func testAccSubnetIdV2DataSource_subnetids(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnet_ids" "test" {
  vpc_id = huaweicloud_vpc.test.id
}
`, testAccSubnetIdV2DataSource_vpcsubnet(rName))
}
