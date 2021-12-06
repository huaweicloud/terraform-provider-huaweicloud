package dms

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDmsAZV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsAZV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsAZV1DataSourceID("data.huaweicloud_dms_az.az1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_az.az1", "code", "cn-north-4a"),
				),
			},
		},
	})
}

func testAccCheckDmsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find Dms az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Dms az data source ID not set")
		}

		return nil
	}
}

var testAccDmsAZV1DataSource_basic = fmt.Sprintf(`
data "huaweicloud_dms_az" "az1" {
  code = "cn-north-4a"
}
`)
