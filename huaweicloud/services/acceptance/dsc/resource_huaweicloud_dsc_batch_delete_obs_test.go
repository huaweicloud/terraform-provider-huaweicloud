package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceBatchDeleteObs_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscObsId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBatchDeleteObs_basic(),
			},
		},
	})
}

func testAccResourceBatchDeleteObs_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_batch_delete_obs" "test" {
  obs_ids = split(",", "%s")
}
`, acceptance.HW_DSC_OBS_ID)
}
