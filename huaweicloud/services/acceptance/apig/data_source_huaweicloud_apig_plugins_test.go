package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePlugins_basic(t *testing.T) {
	var (
		rName       = acceptance.RandomAccResourceName()
		randomId, _ = uuid.GenerateUUID()

		dataSource = "data.huaweicloud_apig_plugins.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byPluginId   = "data.huaweicloud_apig_plugins.filter_by_plugin_id"
		dcByPluginId = acceptance.InitDataSourceCheck(byPluginId)

		byNotFoundPluginId   = "data.huaweicloud_apig_plugins.filter_by_not_found_plugin_id"
		dcByNotFoundPluginId = acceptance.InitDataSourceCheck(byNotFoundPluginId)

		byPluginNameFuzzy   = "data.huaweicloud_apig_plugins.fuzzy_filter_by_plugin_name"
		dcByPluginNameFuzzy = acceptance.InitDataSourceCheck(byPluginNameFuzzy)

		byPluginName   = "data.huaweicloud_apig_plugins.filter_by_plugin_name"
		dcByPluginName = acceptance.InitDataSourceCheck(byPluginName)

		byPluginType   = "data.huaweicloud_apig_plugins.filter_by_plugin_type"
		dcByPluginType = acceptance.InitDataSourceCheck(byPluginType)

		byPluginScope   = "data.huaweicloud_apig_plugins.filter_by_plugin_scope"
		dcByPluginScope = acceptance.InitDataSourceCheck(byPluginScope)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourcePlugins_basic_step1(randomId),
				ExpectError: regexp.MustCompile(`The instance does not exist`),
			},
			{
				Config: testAccDataSourcePlugins_basic_step2(rName, randomId),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "plugins.#", regexp.MustCompile(`[1-9]\d*`)),
					dcByPluginId.CheckResourceExists(),
					resource.TestCheckOutput("is_plugin_id_filter_useful", "true"),
					dcByNotFoundPluginId.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_plugin_id_useful", "true"),
					dcByPluginNameFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_plugin_name_fuzzy_filter_useful", "true"),
					dcByPluginName.CheckResourceExists(),
					resource.TestCheckOutput("is_plugin_name_filter_useful", "true"),
					dcByPluginType.CheckResourceExists(),
					resource.TestCheckOutput("is_plugin_type_filter_useful", "true"),
					dcByPluginScope.CheckResourceExists(),
					resource.TestCheckOutput("is_plugin_scope_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byPluginId, "plugins.0.content"),
					resource.TestCheckResourceAttrPair(byPluginId, "plugins.0.description", "huaweicloud_apig_plugin.test", "description"),
					resource.TestMatchResourceAttr(byPluginId, "plugins.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byPluginId, "plugins.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourcePlugins_basic_step1(randomId string) string {
	return fmt.Sprintf(`
data "huaweicloud_apig_plugins" "test" {
  instance_id = "%[1]s"
}
`, randomId)
}

func testAccDataSourcePlugins_basic_step2(name, randomId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_plugin" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by acc test"
  type        = "cors"

  content = jsonencode(
    {
      allow_origin      = "*"
      allow_methods     = "GET,PUT,DELETE,HEAD,PATCH"
      allow_headers     = "Content-Type,Accept,Cache-Control"
      expose_headers    = "X-Request-Id,X-Apig-Latency"
      max_age           = 12700
      allow_credentials = true
    }
  )
}

data "huaweicloud_apig_plugins" "test" {
  instance_id = "%[1]s"

  depends_on = [huaweicloud_apig_plugin.test]
}

# Filter by plugin ID.
locals {
  plugin_id = huaweicloud_apig_plugin.test.id
}

data "huaweicloud_apig_plugins" "filter_by_plugin_id" {
  instance_id = "%[1]s"
  plugin_id   = local.plugin_id
}

locals {
  plugin_id_filter_result = [for v in data.huaweicloud_apig_plugins.filter_by_plugin_id.plugins : v.id == local.plugin_id]
}

output "is_plugin_id_filter_useful" {
  value = length(local.plugin_id_filter_result) > 0 && alltrue(local.plugin_id_filter_result)
}

# Filter by plugin IDs that don't exist.
data "huaweicloud_apig_plugins" "filter_by_not_found_plugin_id" {
  instance_id = "%[1]s"
  plugin_id   = "%[3]s"

  depends_on = [huaweicloud_apig_plugin.test]
}

output "is_not_found_plugin_id_useful" {
  value = length(data.huaweicloud_apig_plugins.filter_by_not_found_plugin_id.plugins) == 0
}

# Filter fuzzy search by plugin name.
data "huaweicloud_apig_plugins" "fuzzy_filter_by_plugin_name" {
  instance_id = "%[1]s"
  name        = "tf_test"

  depends_on = [huaweicloud_apig_plugin.test]
}

locals {
  plugin_name_fuzzy_filter_result = [for v in data.huaweicloud_apig_plugins.fuzzy_filter_by_plugin_name.plugins : strcontains(v.name, "tf_test")]
}

output "is_plugin_name_fuzzy_filter_useful" {
  value = length(local.plugin_name_fuzzy_filter_result) > 0 && alltrue(local.plugin_name_fuzzy_filter_result)
}

# Filter exact search by plugin name.
locals {
  plugin_name = huaweicloud_apig_plugin.test.name
}

data "huaweicloud_apig_plugins" "filter_by_plugin_name" {
  instance_id    = "%[1]s"
  name           = local.plugin_name
  precise_search = "name"

  depends_on = [huaweicloud_apig_plugin.test]
}

locals {
  plugin_name_filter_result = [for v in data.huaweicloud_apig_plugins.filter_by_plugin_name.plugins : v.name == local.plugin_name]
}

output "is_plugin_name_filter_useful" {
  value = length(local.plugin_name_filter_result) > 0 && alltrue(local.plugin_name_filter_result)
}

# Filter by plugin type.
locals {
  plugin_type = huaweicloud_apig_plugin.test.type
}

data "huaweicloud_apig_plugins" "filter_by_plugin_type" {
  instance_id = "%[1]s"
  type        = local.plugin_type

  depends_on = [huaweicloud_apig_plugin.test]
}

locals {
  plugin_type_filter_result = [for v in data.huaweicloud_apig_plugins.filter_by_plugin_type.plugins : v.type == local.plugin_type]
}

output "is_plugin_type_filter_useful" {
  value = length(local.plugin_type_filter_result) > 0 && alltrue(local.plugin_type_filter_result)
}

# Filter by plugin scope.
locals {
  plugin_scope = data.huaweicloud_apig_plugins.test.plugins[0].plugin_scope
}

data "huaweicloud_apig_plugins" "filter_by_plugin_scope" {
  instance_id  = "%[1]s"
  plugin_scope = local.plugin_scope
}

locals {
  plugin_scope_filter_result = [for v in data.huaweicloud_apig_plugins.filter_by_plugin_scope.plugins : v.plugin_scope == local.plugin_scope]
}

output "is_plugin_scope_filter_useful" {
  value = length(local.plugin_scope_filter_result) > 0 && alltrue(local.plugin_scope_filter_result)
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, randomId)
}
