package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2Deployments_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_deployments.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2Deployments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "deployments.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployments.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployments.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployments.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployments.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployments.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployments.0.uid"),
				),
			},
		},
	})
}

func testAccDataSourceV2Deployments_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cciv2_deployments" "test" {
  depends_on = [huaweicloud_cciv2_deployment.test]

  namespace = huaweicloud_cciv2_namespace.test.name
}
`, testAccV2Deployment_basic(rName))
}
