package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAiComponentStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_ai_component_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAiComponentStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
				),
			},
		},
	})
}

func testDataSourceAiComponentStatistics_basic() string {
	return `
data "huaweicloud_hss_ai_component_statistics" "test" {
  enterprise_project_id = "all_granted_eps"
  category              = "host"
  catalogue             = "app"
}
`
}
