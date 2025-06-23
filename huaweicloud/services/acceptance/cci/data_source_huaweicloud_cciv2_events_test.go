package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2Events_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_events.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2Events_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "events.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "events.0.metadata.name"),
				),
			},
		},
	})
}

func testAccDataSourceV2Events_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cciv2_events" "test" {
  depends_on = [huaweicloud_cciv2_deployment.test]

  namespace = huaweicloud_cciv2_namespace.test.name
}
`, testAccV2Deployment_basic(rName))
}
