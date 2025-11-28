package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwrImageManualSync_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrRepository(t)
			acceptance.TestAccPreCheckSwrOrigination(t)
			acceptance.TestAccPreCheckSwrTargetRegion(t)
			acceptance.TestAccPreCheckSwrTargetOrigination(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrImageManualSync_basic(),
			},
		},
	})
}

func testAccSwrImageManualSync_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_swr_image_tags" "test" {
  organization = "%[1]s"
  repository   = "%[2]s"
}

resource "huaweicloud_swr_image_manual_sync" "test" {
  organization        = "%[1]s"
  repository          = "%[2]s"
  image_tag           = [data.huaweicloud_swr_image_tags.test.image_tags[0].name]
  target_region       = "%[3]s"
  target_organization = "%[4]s"
  override            = true
}
`, acceptance.HW_SWR_ORGANIZATION, acceptance.HW_SWR_REPOSITORY, acceptance.HW_SWR_TARGET_REGION, acceptance.HW_SWR_TARGET_ORGANIZATION)
}
