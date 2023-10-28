package dcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceTemplates_basic(t *testing.T) {
	rName := "data.huaweicloud_dcs_templates.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "templates.0.template_id"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.name"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.type"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.engine"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.engine_version"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.cache_mode"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.product_type"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.storage_type"),

					resource.TestCheckOutput("template_id_filter_is_useful", "true"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("engine_filter_is_useful", "true"),

					resource.TestCheckOutput("engine_version_filter_is_useful", "true"),

					resource.TestCheckOutput("cache_mode_filter_is_useful", "true"),

					resource.TestCheckOutput("product_type_filter_is_useful", "true"),

					resource.TestCheckOutput("storage_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceTemplates_basic() string {
	return `
data "huaweicloud_dcs_templates" "test" {
  type = "sys"
}

data "huaweicloud_dcs_templates" "template_id_filter" {
  type        = "sys"
  template_id = "6"

  depends_on = [data.huaweicloud_dcs_templates.test]
}
output "template_id_filter_is_useful" {
  value = length(data.huaweicloud_dcs_templates.template_id_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_templates.template_id_filter.templates[*].template_id : v == "6"]
  )  
}

data "huaweicloud_dcs_templates" "name_filter" {
  type = "sys"
  name = "Default-Redis-6.0-ha-enterprise-SSD"

  depends_on = [data.huaweicloud_dcs_templates.test]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dcs_templates.name_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_templates.name_filter.templates[*].name : v == "Default-Redis-6.0-ha-enterprise-SSD"]
  )  
}

data "huaweicloud_dcs_templates" "engine_filter" {
  type   = "sys"
  engine = "Redis"

  depends_on = [data.huaweicloud_dcs_templates.test]
}
output "engine_filter_is_useful" {
  value = length(data.huaweicloud_dcs_templates.engine_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_templates.engine_filter.templates[*].engine : v == "Redis"]
  )  
}

data "huaweicloud_dcs_templates" "engine_version_filter" {
  type           = "sys"
  engine_version = "6.0"

  depends_on = [data.huaweicloud_dcs_templates.test]
}
output "engine_version_filter_is_useful" {
  value = length(data.huaweicloud_dcs_templates.engine_version_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_templates.engine_version_filter.templates[*].engine_version : v == "6.0"]
  )  
}

data "huaweicloud_dcs_templates" "cache_mode_filter" {
  type       = "sys"
  cache_mode = "single"

  depends_on = [data.huaweicloud_dcs_templates.test]
}
output "cache_mode_filter_is_useful" {
  value = length(data.huaweicloud_dcs_templates.cache_mode_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_templates.cache_mode_filter.templates[*].cache_mode : v == "single"]
  )  
}

data "huaweicloud_dcs_templates" "product_type_filter" {
  type         = "sys"
  product_type = "generic"

  depends_on = [data.huaweicloud_dcs_templates.test]
}
output "product_type_filter_is_useful" {
  value = length(data.huaweicloud_dcs_templates.product_type_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_templates.product_type_filter.templates[*].product_type : v == "generic"]
  )  
}

data "huaweicloud_dcs_templates" "storage_type_filter" {
  type         = "sys"
  storage_type = "DRAM"

  depends_on = [data.huaweicloud_dcs_templates.test]
}
output "storage_type_filter_is_useful" {
  value = length(data.huaweicloud_dcs_templates.storage_type_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_templates.storage_type_filter.templates[*].storage_type : v == "DRAM"]
  )  
}
`
}
