package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocApplications_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_applications.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocApplications_basic(rName),
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
					resource.TestCheckOutput("id_list_filter_is_useful", "true"),
					resource.TestCheckOutput("parent_id_filter_is_useful", "true"),
					resource.TestCheckOutput("code_filter_is_useful", "true"),
					resource.TestCheckOutput("name_like_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocApplications_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_applications" "test" {
  id_list = [huaweicloud_coc_application.test.id]
}

output "id_list_filter_is_useful" {
  value = length(data.huaweicloud_coc_applications.test.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_applications.test.data[*].id : v == huaweicloud_coc_application.test.id]
  )
}

data "huaweicloud_coc_applications" "parent_id_filter" {
  parent_id = huaweicloud_coc_application.test.parent_id
}

output "parent_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_applications.parent_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_applications.parent_id_filter.data[*].parent_id :
      v == huaweicloud_coc_application.test.parent_id]
  )
}

data "huaweicloud_coc_applications" "code_filter" {
  code = huaweicloud_coc_application.test.code
}

output "code_filter_is_useful" {
  value = data.huaweicloud_coc_applications.code_filter.data[0].code == huaweicloud_coc_application.test.code
}

data "huaweicloud_coc_applications" "name_like_filter" {
  name_like = huaweicloud_coc_application.test.name
}

output "name_like_filter_is_useful" {
  value = length(data.huaweicloud_coc_applications.name_like_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_applications.name_like_filter.data[*].name :
      v == huaweicloud_coc_application.test.name]
  )
}
`, testAccApplication_parent_id(name))
}
