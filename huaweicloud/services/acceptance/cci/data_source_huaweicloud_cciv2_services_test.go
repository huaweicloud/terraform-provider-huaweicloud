package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2Services_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_services.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2Services_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "services.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "services.0.name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.finalizers.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.status.#"),
				),
			},
		},
	})
}

func testAccDataSourceV2Services_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cciv2_services" "test" {
  depends_on = [huaweicloud_cciv2_service.test]

  namespace = huaweicloud_cciv2_namespace.test.name
}
`, testAccV2Service_basic(rName), rName)
}
