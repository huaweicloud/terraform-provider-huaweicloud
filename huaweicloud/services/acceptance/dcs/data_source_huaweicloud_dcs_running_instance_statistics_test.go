package dcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsRunningInstanceStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_running_instance_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcsRunningInstanceStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.#"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.input_kbps"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.output_kbps"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.keys"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.used_memory"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.max_memory"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.cmd_get_count"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.cmd_set_count"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.used_cpu"),
				),
			},
		},
	})
}

func testDataSourceDcsRunningInstanceStatistics_basic() string {
	return `
data "huaweicloud_dcs_running_instance_statistics" "test" {}
`
}
