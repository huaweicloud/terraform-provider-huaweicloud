package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDirUsage_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_sfs_turbo_dir_usage.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, please prepare a SFS Turbo file system which the type is **STANDARD**.
			acceptance.TestAccPrecheckSFSTurboShareId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDirUsage_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "dir_usage.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dir_usage.0.used_capacity"),
				),
			},
		},
	})
}

func testDataSourceDirUsage_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sfs_turbo_dir" "test" {
  share_id = "%[1]s"
  path     = "/temp"
}

data "huaweicloud_sfs_turbo_dir_usage" "test" {
  share_id = "%[1]s"
  path     = huaweicloud_sfs_turbo_dir.test.path
}
`, acceptance.HW_SFS_TURBO_SHARE_ID)
}
