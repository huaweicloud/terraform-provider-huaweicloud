package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCloudLogPlatforms_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_cloud_log_platforms.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudLogPlatforms_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "platforms.#"),
					resource.TestCheckResourceAttrSet(dataSource, "platforms.0.tenant_managed_domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "platforms.0.dw_region"),
					resource.TestCheckResourceAttrSet(dataSource, "platforms.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "platforms.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "platforms.0.publish_status"),
					resource.TestCheckResourceAttrSet(dataSource, "platforms.0.white_list"),
				),
			},
		},
	})
}

func testAccDataSourceCloudLogPlatforms_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_cloud_log_platforms" "test" {
  workspace_id = "%[1]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
