package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAiComponentDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_ai_component_detail.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAiComponentDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					// The query API has no response.
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
				),
			},
		},
	})
}

func testDataSourceAiComponentDetail_basic() string {
	return `
data "huaweicloud_hss_ai_component_detail" "test" {
  category         = "host"
  catalogue        = "app"
  first_scan_time  = "1774076220000"
  latest_scan_time = "1773730600349"
}
`
}
