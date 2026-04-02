package taurusdb

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBDehResourceDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBDehResourceDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBDehResourceDataSourceID("data.huaweicloud_taurusdb_dedicated_resource.test"),
				),
			},
		},
	})
}

func testAccCheckTaurusDBDehResourceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TaurusDB dedicated resource data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("the TaurusDB dedicated resource data source ID not set")
		}

		return nil
	}
}

const testAccTaurusDBDehResourceDataSource_basic = `
data "huaweicloud_taurusdb_dedicated_resource" "test" {
}
`
