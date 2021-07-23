package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpcEipDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEipConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcEipDataSourceID("data.huaweicloud_vpc_eip.by_public_ip"),
				),
			},
		},
	})
}

func testAccCheckVpcEipDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find eip data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Eip data source ID not set")
		}

		return nil
	}
}

func testAccDataSourceVpcEipConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

data "huaweicloud_vpc_eip" "by_public_ip" {
  public_ip = huaweicloud_vpc_eip.test.address
}
`, rName)
}
