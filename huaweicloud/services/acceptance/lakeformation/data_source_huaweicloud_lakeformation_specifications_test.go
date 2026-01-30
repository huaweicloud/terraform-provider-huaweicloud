package lakeformation

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSpecifications_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_lakeformation_specifications.all"
		dc  = acceptance.InitDataSourceCheck(all)

		bySpecCode   = "data.huaweicloud_lakeformation_specifications.filter_by_spec_code"
		dcBySpecCode = acceptance.InitDataSourceCheck(bySpecCode)

		byNotFoundSpecCode   = "data.huaweicloud_lakeformation_specifications.filter_by_not_found_spec_code"
		dcByNotFoundSpecCode = acceptance.InitDataSourceCheck(byNotFoundSpecCode)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpecifications_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "specifications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcBySpecCode.CheckResourceExists(),
					resource.TestCheckOutput("is_spec_code_filter_useful", "true"),
					dcByNotFoundSpecCode.CheckResourceExists(),
					resource.TestCheckOutput("is_spec_code_not_found_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceSpecifications_basic = `
data "huaweicloud_lakeformation_specifications" "all" {}

# Filter by spec_code
locals {
  first_spec_code = data.huaweicloud_lakeformation_specifications.all.specifications[0].spec_code
}

data "huaweicloud_lakeformation_specifications" "filter_by_spec_code" {
  spec_code = local.first_spec_code
}

locals {
  spec_code_filter_result = [
    for v in data.huaweicloud_lakeformation_specifications.filter_by_spec_code.specifications[*].spec_code :
	  v == local.first_spec_code
  ]
}

output "is_spec_code_filter_useful" {
  value = length(local.spec_code_filter_result) > 0 && alltrue(local.spec_code_filter_result)
}

# Filter by spec_code (not found)
locals {
  not_found_spec_code = "NOT_FOUND"
}

data "huaweicloud_lakeformation_specifications" "filter_by_not_found_spec_code" {
  spec_code = local.not_found_spec_code
}

locals {
  not_found_spec_code_filter_result = [
    for v in data.huaweicloud_lakeformation_specifications.filter_by_not_found_spec_code.specifications[*].spec_code :
	  strcontains(v, local.not_found_spec_code)
  ]
}

output "is_spec_code_not_found_filter_useful" {
  value = length(local.not_found_spec_code_filter_result) == 0
}
`
