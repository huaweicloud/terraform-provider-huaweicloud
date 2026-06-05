package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTrainingJobFlavors_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelarts_training_job_flavors.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byFlavorType   = "data.huaweicloud_modelarts_training_job_flavors.filter_by_flavor_type"
		dcByFlavorType = acceptance.InitDataSourceCheck(byFlavorType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTrainingJobFlavors_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "flavors.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_name"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_type"),
					resource.TestCheckResourceAttr(all, "flavors.0.billing.#", "1"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.billing.0.code"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.billing.0.unit_num"),
					resource.TestCheckResourceAttr(all, "flavors.0.flavor_info.#", "1"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_info.0.max_num"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_info.0.cpu.0.arch"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_info.0.cpu.0.core_num"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_info.0.memory.0.size"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_info.0.memory.0.unit"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_info.0.disk.0.size"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_info.0.disk.0.unit"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.support_engines"),
					// Filter by 'flavor_type' parameter.
					dcByFlavorType.CheckResourceExists(),
					resource.TestCheckOutput("is_flavor_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataTrainingJobFlavors_basic = `
# Without any filter parameters.
data "huaweicloud_modelarts_training_job_flavors" "test" {}

# Filter by 'flavor_type' parameter.
locals {
  flavor_type = try(data.huaweicloud_modelarts_training_job_flavors.test.flavors[0].flavor_type, "")
}

data "huaweicloud_modelarts_training_job_flavors" "filter_by_flavor_type" {
  flavor_type = local.flavor_type
}

locals {
  flavor_type_filter_result = [for v in data.huaweicloud_modelarts_training_job_flavors.filter_by_flavor_type.flavors[*].flavor_type :
    v == local.flavor_type
  ]
}

output "is_flavor_type_filter_useful" {
  value = length(local.flavor_type_filter_result) > 0 && alltrue(local.flavor_type_filter_result)
}
`
