package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsBigkeyAnalyses_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_bigkey_analyses.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDcsBigkeyAnalyses_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.scan_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.started_at"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.finished_at"),

					resource.TestCheckOutput("analysis_id_filter_is_useful", "true"),
					resource.TestCheckOutput("scan_type_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDcsBigkeyAnalyses_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_bigkey_analyses" "test" {
  depends_on  = [huaweicloud_dcs_bigkey_analysis.test]
  instance_id = huaweicloud_dcs_instance.instance_1.id
}

locals {
  analysis_id = huaweicloud_dcs_bigkey_analysis.test.id
}
data "huaweicloud_dcs_bigkey_analyses" "analysis_id_filter" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  analysis_id = huaweicloud_dcs_bigkey_analysis.test.id
}
output "analysis_id_filter_is_useful" {
  value = length(data.huaweicloud_dcs_bigkey_analyses.analysis_id_filter.records) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_bigkey_analyses.analysis_id_filter.records[*].id : v == local.analysis_id]
  )
}

locals {
  scan_type = huaweicloud_dcs_bigkey_analysis.test.scan_type
}
data "huaweicloud_dcs_bigkey_analyses" "scan_type_filter" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  scan_type   = huaweicloud_dcs_bigkey_analysis.test.scan_type
}
output "scan_type_filter_is_useful" {
  value = length(data.huaweicloud_dcs_bigkey_analyses.scan_type_filter.records) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_bigkey_analyses.scan_type_filter.records[*].scan_type : v == local.scan_type]
  )
}

locals {
  status = huaweicloud_dcs_bigkey_analysis.test.status
}
data "huaweicloud_dcs_bigkey_analyses" "status_filter" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  status      = huaweicloud_dcs_bigkey_analysis.test.status
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_dcs_bigkey_analyses.status_filter.records) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_bigkey_analyses.status_filter.records[*].status : v == local.status]
  )
}
`, testBigKeyAnalysis_basic(name))
}
