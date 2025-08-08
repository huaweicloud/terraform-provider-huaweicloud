package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSfsTurboMountedIps_basic(t *testing.T) {
	resourceName := "data.huaweicloud_sfs_turbo_mounted_ips.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// create a sfs_turbo file system and ecs compute instance which mount the file system
			acceptance.TestAccPrecheckSFSTurboShareId(t)
			acceptance.TestAccPrecheckSFSTurboECSMoutShareIp(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSfsTurboMountedIps_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "ips_attribute.#"),
				),
			},
		},
	})
}

func testAccDataSourceSfsTurboMountedIps_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_sfs_turbo_mounted_ips" "test" {
  share_id = "%[1]s"
  ips      = "%[2]s"
}
`, acceptance.HW_SFS_TURBO_SHARE_ID, acceptance.HW_SFS_TURBO_ECS_MOUNT_SHARE_IP)
}
