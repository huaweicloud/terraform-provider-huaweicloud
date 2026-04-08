package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOfflineKeyAnalyses_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_offline_key_analyses.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOfflineKeyAnalyses_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.started_at"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.finished_at"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOfflineKeyAnalyses_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_offline_key_analyses" "test" {
  depends_on  = [huaweicloud_dcs_offline_key_analysis.test]

  instance_id = huaweicloud_dcs_instance.test.id
}

locals {
  status = data.huaweicloud_dcs_offline_key_analyses.test.records[0].status
}
data "huaweicloud_dcs_offline_key_analyses" "status_filter" {
  depends_on  = [huaweicloud_dcs_offline_key_analysis.test]

  instance_id = huaweicloud_dcs_instance.test.id
  status      = data.huaweicloud_dcs_offline_key_analyses.test.records[0].status
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_dcs_offline_key_analyses.status_filter.records) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_offline_key_analyses.status_filter.records[*].status : v == local.status]
  )
}
`, testOfflineKeyAnalysis_basic(name))
}
