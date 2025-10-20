package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRecycleBinPolicy_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_evs_recycle_bin_policy.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRecycleBinPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "switch"),
					resource.TestCheckResourceAttrSet(dataSourceName, "threshold_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "keep_time"),
				),
			},
		},
	})
}

func testDataSourceRecycleBinPolicy_basic() string {
	return `
data "huaweicloud_evs_recycle_bin_policy" "test" {}
`
}
