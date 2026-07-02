package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDevServerImages_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelarts_devserver_images.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byServerType   = "data.huaweicloud_modelarts_devserver_images.filter_by_server_type"
		dcByServerType = acceptance.InitDataSourceCheck(byServerType)

		byResourceFlavorName   = "data.huaweicloud_modelarts_devserver_images.filter_by_resource_flavor_name"
		dcByResourceFlavorName = acceptance.InitDataSourceCheck(byResourceFlavorName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDevServerImages_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "images.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "images.0.id"),
					resource.TestCheckResourceAttrSet(all, "images.0.name"),
					resource.TestCheckResourceAttrSet(all, "images.0.server_type"),
					resource.TestCheckResourceAttrSet(all, "images.0.arch"),
					resource.TestCheckResourceAttrSet(all, "images.0.status"),

					dcByServerType.CheckResourceExists(),
					resource.TestCheckOutput("is_server_type_filter_useful", "true"),

					dcByResourceFlavorName.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_flavor_name_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceDevServerImages_basic string = `
# Query all DevServer images without any filter.
data "huaweicloud_modelarts_devserver_images" "all" {}

# Filter by parameter 'server_type'.
locals {
  server_type = data.huaweicloud_modelarts_devserver_images.all.images[0].server_type
}

data "huaweicloud_modelarts_devserver_images" "filter_by_server_type" {
  server_type = local.server_type
}

locals {
  filter_by_server_type_result = [
    for v in data.huaweicloud_modelarts_devserver_images.filter_by_server_type.images[*].server_type : v == local.server_type
  ]
}

output "is_server_type_filter_useful" {
  value = length(local.filter_by_server_type_result) > 0 && alltrue(local.filter_by_server_type_result)
}

# Filter by parameter 'resource_flavor_name'.
data "huaweicloud_modelarts_devserver_flavors" "test" {}

locals {
  resource_flavor_name = try(data.huaweicloud_modelarts_devserver_flavors.test.flavors[0].flavor, "")
}

data "huaweicloud_modelarts_devserver_images" "filter_by_resource_flavor_name" {
  resource_flavor_name = local.resource_flavor_name
}

locals {
  filter_by_resource_flavor_name_result = length(data.huaweicloud_modelarts_devserver_images.filter_by_resource_flavor_name.images) > 0
}

output "is_resource_flavor_name_filter_useful" {
  value = local.filter_by_resource_flavor_name_result
}
`
