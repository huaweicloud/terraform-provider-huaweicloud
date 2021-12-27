package deprecated

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDcsProductV1DataSource_basic(t *testing.T) {
	flavor := "redis.cluster.xu1.large.r2.4"
	resourceName := "data.huaweicloud_dcs_product.product1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsProductV1DataSource_basic(flavor),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsProductV1DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spec_code", flavor),
					resource.TestCheckResourceAttrSet(resourceName, "engine"),
					resource.TestCheckResourceAttrSet(resourceName, "engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_mode"),
				),
			},
		},
	})
}

func testAccCheckDcsProductV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DCS product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DCS product data source ID not set")
		}

		return nil
	}
}

func testAccDcsProductV1DataSource_basic(flavor string) string {
	return fmt.Sprintf(`
data "huaweicloud_dcs_product" "product1" {
  spec_code = "%s"
}
`, flavor)
}
