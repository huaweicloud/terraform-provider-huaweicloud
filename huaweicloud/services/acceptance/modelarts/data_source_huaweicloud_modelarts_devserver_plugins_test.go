package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDevServerPlugins_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelarts_devserver_plugins.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_modelarts_devserver_plugins.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDevServerPlugins_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "plugins.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(all, "plugins.0.name"),
					resource.TestMatchResourceAttr(all, "plugins.0.infos.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(all, "plugins.0.infos.0.id"),
					resource.TestCheckResourceAttrSet(all, "plugins.0.infos.0.version"),
					resource.TestCheckResourceAttrSet(all, "plugins.0.infos.0.status"),
					resource.TestCheckResourceAttrSet(all, "plugins.0.infos.0.url"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDevServerPlugins_basic() string {
	return `
# Query all DevServer plugins and without any filter
data "huaweicloud_modelarts_devserver_plugins" "all" {}

# Filter by name
locals {
  plugin_name = data.huaweicloud_modelarts_devserver_plugins.all.plugins[0].name
}

data "huaweicloud_modelarts_devserver_plugins" "filter_by_name" {
  name = local.plugin_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_modelarts_devserver_plugins.filter_by_name.plugins[*].name :
      v == local.plugin_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`
}
