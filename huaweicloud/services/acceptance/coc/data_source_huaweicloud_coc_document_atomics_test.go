package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDocumentAtomics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_document_atomics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDocumentAtomics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.atomic_unique_key"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.atomic_name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.atomic_name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.tags.#"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDocumentAtomics_basic() string {
	return `
data "huaweicloud_coc_document_atomics" "test" {}
`
}
