package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocDocumentExecutionOperation_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocDocumentExecutionOperation_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocDocumentExecutionOperation_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_coc_document_execution_operation" "test" {
  execution_id = huaweicloud_coc_document_execute.test.id
  operate_type = "TERMINATE"
}
`, testCocDocumentExecute_basic(name))
}
