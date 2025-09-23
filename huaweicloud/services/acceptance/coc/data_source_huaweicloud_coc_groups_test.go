package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_groups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.vendor"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.application_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sync_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sync_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sync_rules.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sync_rules.0.rule_tags"),
					resource.TestCheckOutput("component_id_filter_is_useful", "true"),
					resource.TestCheckOutput("id_list_filter_is_useful", "true"),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_like_filter_is_useful", "true"),
					resource.TestCheckOutput("code_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_groups" "test" {
  component_id = huaweicloud_coc_component.test.id

  depends_on = [huaweicloud_coc_group.test]
}

output "component_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_groups.test.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_groups.test.data[*].component_id : v == huaweicloud_coc_component.test.id]
  )
}

data "huaweicloud_coc_groups" "id_list_filter" {
  component_id = huaweicloud_coc_component.test.id
  id_list      = [huaweicloud_coc_group.test.id]
}

output "id_list_filter_is_useful" {
  value = length(data.huaweicloud_coc_groups.id_list_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_groups.id_list_filter.data[*].id : v == huaweicloud_coc_group.test.id]
  )
}

data "huaweicloud_coc_groups" "application_id_filter" {
  component_id   = huaweicloud_coc_component.test.id
  application_id = huaweicloud_coc_application.test.id

  depends_on = [huaweicloud_coc_group.test]
}

output "application_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_groups.application_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_groups.application_id_filter.data[*].application_id :
      v == huaweicloud_coc_application.test.id]
  )
}

data "huaweicloud_coc_groups" "name_like_filter" {
  component_id = huaweicloud_coc_component.test.id
  name_like    = huaweicloud_coc_group.test.name
}

output "name_like_filter_is_useful" {
  value = length(data.huaweicloud_coc_groups.name_like_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_groups.name_like_filter.data[*].name : v == huaweicloud_coc_group.test.name]
  )
}

data "huaweicloud_coc_groups" "code_filter" {
  component_id = huaweicloud_coc_component.test.id
  code         = huaweicloud_coc_group.test.code
}

output "code_filter_is_useful" {
  value = length(data.huaweicloud_coc_groups.code_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_groups.code_filter.data[*].code : v == huaweicloud_coc_group.test.code]
  )
}
`, testAccGroup_basic(name))
}
