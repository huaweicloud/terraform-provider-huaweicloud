package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2Secrets_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_secrets.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2Secrets_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.data.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.type"),
				),
			},
		},
	})
}

func testAccDataSourceV2Secrets_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cciv2_secrets" "test" {
  depends_on = [huaweicloud_cciv2_secret.test]

  namespace = huaweicloud_cciv2_namespace.test.name
}
`, testAccV2Secret_basic(rName))
}
