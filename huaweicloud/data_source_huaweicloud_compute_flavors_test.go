package huaweicloud

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccEcsFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsFlavorDataSourceID("data.huaweicloud_compute_flavors.this"),
				),
			},
		},
	})
}

func testAccCheckEcsFlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find compute flavors data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Compute Flavors data source ID not set")
		}

		return nil
	}
}

const testAccEcsFlavorsDataSource_basic = `
data "huaweicloud_compute_flavors" "this" {
	performance_type = "normal"
	cpu_core_count   = 2
	memory_size      = 4
}
`
