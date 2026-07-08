package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBRedistributionParametersModify_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBRedistributionParametersModify_basic(),
			},
		},
	})
}

func testAccGaussDBRedistributionParametersModify_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_redistribution_parameters_modify" "test" {
  instance_id          = "%s"
  redis_parallel_jobs  = 4
  redis_resource_level = "l"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
