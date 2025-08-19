package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDocuments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_documents.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDocuments_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.document_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enterprise_project_id"),
					resource.TestCheckOutput("name_like_filter_is_useful", "true"),
					resource.TestCheckOutput("creator_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("document_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDocuments_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_documents" "test" {}

data "huaweicloud_coc_documents" "name_like_filter" {
  name_like = huaweicloud_coc_document.test.name
}

output "name_like_filter_is_useful" {
  value = length(data.huaweicloud_coc_documents.name_like_filter.data) > 0
}

data "huaweicloud_coc_documents" "creator_filter" {
  creator = huaweicloud_coc_document.test.creator
}

output "creator_filter_is_useful" {
  value = length(data.huaweicloud_coc_documents.creator_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_documents.creator_filter.data[*].creator : v == huaweicloud_coc_document.test.creator]
  )
}

data "huaweicloud_coc_documents" "enterprise_project_id_filter" {
  enterprise_project_id = huaweicloud_coc_document.test.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_documents.enterprise_project_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_documents.enterprise_project_id_filter.data[*].enterprise_project_id :
      v == huaweicloud_coc_document.test.enterprise_project_id]
  )
}

data "huaweicloud_coc_documents" "document_type_filter" {
  document_type = "PUBLIC"
}

output "document_type_filter_is_useful" {
  value = length(data.huaweicloud_coc_documents.document_type_filter.data) > 0
}
`, testDocument_basic(name))
}
