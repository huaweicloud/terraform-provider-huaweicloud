package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDiagnosisTaskSummary_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_diagnosis_task_summary.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDiagnosisTaskSummary_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
					resource.TestCheckResourceAttrSet(dataSource, "start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_summary.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_summary.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_summary.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_summary.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_summary.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_summary.0.normal_item_num"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_summary.0.abnormal_item_num"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDiagnosisTaskSummary_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_diagnosis_task_summary" "test" {
  task_id = huaweicloud_coc_diagnosis_task.test.id
}
`, testCocDiagnosisTask_basic())
}
