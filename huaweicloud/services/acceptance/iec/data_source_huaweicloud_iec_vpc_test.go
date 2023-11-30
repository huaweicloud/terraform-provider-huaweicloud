package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-vpc-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccVpcDataSourceCheck("data.huaweicloud_iec_vpc.by_id", rName),
					testAccVpcDataSourceCheck("data.huaweicloud_iec_vpc.by_name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_iec_vpc.by_id", "mode", "SYSTEM"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_iec_vpc.by_name", "cidr", "192.168.0.0/16"),
				),
			},
		},
	})
}

func testAccVpcDataSourceCheck(n, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		vpcRs, ok := s.RootModule().Resources["huaweicloud_iec_vpc.test"]
		if !ok {
			return fmt.Errorf("can't find IEC vpc in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != vpcRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				vpcRs.Primary.Attributes["id"],
			)
		}

		if attr["name"] != rName {
			return fmt.Errorf("bad IEC vpc name %s", attr["name"])
		}

		return nil
	}
}

func testAccVpcDataSource_basic(rName string) string {
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
