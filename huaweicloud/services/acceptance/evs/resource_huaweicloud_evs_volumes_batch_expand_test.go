package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEvsVolumesBatchExpand_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test needs to create an EVS prepaid volume and an EVS postpaid volume before running.
			// The HW_EVS_VOLUME_NEW_SIZE and HW_EVS_PREPAID_VOLUME_NEW_SIZE needs bigger than now.
			acceptance.TestAccPreCheckEVSVolumeID(t)
			acceptance.TestAccPreCheckEVSVolumeNewSize(t)
			acceptance.TestAccPreCheckEVSPrepaidVolumeID(t)
			acceptance.TestAccPreCheckEVSPrepaidVolumeNewSize(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolumesBatchExpand_postpaid(),
			},
			{
				Config: testAccEvsVolumesBatchExpand_prepaid(),
			},
		},
	})
}

func testAccEvsVolumesBatchExpand_postpaid() string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volumes_batch_expand" "postpaid_test" {
  volumes  {
    id       = "%[1]s"
    new_size = "%[2]s"
  }

  is_auto_pay = true
}
`, acceptance.HW_EVS_VOLUME_ID, acceptance.HW_EVS_VOLUME_NEW_SIZE)
}

func testAccEvsVolumesBatchExpand_prepaid() string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volumes_batch_expand" "prepaid_test" {
  volumes  {
    id       = "%[1]s"
    new_size = "%[2]s"
  }

  is_auto_pay = true
}
`, acceptance.HW_EVS_PREPAID_VOLUME_ID, acceptance.HW_EVS_PREPAID_VOLUME_NEW_SIZE)
}
