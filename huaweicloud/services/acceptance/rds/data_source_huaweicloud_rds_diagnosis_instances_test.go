package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsDiagnosisInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_diagnosis_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsDiagnosisInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
				),
			},
		},
	})
}

func testDataSourceRdsDiagnosisInstances_basic(_ string) string {
	return `
data "huaweicloud_rds_diagnosis_instances" "test" {
  engine    = "sqlserver"
  diagnosis = "insufficient_capacity"
}
`
}
