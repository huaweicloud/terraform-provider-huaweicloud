package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBDiagnosisStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_diagnosis_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBDiagnosisStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_info.0.metric_name"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_info.0.count"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBDiagnosisStatistics_basic() string {
	return `data "huaweicloud_taurusdb_diagnosis_statistics" "test" {}`
}
