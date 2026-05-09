package modelarts

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV2PluginTemplates_basic(t *testing.T) {
	var (
		name = strings.ReplaceAll(acceptance.RandomAccResourceName(), "_", "-")

		all = "data.huaweicloud_modelartsv2_plugin_templates.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByTemplateName   = "data.huaweicloud_modelartsv2_plugin_templates.filter_by_template_name"
		dcFilterByTemplateName = acceptance.InitDataSourceCheck(filterByTemplateName)

		filterByPoolName   = "data.huaweicloud_modelartsv2_plugin_templates.filter_by_pool_name"
		dcFilterByPoolName = acceptance.InitDataSourceCheck(filterByPoolName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV2PluginTemplates_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "plugin_templates.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// filter by template name
					dcFilterByTemplateName.CheckResourceExists(),
					resource.TestCheckOutput("is_template_name_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(filterByTemplateName, "plugin_templates.0.metadata.0.name",
						all, "plugin_templates.0.metadata.0.name"),
					resource.TestMatchResourceAttr(filterByTemplateName, "plugin_templates.0.metadata.0.annotations.%",
						regexp.MustCompile(`^[1-9][0-9]*$`)),
					resource.TestCheckResourceAttrPair(filterByTemplateName, "plugin_templates.0.spec.0.optional",
						all, "plugin_templates.0.spec.0.optional"),
					resource.TestCheckResourceAttrPair(filterByTemplateName, "plugin_templates.0.spec.0.type",
						all, "plugin_templates.0.spec.0.type"),
					resource.TestCheckResourceAttrPair(filterByTemplateName, "plugin_templates.0.spec.0.logo_url",
						all, "plugin_templates.0.spec.0.logo_url"),
					resource.TestCheckResourceAttrPair(filterByTemplateName, "plugin_templates.0.spec.0.description",
						all, "plugin_templates.0.spec.0.description"),
					resource.TestMatchResourceAttr(filterByTemplateName, "plugin_templates.0.spec.0.versions.#",
						regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// filter by pool name
					dcFilterByPoolName.CheckResourceExists(),
					resource.TestCheckOutput("is_pool_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataV2PluginTemplates_basic_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_plugin_templates" "all" {}

resource "huaweicloud_modelarts_network" "test" {
  name = "%[1]s"
  cidr = "172.16.0.0/12"
}

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "%[1]s"
  description = "This is a demo"
  scope       = ["Train", "Infer", "Notebook"]
  network_id  = huaweicloud_modelarts_network.test.id

  resources {
    flavor_id = "modelarts.vm.cpu.16u64g.d"
    count     = 1
  }
}
`, name)
}

func testAccDataV2PluginTemplates_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Filter by template name
locals {
  template_name = (
    length(data.huaweicloud_modelartsv2_plugin_templates.all.plugin_templates) > 0 ? 
  data.huaweicloud_modelartsv2_plugin_templates.all.plugin_templates[0].metadata[0].name : ""
  )
}

data "huaweicloud_modelartsv2_plugin_templates" "filter_by_template_name" {
  depends_on = [
    data.huaweicloud_modelartsv2_plugin_templates.all,
  ]

  template_name = local.template_name
}

locals {
  template_name_filter_result = [
    for v in data.huaweicloud_modelartsv2_plugin_templates.filter_by_template_name.plugin_templates : v.metadata[0].name == local.template_name
  ]
}

output "is_template_name_filter_useful" {
  value = length(local.template_name_filter_result) > 0 && alltrue(local.template_name_filter_result)
}

# Filter by pool name
locals {
  pool_name = huaweicloud_modelarts_resource_pool.test.id
}

data "huaweicloud_modelartsv2_plugin_templates" "filter_by_pool_name" {
  depends_on = [
    huaweicloud_modelarts_resource_pool.test,
  ]

  pool_name = local.pool_name
}

locals {
  pool_name_filter_result = length(data.huaweicloud_modelartsv2_plugin_templates.filter_by_pool_name.plugin_templates) > 0
}

output "is_pool_name_filter_useful" {
  value = local.pool_name_filter_result
}
`, testAccDataV2PluginTemplates_basic_base(name))
}
