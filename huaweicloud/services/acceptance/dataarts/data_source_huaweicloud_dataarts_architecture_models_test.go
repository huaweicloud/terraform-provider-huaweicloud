package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureModels_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_models.test"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByType   = "data.huaweicloud_dataarts_architecture_models.filter_by_type"
		dcFilterByType = acceptance.InitDataSourceCheck(filterByType)

		filterByDwType   = "data.huaweicloud_dataarts_architecture_models.filter_by_dw_type"
		dcFilterByDwType = acceptance.InitDataSourceCheck(filterByDwType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchitectureModels_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "models.#", regexp.MustCompile("^[1-9]([0-9]*)?$")),
					resource.TestCheckResourceAttrSet(all, "models.0.id"),
					resource.TestCheckResourceAttrSet(all, "models.0.name"),
					resource.TestCheckResourceAttrSet(all, "models.0.type"),
					resource.TestCheckResourceAttrSet(all, "models.0.dw_type"),
					resource.TestCheckResourceAttrSet(all, "models.0.is_physical"),
					resource.TestCheckResourceAttrSet(all, "models.0.create_time"),
					resource.TestCheckResourceAttrSet(all, "models.0.update_time"),
					resource.TestCheckResourceAttrSet(all, "models.0.create_by"),
					resource.TestCheckResourceAttrSet(all, "models.0.update_by"),

					// filter by type
					dcFilterByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),

					// filter by dw type
					dcFilterByDwType.CheckResourceExists(),
					resource.TestCheckOutput("is_dw_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureModels_base(name string) string {
	return fmt.Sprintf(`
variable "models" {
  type    = list(object({
    name = string
    type = string
  }))
  default = [
    {
      name = "%[1]s_third_nf"
      type = "THIRD_NF"
    },
    {
      name = "%[1]s_dimension"
      type = "DIMENSION"
    }
  ]
}

resource "huaweicloud_dataarts_architecture_model" "test" {
  count = length(var.models)

  workspace_id = "%[2]s"
  name         = var.models[count.index].name
  type         = var.models[count.index].type
  physical     = true
  dw_type      = "DWS"
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDataArchitectureModels_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all models in the workspace.
data "huaweicloud_dataarts_architecture_models" "test" {
  workspace_id = "%[2]s"

  depends_on = [huaweicloud_dataarts_architecture_model.test]
}

# Filter by workspace type.
locals {
  workspace_type = "THIRD_NF"
}

data "huaweicloud_dataarts_architecture_models" "filter_by_type" {
  workspace_id   = "%[2]s"
  workspace_type = local.workspace_type

  depends_on = [huaweicloud_dataarts_architecture_model.test]
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_models.filter_by_type.models :
	v.type == local.workspace_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by dw type.
locals {
  dw_type = "DWS"
}

data "huaweicloud_dataarts_architecture_models" "filter_by_dw_type" {
  workspace_id = "%[2]s"
  dw_type      = local.dw_type

  depends_on = [huaweicloud_dataarts_architecture_model.test]
}

locals {
  dw_type_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_models.filter_by_dw_type.models :
	v.dw_type == local.dw_type
  ]
}

output "is_dw_type_filter_useful" {
  value = length(local.dw_type_filter_result) > 0 && alltrue(local.dw_type_filter_result)
}
`, testAccDataArchitectureModels_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
