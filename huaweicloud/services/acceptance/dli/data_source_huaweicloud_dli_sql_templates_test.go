package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDliSqlTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dli_sql_templates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDliSqlTemplates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.group"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.owner"),

					resource.TestCheckOutput("template_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("group_filter_is_useful", "true"),
					resource.TestCheckOutput("owner_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDliSqlTemplates_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dli_sql_templates" "test" {
  depends_on = [
    huaweicloud_dli_sql_template.test
  ]
}

data "huaweicloud_dli_sql_templates" "template_id_filter" {
  template_id = local.template_id
}

locals {
  template_id = data.huaweicloud_dli_sql_templates.test.templates[0].id
}

output "template_id_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_templates.template_id_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_templates.template_id_filter.templates[*].id : v == local.template_id]
  )
}

data "huaweicloud_dli_sql_templates" "name_filter" {
  name = local.name
}

locals {
  name = data.huaweicloud_dli_sql_templates.test.templates[0].name
}
 
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_templates.name_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_templates.name_filter.templates[*].name : v == local.name]
  )
}

data "huaweicloud_dli_sql_templates" "group_filter" {
  group = local.group
}

locals {
  group = data.huaweicloud_dli_sql_templates.test.templates[0].group
}
 
output "group_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_templates.group_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_templates.group_filter.templates[*].group : v == local.group]
  )
}

data "huaweicloud_dli_sql_templates" "owner_filter" {
  owner = local.owner
}

locals {
  owner = data.huaweicloud_dli_sql_templates.test.templates[0].owner
}

output "owner_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_templates.owner_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_templates.owner_filter.templates[*].owner : v == local.owner]
  )
}
`, testSQLTemplate_basic(name))
}
