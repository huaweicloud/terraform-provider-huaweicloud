package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2PersistentVolumeClaims_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_persistent_volumes.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2PersistentVolumeClaims_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "pvcs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pvcs.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pvcs.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pvcs.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pvcs.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pvcs.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pvcs.0.status.#"),
				),
			},
		},
	})
}

func testAccDataSourceV2PersistentVolumeClaims_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cciv2_persistent_volume_claims" "test" {
  depends_on = [huaweicloud_cciv2_persistent_volume_claim.test]

  namespace = huaweicloud_cciv2_namespace.test.name
}
`, testAccV2PersistentVolumeClaim_basic(rName), rName)
}
