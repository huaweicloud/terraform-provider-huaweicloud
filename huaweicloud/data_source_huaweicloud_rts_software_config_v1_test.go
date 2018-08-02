package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRtsSoftwareConfigV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRtsSoftwareConfigV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsSoftwareConfigV1DataSourceID("data.huaweicloud_rts_software_config_v1.configs"),
					resource.TestCheckResourceAttr("data.huaweicloud_rts_software_config_v1.configs", "name", "huaweicloud-config"),
					resource.TestCheckResourceAttr("data.huaweicloud_rts_software_config_v1.configs", "group", "script"),
				),
			},
		},
	})
}

func testAccCheckRtsSoftwareConfigV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find software config data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("RTS software config data source ID not set ")
		}

		return nil
	}
}

var testAccRtsSoftwareConfigV1DataSource_basic = `
resource "huaweicloud_rts_software_config_v1" "config_1" {
  name = "huaweicloud-config"
  output_values = [{
    type = "String"
    name = "result"
    error_output = "false"
    description = "value1"
  }]
  input_values = [{
    default = "0"
    type = "String"
    name = "foo"
    description = "value2"
  }]
  group = "script"
}

data "huaweicloud_rts_software_config_v1" "configs" {
  id = "${huaweicloud_rts_software_config_v1.config_1.id}"
}
`
