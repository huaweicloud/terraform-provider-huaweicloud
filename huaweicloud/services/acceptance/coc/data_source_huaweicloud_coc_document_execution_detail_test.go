package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDocumentExecutionDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_document_execution_detail.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDocumentExecutionDetail_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "document_name"),
					resource.TestCheckResourceAttrSet(dataSource, "document_id"),
					resource.TestCheckResourceAttrSet(dataSource, "document_version_id"),
					resource.TestCheckResourceAttrSet(dataSource, "document_version"),
					resource.TestCheckResourceAttrSet(dataSource, "start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "creator"),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "parameters.#"),
					resource.TestCheckResourceAttrSet(dataSource, "parameters.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "parameters.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "type"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDocumentExecutionDetail_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_document_execution_detail" "test" {
  execution_id = huaweicloud_coc_document_execute.test.id
}
`, testCocDocumentExecute_basic(name))
}
