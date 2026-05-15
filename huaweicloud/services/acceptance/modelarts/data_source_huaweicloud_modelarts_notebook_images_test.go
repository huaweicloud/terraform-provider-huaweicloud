package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNotebookImages_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelarts_notebook_images.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_modelarts_notebook_images.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_modelarts_notebook_images.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byWorkspaceId   = "data.huaweicloud_modelarts_notebook_images.filter_by_workspace_id"
		dcByWorkspaceId = acceptance.InitDataSourceCheck(byWorkspaceId)

		byCpuArch   = "data.huaweicloud_modelarts_notebook_images.filter_by_cpu_arch"
		dcByCpuArch = acceptance.InitDataSourceCheck(byCpuArch)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceImages_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "images.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(all, "images.0.type", "BUILD_IN"),
					resource.TestCheckResourceAttr(all, "images.0.cpu_arch", "x86_64"),
					resource.TestCheckResourceAttrSet(all, "images.0.id"),
					resource.TestCheckResourceAttrSet(all, "images.0.name"),
					resource.TestCheckResourceAttrSet(all, "images.0.swr_path"),
					resource.TestCheckResourceAttrSet(all, "images.0.description"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByWorkspaceId.CheckResourceExists(),
					resource.TestCheckOutput("is_workspace_id_filter_useful", "true"),
					dcByCpuArch.CheckResourceExists(),
					resource.TestCheckOutput("is_cpu_arch_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceImages_basic string = `
# Query all notebook images without any filter.
data "huaweicloud_modelarts_notebook_images" "all" {}

# Filter by parameter 'name'.
locals {
  image_name = data.huaweicloud_modelarts_notebook_images.all.images[0].name
}

data "huaweicloud_modelarts_notebook_images" "filter_by_name" {
  name = local.image_name
}

locals {
  filter_by_name_result = [
    for v in data.huaweicloud_modelarts_notebook_images.filter_by_name.images[*].name : v == local.image_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.filter_by_name_result) > 0 && alltrue(local.filter_by_name_result)
}

# Filter by parameter 'type'.
locals {
  image_type = data.huaweicloud_modelarts_notebook_images.all.images[0].type
}

data "huaweicloud_modelarts_notebook_images" "filter_by_type" {
  type = local.image_type
}

locals {
  filter_by_type_result = [
    for v in data.huaweicloud_modelarts_notebook_images.filter_by_type.images[*].type : v == local.image_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.filter_by_type_result) > 0 && alltrue(local.filter_by_type_result)
}

# Filter by parameter 'workspace_id'.
locals {
  image_workspace_id = data.huaweicloud_modelarts_notebook_images.all.images[0].workspace_id
}

data "huaweicloud_modelarts_notebook_images" "filter_by_workspace_id" {
  workspace_id = local.image_workspace_id
}

locals {
  filter_by_workspace_id_result = [
    for v in data.huaweicloud_modelarts_notebook_images.filter_by_workspace_id.images[*].workspace_id : v == local.image_workspace_id
  ]
}

output "is_workspace_id_filter_useful" {
  value = length(local.filter_by_workspace_id_result) > 0 && alltrue(local.filter_by_workspace_id_result)
}

# Filter by parameter 'cpu_arch'.
locals {
  image_cpu_arch = data.huaweicloud_modelarts_notebook_images.all.images[0].cpu_arch
}

data "huaweicloud_modelarts_notebook_images" "filter_by_cpu_arch" {
  cpu_arch = local.image_cpu_arch
}

locals {
  filter_by_cpu_arch_result = [
    for v in data.huaweicloud_modelarts_notebook_images.filter_by_cpu_arch.images[*].cpu_arch : v == local.image_cpu_arch
  ]
}

output "is_cpu_arch_filter_useful" {
  value = length(local.filter_by_cpu_arch_result) > 0 && alltrue(local.filter_by_cpu_arch_result)
}
`
