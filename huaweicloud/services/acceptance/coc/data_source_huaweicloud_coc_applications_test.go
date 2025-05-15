package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocApplications_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_applications.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocApplicationID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocApplications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.path"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckOutput("id_list_filter_useful", "true"),
					resource.TestCheckOutput("parent_id_filter_useful", "true"),
					resource.TestCheckOutput("code_filter_useful", "true"),
					resource.TestCheckOutput("name_like_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocApplications_basic() string {
	return `
data "huaweicloud_coc_applications" "test" {}

locals {
  ids = data.huaweicloud_coc_applications.test.data[*].id
}

data "huaweicloud_coc_applications" "id_list_filter" {
    id_list = local.ids
}

locals {
  id_list_filter_result = [
    for v in data.huaweicloud_coc_applications.id_list_filter.data[*].id : v if contains(local.ids, v)
  ]
}
output "id_list_filter_useful" {
  value = join(",", sort(local.id_list_filter_result)) == join(",", sort(local.ids))
}

locals {
  parent_id = [for v in data.huaweicloud_coc_applications.test.data[*].parent_id : v if v != ""][0]
}

data "huaweicloud_coc_applications" "parent_id_filter" {
  parent_id = local.parent_id
}

output "parent_id_filter_useful" {
  value = length(data.huaweicloud_coc_applications.parent_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_applications.parent_id_filter.data[*].parent_id : v == local.parent_id]
  )
}

locals {
  code = [for v in data.huaweicloud_coc_applications.test.data[*].code : v if v != ""][0]
}

data "huaweicloud_coc_applications" "code_filter" {
  code = local.code
}

output "code_filter_useful" {
  value = data.huaweicloud_coc_applications.code_filter.data[0].code == local.code
}

locals {
  name_like = split("_", [for v in data.huaweicloud_coc_applications.test.data[*].name : v if v != ""][0])[0]
}

data "huaweicloud_coc_applications" "name_like_filter" {
  name_like = local.name_like
}

output "name_like_filter_is_useful" {
  value = length(data.huaweicloud_coc_applications.name_like_filter.data) > 0
}
`
}
