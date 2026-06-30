package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBParameterTemplates_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_parameter_templates.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBParameterTemplates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.datastore_version_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.created"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.updated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.user_defined"),

					resource.TestCheckOutput("is_datastore_name_filter_useful", "true"),
					resource.TestCheckOutput("is_mode_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccGeminiDBParameterTemplates_basic = `
data "huaweicloud_geminidb_parameter_templates" "test" {}

data "huaweicloud_geminidb_parameter_templates" "datastore_name_filter" {
  datastore_name = "cassandra"
}

output "is_datastore_name_filter_useful" {
  value = length(data.huaweicloud_geminidb_parameter_templates.datastore_name_filter.configurations) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_parameter_templates.datastore_name_filter.configurations[*].datastore_name :
    v == "cassandra"]
  )
}

data "huaweicloud_geminidb_parameter_templates" "mode_filter" {
  mode = "Cluster"
}

output "is_mode_filter_useful" {
  value = length(data.huaweicloud_geminidb_parameter_templates.mode_filter.configurations) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_parameter_templates.mode_filter.configurations[*].mode :
    v == "Cluster"]
  )
}

data "huaweicloud_geminidb_parameter_templates" "user_defined_filter" {
  user_defined = "false"
}

output "user_defined_filter_useful" {
  value = length(data.huaweicloud_geminidb_parameter_templates.user_defined_filter.configurations) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_parameter_templates.user_defined_filter.configurations[*].user_defined :
    v == "false"]
  )
}
`
