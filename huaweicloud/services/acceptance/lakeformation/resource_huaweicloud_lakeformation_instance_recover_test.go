package lakeformation

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccLakeFormationInstanceRecover_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLakeFormationInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccLakeFormationInstanceRecover_basic(),
			},
		},
	})
}

func testAccLakeFormationInstanceRecover_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_lakeformation_instance_recover" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_LAKE_FORMATION_INSTANCE_ID)
}
