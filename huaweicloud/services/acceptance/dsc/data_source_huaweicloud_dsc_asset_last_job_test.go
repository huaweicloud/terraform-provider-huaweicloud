package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscAssetLastJob_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_asset_last_job.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscAssetId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDscAssetLastJob_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testDataSourceDscAssetLastJob_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_asset_last_job" "test" {
  asset_id = "%[1]s"
}
`, acceptance.HW_DSC_ASSET_ID)
}
