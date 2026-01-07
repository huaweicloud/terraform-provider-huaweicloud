package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCloudLogResources_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_cloud_log_resources.test"
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
				Config: testAccDataSourceCloudLogResources_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.0.alert"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.0.allow_alert"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.0.allow_lts"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.0.enable"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.0.update_time"),
				),
			},
		},
	})
}

func testAccDataSourceCloudLogResources_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_cloud_log_resources" "test" {
  workspace_id = "%s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
