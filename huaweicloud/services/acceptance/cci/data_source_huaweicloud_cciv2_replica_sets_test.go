package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2ReplicaSets_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_replica_sets.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2ReplicaSets_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "replica_sets.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replica_sets.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replica_sets.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replica_sets.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replica_sets.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replica_sets.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replica_sets.0.status.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replica_sets.0.template.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replica_sets.0.selector.#"),
				),
			},
		},
	})
}

func testAccDataSourceV2ReplicaSets_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cciv2_replica_sets" "test" {
  depends_on = [huaweicloud_cciv2_deployment.test]

  name = "%[2]s"
}
`, testAccV2Deployment_basic(rName), rName)
}
