package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDocumentAtomicDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_document_atomic_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDocumentAtomicDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "atomic_name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "atomic_name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.0.param_key"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.0.param_name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.0.param_name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.0.required"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.0.param_type"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.0.min"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.0.max"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.0.min_len"),
					resource.TestCheckResourceAttrSet(dataSource, "inputs.0.max_len"),
					resource.TestCheckResourceAttrSet(dataSource, "outputs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "outputs.0.supported"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDocumentAtomicDetail_basic() string {
	return `
data "huaweicloud_coc_document_atomic_detail" "test" {
  atomic_unique_key = "coc_step_sleep"
}
`
}
