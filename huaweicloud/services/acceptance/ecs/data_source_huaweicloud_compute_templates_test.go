package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_templates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeTemplates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "launch_templates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_templates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_templates.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_templates.0.default_version"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_templates.0.latest_version"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_templates.0.created_at"),
					resource.TestCheckOutput("launch_template_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeTemplates_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_compute_templates" "test" {
  depends_on = [huaweicloud_compute_template.test]
}

data "huaweicloud_compute_templates" "launch_template_id_filter" {
  depends_on = [huaweicloud_compute_template.test]

  launch_template_id = [data.huaweicloud_compute_templates.test.launch_templates[0].id]
}
locals {
  launch_template_id = data.huaweicloud_compute_templates.test.launch_templates[0].id
}
output "launch_template_id_filter_is_useful" {
  value = length(data.huaweicloud_compute_templates.launch_template_id_filter.launch_templates) > 0 && alltrue(
  [for v in data.huaweicloud_compute_templates.launch_template_id_filter.launch_templates[*].id : v == local.launch_template_id]
  )
}

data "huaweicloud_compute_templates" "name_filter" {
  depends_on = [huaweicloud_compute_template.test]

  name = [data.huaweicloud_compute_templates.test.launch_templates[0].name]
}
locals {
  name = data.huaweicloud_compute_templates.test.launch_templates[0].name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_compute_templates.name_filter.launch_templates) > 0 && alltrue(
  [for v in data.huaweicloud_compute_templates.name_filter.launch_templates[*].name : v == local.name]
  )
}
`, testAccComputeTemplate_basic(name))
}
