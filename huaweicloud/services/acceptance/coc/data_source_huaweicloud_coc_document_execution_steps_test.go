package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDocumentExecutionSteps_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_document_execution_steps.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDocumentExecutionSteps_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.execution_step_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.inputs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.inputs.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.inputs.0.value"),
					resource.TestCheckOutput("execution_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDocumentExecutionSteps_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_document_execution_steps" "test" {
  execution_id = huaweicloud_coc_document_execute.test.id
}

output "execution_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_document_execution_steps.test.data) > 0
}
`, testCocDocumentExecute_basic(name))
}
