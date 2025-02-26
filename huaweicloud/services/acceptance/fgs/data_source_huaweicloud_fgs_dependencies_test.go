package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDependencies_basic(t *testing.T) {
	var (
		base = "huaweicloud_fgs_dependency_version.test"

		all                  = "data.huaweicloud_fgs_dependencies.all"
		dcForAllDependencies = acceptance.InitDataSourceCheck(all)

		byType   = "data.huaweicloud_fgs_dependencies.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byName                    = "data.huaweicloud_fgs_dependencies.filter_by_name"
		dcByName                  = acceptance.InitDataSourceCheck(byName)
		byNotFoundName            = "data.huaweicloud_fgs_dependencies.filter_by_not_found_name"
		dcByNotFoundName          = acceptance.InitDataSourceCheck(byNotFoundName)
		byNameWithVersionsQuery   = "data.huaweicloud_fgs_dependencies.filter_by_name_with_versions_query"
		dcByNameWithVersionsQuery = acceptance.InitDataSourceCheck(byNameWithVersionsQuery)

		byRuntime   = "data.huaweicloud_fgs_dependencies.filter_by_runtime"
		dcByRuntime = acceptance.InitDataSourceCheck(byRuntime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsDependencyLink(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataDependencies_invalidRuntime,
				ExpectError: regexp.MustCompile(`Invalid runtime.`),
			},
			{
				Config: testAccDataDependencies_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcForAllDependencies.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "packages.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by dependency type.
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					// Filter by dependency name.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
					dcByNameWithVersionsQuery.CheckResourceExists(),
					resource.TestCheckOutput("is_versions_query_allowed_param_useful", "true"),
					// Filter by dependency runtime.
					dcByRuntime.CheckResourceExists(),
					resource.TestCheckOutput("is_runtime_filter_useful", "true"),
					// Check the attributes.
					// The behavior of the filter parameter 'name' is exact match.
					resource.TestCheckResourceAttrPair(byName, "packages.0.id", base, "dependency_id"),
					resource.TestCheckResourceAttrPair(byName, "packages.0.name", base, "name"),
					resource.TestCheckResourceAttrPair(byName, "packages.0.owner", base, "owner"),
					// The link will be replaced with a new link which belongs to the FunctionGraph service.
					resource.TestCheckResourceAttrSet(byName, "packages.0.link"),
					resource.TestCheckResourceAttrPair(byName, "packages.0.etag", base, "etag"),
					resource.TestCheckResourceAttrPair(byName, "packages.0.size", base, "size"),
					resource.TestCheckResourceAttrPair(byName, "packages.0.runtime", base, "runtime"),
					// The dependency version does not have the parameter 'file_name'.
					resource.TestCheckResourceAttrSet(byName, "packages.0.file_name"),
				),
			},
		},
	})
}

const testAccDataDependencies_invalidRuntime = `
data "huaweicloud_fgs_dependencies" "invalid_runtime" {
  runtime = "runtime_not_found"
}
`

func testAccDataDependencies_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_fgs_dependencies" "all" {}

# Filter by dependency type.
locals {
  dependency_type = "public"
}

data "huaweicloud_fgs_dependencies" "filter_by_type" {
  type = local.dependency_type
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_fgs_dependencies.filter_by_type.packages) > 0
}

# Filter by dependency name.
locals {
  dependency_name = huaweicloud_fgs_dependency_version.test.name
}

data "huaweicloud_fgs_dependencies" "filter_by_name" {
  # The behavior of parameter 'name' of the application resource is 'Required', means this parameter does not
  # have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_dependency_version.test,
  ]

  name = local.dependency_name
}

data "huaweicloud_fgs_dependencies" "filter_by_not_found_name" {
  name = "name_not_found"
}

data "huaweicloud_fgs_dependencies" "filter_by_name_with_versions_query" {
  # The behavior of parameter 'name' of the application resource is 'Required', means this parameter does not
  # have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_dependency_version.test,
  ]

  name                      = huaweicloud_fgs_dependency_version.test.name
  is_versions_query_allowed = true
}

locals {
  name_filter_result = [for v in data.huaweicloud_fgs_dependencies.filter_by_name.packages[*].name :
    v == local.dependency_name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

output "name_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_dependencies.filter_by_not_found_name.packages) == 0
}

locals {
  name_filter_with_versions_query_result = data.huaweicloud_fgs_dependencies.filter_by_name_with_versions_query.packages[0].versions
}

output "is_versions_query_allowed_param_useful" {
  value = length(local.name_filter_with_versions_query_result) > 0 && alltrue([
	for v in local.name_filter_with_versions_query_result : v.id != "" && v.version != 0
  ])
}	

# Filter by dependency runtime.
locals {
  dependency_runtime = data.huaweicloud_fgs_dependencies.all.packages[0].runtime
}

data "huaweicloud_fgs_dependencies" "filter_by_runtime" {
  runtime = local.dependency_runtime
}

locals {
  runtime_filter_result = [for v in data.huaweicloud_fgs_dependencies.filter_by_runtime.packages[*].runtime :
    v == local.dependency_runtime]
}

output "is_runtime_filter_useful" {
  value = length(local.runtime_filter_result) > 0 && alltrue(local.runtime_filter_result)
}
`, testAccDependencyVersion_basic(name))
}
