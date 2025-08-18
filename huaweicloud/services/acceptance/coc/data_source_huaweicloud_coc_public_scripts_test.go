package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocPublicScripts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_public_scripts.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocPublicScripts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.script_uuid"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.gmt_created"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.0.risk_level"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.0.version"),
					resource.TestCheckOutput("name_like_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("risk_level_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocPublicScripts_basic() string {
	return `
data "huaweicloud_coc_public_scripts" "test" {}

locals {
  name_like = substr([for v in data.huaweicloud_coc_public_scripts.test.data[*].name : v if v != ""][0], 0, 5)
}

data "huaweicloud_coc_public_scripts" "name_like_filter" {
  name_like = local.name_like
}

output "name_like_filter_is_useful" {
  value = length(data.huaweicloud_coc_public_scripts.name_like_filter.data) > 0
}

locals {
  name = [for v in data.huaweicloud_coc_public_scripts.test.data[*].name : v if v != ""][0]
}

data "huaweicloud_coc_public_scripts" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_coc_public_scripts.name_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_public_scripts.name_filter.data[*].name : v == local.name]
  )
}

data "huaweicloud_coc_public_scripts" "risk_level_filter" {
  risk_level = "HIGH"
}

output "risk_level_filter_is_useful" {
  value = length(data.huaweicloud_coc_public_scripts.risk_level_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_public_scripts.risk_level_filter.data[*].properties[0].risk_level : v == "HIGH"]
  )
}

locals {
  type = [for v in data.huaweicloud_coc_public_scripts.test.data[*].type : v if v != ""][0]
}

data "huaweicloud_coc_public_scripts" "type_filter" {
  type = local.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_coc_public_scripts.type_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_public_scripts.type_filter.data[*].type : v == local.type]
  )
}
`
}
