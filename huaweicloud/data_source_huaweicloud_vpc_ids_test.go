package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcIdsV1DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcIdsV1Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccVpcIdsV2DataSourceID("data.huaweicloud_vpc_ids.test"),
				),
			},
		},
	})
}

func testAccVpcIdsV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find vpc data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Vpc data source ID not set")
		}

		return nil
	}
}

func testAccDataSourceVpcIdsV1Config(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "172.16.9.0/24"
}

data "huaweicloud_vpc_ids" "test" {
}
`, rName)
}
