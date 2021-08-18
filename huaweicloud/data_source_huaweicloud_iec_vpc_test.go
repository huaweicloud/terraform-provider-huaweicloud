package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIECVpcDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-vpc-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIECVpc_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceIECVpcCheck("data.huaweicloud_iec_vpc.by_id", rName),
					testAccDataSourceIECVpcCheck("data.huaweicloud_iec_vpc.by_name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_iec_vpc.by_id", "mode", "SYSTEM"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_iec_vpc.by_name", "cidr", "192.168.0.0/16"),
				),
			},
		},
	})
}

func testAccDataSourceIECVpcCheck(n, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("root module has no resource called %s", n)
		}

		vpcRs, ok := s.RootModule().Resources["huaweicloud_iec_vpc.test"]
		if !ok {
			return fmtp.Errorf("can't find huaweicloud_iec_vpc.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != vpcRs.Primary.Attributes["id"] {
			return fmtp.Errorf(
				"id is %s; want %s",
				attr["id"],
				vpcRs.Primary.Attributes["id"],
			)
		}

		if attr["name"] != rName {
			return fmtp.Errorf("bad iec vpc name %s", attr["name"])
		}

		return nil
	}
}

func testAccDataSourceIECVpc_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

data "huaweicloud_iec_vpc" "by_id" {
  id = huaweicloud_iec_vpc.test.id
}

data "huaweicloud_iec_vpc" "by_name" {
  name = huaweicloud_iec_vpc.test.name
}
`, rName)
}
