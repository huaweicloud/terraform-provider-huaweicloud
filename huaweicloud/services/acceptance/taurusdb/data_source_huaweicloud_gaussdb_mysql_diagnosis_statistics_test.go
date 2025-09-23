package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlDiagnosisStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_diagnosis_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlDiagnosisStatistics_basic(),
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

func testDataSourceGaussdbMysqlDiagnosisStatistics_basic() string {
	return `data "huaweicloud_gaussdb_mysql_diagnosis_statistics" "test" {}`
}
