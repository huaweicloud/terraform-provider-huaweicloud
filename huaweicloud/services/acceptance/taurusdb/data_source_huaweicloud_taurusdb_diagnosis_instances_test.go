package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBDiagnosisInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_diagnosis_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBDiagnosisInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_infos.#"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBDiagnosisInstances_basic() string {
	return `
data "huaweicloud_taurusdb_diagnosis_statistics" "test" {}

data "huaweicloud_taurusdb_diagnosis_instances" "test" {
  metric_name = data.huaweicloud_taurusdb_diagnosis_statistics.test.diagnosis_info[0].metric_name
}
`
}
