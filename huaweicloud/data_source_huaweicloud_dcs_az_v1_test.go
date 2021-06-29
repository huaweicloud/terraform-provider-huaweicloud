package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDcsAZV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsAZV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsAZV1DataSourceID("data.huaweicloud_dcs_az_v1.az1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dcs_az_v1.az1", "port", "8403"),
				),
			},
		},
	})
}

func testAccCheckDcsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find Dcs az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Dcs az data source ID not set")
		}

		return nil
	}
}

var testAccDcsAZV1DataSource_basic = fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dcs_az_v1" "az1" {
  code = data.huaweicloud_availability_zones.test.names[0]
  port = "8403"
}
`)
