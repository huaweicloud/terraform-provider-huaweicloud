package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDiagnosisTaskNodeDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_diagnosis_task_node_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDiagnosisTaskNodeDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "message"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnostic_task_node_id"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDiagnosisTaskNodeDetail_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_coc_diagnosis_task_node_detail" "test" {
  task_id     = huaweicloud_coc_diagnosis_task.test.id
  code        = "holmesInstall"
  instance_id = "%[2]s"
}
`, testCocDiagnosisTask_basic(), acceptance.HW_COC_INSTANCE_ID)
}
