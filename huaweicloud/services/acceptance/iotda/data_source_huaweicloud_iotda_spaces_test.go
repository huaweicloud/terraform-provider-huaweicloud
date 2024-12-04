package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSpaces_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_spaces.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaces_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "spaces.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "spaces.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "spaces.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "spaces.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "spaces.0.is_default"),

					resource.TestCheckOutput("is_space_id_filter_useful", "true"),
					resource.TestCheckOutput("is_space_name_filter_useful", "true"),
					resource.TestCheckOutput("is_default_true_filter_useful", "true"),
					resource.TestCheckOutput("is_default_false_filter_useful", "true"),
					resource.TestCheckOutput("is_default_empty_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSpaces_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_spaces" "test" {
  depends_on = [huaweicloud_iotda_space.test]
}

# Filter using space ID.
locals {
  space_id = data.huaweicloud_iotda_spaces.test.spaces[0].id
}

data "huaweicloud_iotda_spaces" "space_id_filter" {
  space_id = local.space_id
}

output "is_space_id_filter_useful" {
  value = length(data.huaweicloud_iotda_spaces.space_id_filter.spaces) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_spaces.space_id_filter.spaces[*].id : v == local.space_id]
  )
}

# Filter using space name.
locals {
  space_name = data.huaweicloud_iotda_spaces.test.spaces[0].name
}

data "huaweicloud_iotda_spaces" "space_name_filter" {
  space_name = local.space_name
}

output "is_space_name_filter_useful" {
  value = length(data.huaweicloud_iotda_spaces.space_name_filter.spaces) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_spaces.space_name_filter.spaces[*].name : v == local.space_name]
  )
}

# Filter with is_default field set to true. There will definitely be a default space.
data "huaweicloud_iotda_spaces" "is_default_true_filter" {
  depends_on = [huaweicloud_iotda_space.test]

  is_default = "true"
}

output "is_default_true_filter_useful" {
  value = length(data.huaweicloud_iotda_spaces.is_default_true_filter.spaces) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_spaces.is_default_true_filter.spaces[*].is_default : v == true]
  )
}

# Filter with is_default field set to false. Set to false to query non default spaces.
data "huaweicloud_iotda_spaces" "is_default_false_filter" {
  depends_on = [huaweicloud_iotda_space.test]

  is_default = "false"
}

output "is_default_false_filter_useful" {
  value = length(data.huaweicloud_iotda_spaces.is_default_false_filter.spaces) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_spaces.is_default_false_filter.spaces[*].is_default : v == false]
  )
}

# Filter with is_default field is empty. Set to empty to query all spaces.
data "huaweicloud_iotda_spaces" "is_default_empty_filter" {
  depends_on = [huaweicloud_iotda_space.test]
}

output "is_default_empty_filter_useful" {
  value = length(data.huaweicloud_iotda_spaces.is_default_empty_filter.spaces) > 1
}

# Filter using non existent space name.
data "huaweicloud_iotda_spaces" "not_found" {
  space_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_spaces.not_found.spaces) == 0
}
`, testSpace_basic(name))
}
