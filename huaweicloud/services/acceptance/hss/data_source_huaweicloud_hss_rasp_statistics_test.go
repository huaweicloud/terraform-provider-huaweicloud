package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRaspStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_rasp_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRaspStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "protect_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "anti_tampering_num"),
				),
			},
		},
	})
}

func testDataSourceRaspStatistics_basic() string {
	return `
data "huaweicloud_hss_rasp_statistics" "test" {
  enterprise_project_id = "0"
}
`
}
