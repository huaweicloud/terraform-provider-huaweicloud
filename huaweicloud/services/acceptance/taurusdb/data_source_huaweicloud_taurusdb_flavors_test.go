package taurusdb

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccHuaweiCloudTaurusDBFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudTaurusDBFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBFlavorsDataSourceID("data.huaweicloud_taurusdb_flavors.flavor"),
				),
			},
		},
	})
}

func testAccCheckTaurusDBFlavorsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TaurusDB data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("the TaurusDB data source ID not set")
		}

		return nil
	}
}

var testAccHuaweiCloudTaurusDBFlavorsDataSource_basic = `
data "huaweicloud_taurusdb_flavors" "flavor" {
  engine = "gaussdb-mysql"
  version = "8.0"
  availability_zone_mode = "multi"
}
`
