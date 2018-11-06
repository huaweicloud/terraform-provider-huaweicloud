package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDmsAZV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDmsAZV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsAZV1DataSourceID("data.huaweicloud_dms_az_v1.az1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_az_v1.az1", "name", "可用区1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_az_v1.az1", "port", "8002"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_az_v1.az1", "code", "cn-north-1a"),
				),
			},
		},
	})
}

func testAccCheckDmsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dms az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dms az data source ID not set")
		}

		return nil
	}
}

var testAccDmsAZV1DataSource_basic = fmt.Sprintf(`
data "huaweicloud_dms_az_v1" "az1" {
name = "可用区1"
port = "8002"
code = "cn-north-1a"
}
`)
