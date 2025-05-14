package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocScripts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_scripts.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocScripts_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.gmt_created"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.gmt_modified"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.script_uuid"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.usage_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.0.risk_level"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.0.version"),
					resource.TestCheckOutput("name_like_filter_is_useful", "true"),
					resource.TestCheckOutput("creator_filter_is_useful", "true"),
					resource.TestCheckOutput("risk_level_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocScripts_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_script" "test" {
  name        = "%[1]s"
  description = "a new demo script"
  risk_level  = "MEDIUM"
  version     = "1.0.1"
  type        = "SHELL"
 
  content = <<EOF
#! /bin/bash
echo "hello $${name}!"
EOF

  parameters {
    name        = "name"
    value       = "world"
    description = "the parameter"
  }
}

data "huaweicloud_coc_scripts" "test" {
  depends_on = [huaweicloud_coc_script.test]
}

data "huaweicloud_coc_scripts" "name_like_filter" {
  name_like = "%[1]s"

  depends_on = [huaweicloud_coc_script.test]
}

output "name_like_filter_is_useful" {
  value = length(data.huaweicloud_coc_scripts.name_like_filter.data) > 0
}

locals {
  creator = [for v in data.huaweicloud_coc_scripts.test.data[*].creator : v if v != ""][0]
}

data "huaweicloud_coc_scripts" "creator_filter" {
  creator = local.creator

  depends_on = [huaweicloud_coc_script.test]
}

output "creator_filter_is_useful" {
  value = length(data.huaweicloud_coc_scripts.creator_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scripts.creator_filter.data[*].creator : v == local.creator]
  )
}

data "huaweicloud_coc_scripts" "risk_level_filter" {
  risk_level = huaweicloud_coc_script.test.risk_level

  depends_on = [huaweicloud_coc_script.test]
}

output "risk_level_filter_is_useful" {
  value = length(data.huaweicloud_coc_scripts.risk_level_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scripts.risk_level_filter.data[*].properties[0].risk_level : v == huaweicloud_coc_script.test.risk_level]
  )
}

data "huaweicloud_coc_scripts" "type_filter" {
  type = huaweicloud_coc_script.test.type

  depends_on = [huaweicloud_coc_script.test]
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_coc_scripts.type_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scripts.type_filter.data[*].type : v == huaweicloud_coc_script.test.type]
  )
}

locals {
  enterprise_project_id = [for v in data.huaweicloud_coc_scripts.test.data[*].enterprise_project_id : v if v != ""][0]
}

data "huaweicloud_coc_scripts" "enterprise_project_id_filter" {
  enterprise_project_id = local.enterprise_project_id

  depends_on = [huaweicloud_coc_script.test]
}

locals {
  enterprise_project_id_filter_result = data.huaweicloud_coc_scripts.enterprise_project_id_filter.data[0].enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = local.enterprise_project_id_filter_result == local.enterprise_project_id
}
`, rName)
}
