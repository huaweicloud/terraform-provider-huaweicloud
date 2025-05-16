package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocScripts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_scripts.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocScriptID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocScripts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.operator"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.gmt_created"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.gmt_modified"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.script_uuid"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.usage_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.0.risk_level"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enterprise_project_id"),
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

func testDataSourceDataSourceCocScripts_basic() string {
	return `
data "huaweicloud_coc_scripts" "test" {}

locals {
  name_like = split("_", [for v in data.huaweicloud_coc_scripts.test.data[*].name : v if v != ""][0])[0]
}

data "huaweicloud_coc_scripts" "name_like_filter" {
  name_like = local.name_like
}

output "name_like_filter_is_useful" {
  value = length(data.huaweicloud_coc_scripts.name_like_filter.data) > 0
}

locals {
  creator = [for v in data.huaweicloud_coc_scripts.test.data[*].creator : v if v != ""][0]
}

data "huaweicloud_coc_scripts" "creator_filter" {
  creator = local.creator
}

output "creator_filter_is_useful" {
  value = length(data.huaweicloud_coc_scripts.creator_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scripts.creator_filter.data[*].creator : v == local.creator]
  )
}

locals {
  risk_level = [for v in data.huaweicloud_coc_scripts.test.data[*].properties[0].risk_level : v if v != ""][0]
}

data "huaweicloud_coc_scripts" "risk_level_filter" {
  risk_level = local.risk_level
}

output "risk_level_filter_is_useful" {
  value = length(data.huaweicloud_coc_scripts.risk_level_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scripts.risk_level_filter.data[*].properties[0].risk_level : v == local.risk_level]
  )
}

locals {
  type = [for v in data.huaweicloud_coc_scripts.test.data[*].type : v if v != ""][0]
}

data "huaweicloud_coc_scripts" "type_filter" {
  type = local.type
}

locals {
  type_filter_result = data.huaweicloud_coc_scripts.type_filter.data[0].type
}

output "type_filter_is_useful" {
  value = local.type_filter_result == local.type
}

locals {
  enterprise_project_id = [for v in data.huaweicloud_coc_scripts.test.data[*].enterprise_project_id : v if v != ""][0]
}

data "huaweicloud_coc_scripts" "enterprise_project_id_filter" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  enterprise_project_id_filter_result = data.huaweicloud_coc_scripts.enterprise_project_id_filter.data[0].enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = local.enterprise_project_id_filter_result == local.enterprise_project_id
}
`
}
