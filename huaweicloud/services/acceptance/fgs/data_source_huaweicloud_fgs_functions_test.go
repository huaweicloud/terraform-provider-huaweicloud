package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFunctions_basic(t *testing.T) {
	var (
		rcName            = "huaweicloud_fgs_function.test"
		all               = "data.huaweicloud_fgs_functions.all"
		dcForAllFunctions = acceptance.InitDataSourceCheck(all)

		byPackageName           = "data.huaweicloud_fgs_functions.filter_by_package_name"
		dcByPackageName         = acceptance.InitDataSourceCheck(byPackageName)
		byNotFoundPackageName   = "data.huaweicloud_fgs_functions.filter_by_not_found_package_name"
		dcByNotFoundPackageName = acceptance.InitDataSourceCheck(byNotFoundPackageName)

		byFunctionUrn           = "data.huaweicloud_fgs_functions.filter_by_function_urn"
		dcByFunctionUrn         = acceptance.InitDataSourceCheck(byFunctionUrn)
		byNotFoundFunctionUrn   = "data.huaweicloud_fgs_functions.filter_by_not_found_function_urn"
		dcByNotFoundFunctionUrn = acceptance.InitDataSourceCheck(byNotFoundFunctionUrn)

		byFunctionName           = "data.huaweicloud_fgs_functions.filter_by_function_name"
		dcByFunctionName         = acceptance.InitDataSourceCheck(byFunctionName)
		byNotFoundFunctionName   = "data.huaweicloud_fgs_functions.filter_by_not_found_function_name"
		dcByNotFoundFunctionName = acceptance.InitDataSourceCheck(byNotFoundFunctionName)

		byRuntime           = "data.huaweicloud_fgs_functions.filter_by_function_runtime"
		dcByRuntime         = acceptance.InitDataSourceCheck(byRuntime)
		byNotFoundRuntime   = "data.huaweicloud_fgs_functions.filter_by_not_found_function_runtime"
		dcByNotFoundRuntime = acceptance.InitDataSourceCheck(byNotFoundRuntime)

		byEpsId            = "data.huaweicloud_fgs_functions.filter_by_eps_id"
		dcByEpsId          = acceptance.InitDataSourceCheck(byEpsId)
		byNotFoundEpsId    = "data.huaweicloud_fgs_functions.filter_by_not_found_eps_id"
		dcByNotFoundPEpsId = acceptance.InitDataSourceCheck(byNotFoundEpsId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFunctions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcForAllFunctions.CheckResourceExists(),
					resource.TestCheckOutput("is_result_contains_function", "true"),
					// Filter by package name.
					dcByPackageName.CheckResourceExists(),
					resource.TestCheckOutput("is_package_name_filter_useful", "true"),
					dcByNotFoundPackageName.CheckResourceExists(),
					resource.TestCheckOutput("package_name_not_found_validation_pass", "true"),
					// Filter by function URN.
					dcByFunctionUrn.CheckResourceExists(),
					resource.TestCheckOutput("is_function_urn_filter_useful", "true"),
					dcByNotFoundFunctionUrn.CheckResourceExists(),
					resource.TestCheckOutput("function_urn_not_found_validation_pass", "true"),
					// Filter by function name.
					dcByFunctionName.CheckResourceExists(),
					resource.TestCheckOutput("is_function_name_filter_useful", "true"),
					dcByNotFoundFunctionName.CheckResourceExists(),
					resource.TestCheckOutput("function_name_not_found_validation_pass", "true"),
					// Filter by function runtime.
					dcByRuntime.CheckResourceExists(),
					resource.TestCheckOutput("is_function_runtime_filter_useful", "true"),
					dcByNotFoundRuntime.CheckResourceExists(),
					resource.TestCheckOutput("function_runtime_not_found_validation_pass", "true"),
					// Filter by enterprise project ID.
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					dcByNotFoundPEpsId.CheckResourceExists(),
					resource.TestCheckOutput("eps_id_not_found_validation_pass", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.name", rcName, "name"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.urn", rcName, "urn"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.package", rcName, "package"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.runtime", rcName, "runtime"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.timeout", rcName, "timeout"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.handler", rcName, "handler"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.memory_size", rcName, "memory_size"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.code_type", rcName, "code_type"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.description", rcName, "description"),
				),
			},
		},
	})
}

func testAccDataFunctions_base(rName string) string {
	return fmt.Sprintf(`
variable "function_code_content" {
  type    = string
  default = <<EOT
def main():  
    print("Hello, World!")  

if __name__ == "__main__":  
    main()
EOT
}

resource "huaweicloud_fgs_function" "test" {
  name                  = "%[1]s"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  handler               = "index.handler"
  package               = "default"
  description           = "fuction test"
  enterprise_project_id = "%[2]s"
  code_type             = "inline"
  func_code             = base64encode(var.function_code_content)
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataFunctions_basic() string {
	randName := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_fgs_functions" "all" {
  depends_on = [
    huaweicloud_fgs_function.test
  ]
}

output "is_result_contains_function" {
  value = contains(data.huaweicloud_fgs_functions.all.functions[*].name, huaweicloud_fgs_function.test.name)
}

# Filter by package name.
locals {
  package_name = huaweicloud_fgs_function.test.package
}

data "huaweicloud_fgs_functions" "filter_by_package_name" {
  // The behavior of parameter 'package_name' of the resource is 'Required', means this parameter does not 
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  package_name = local.package_name
}

data "huaweicloud_fgs_functions" "filter_by_not_found_package_name" {
  // Query functions using a not exist package name after function resource create.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  package_name = "package_name_not_found"
}

locals {
  package_name_filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_package_name.functions[*].package :
    v == local.package_name]
}

output "is_package_name_filter_useful" {
  value = length(local.package_name_filter_result) > 0 && alltrue(local.package_name_filter_result)
}

output "package_name_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.filter_by_not_found_package_name.functions) == 0
}

# Filter by function URN.
locals {
  function_urn = huaweicloud_fgs_function.test.urn
}

data "huaweicloud_fgs_functions" "filter_by_function_urn" {
  urn = local.function_urn
}

data "huaweicloud_fgs_functions" "filter_by_not_found_function_urn" {
  // Query functions using a not exist function URN after function resource create.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  urn = "function_urn_not_found"
}

locals {
  function_urn_filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_function_urn.functions[*].urn :
    v == local.function_urn]
}

output "is_function_urn_filter_useful" {
  value = length(local.function_urn_filter_result) > 0 && alltrue(local.function_urn_filter_result)
}

output "function_urn_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.filter_by_not_found_function_urn.functions) == 0
}

# Filter by function name.
locals {
  function_name = huaweicloud_fgs_function.test.name
}

data "huaweicloud_fgs_functions" "filter_by_function_name" {
  // The behavior of parameter 'name' of the resource is 'Required', means this parameter does not 
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  name = local.function_name
}

data "huaweicloud_fgs_functions" "filter_by_not_found_function_name" {
  // Query functions using a not exist function name after function resource create.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  name = "function_name_not_found"
}

locals {
  function_name_filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_function_urn.functions[*].name :
    v == local.function_name]
}

output "is_function_name_filter_useful" {
  value = length(local.function_name_filter_result) > 0 && alltrue(local.function_name_filter_result)
}

output "function_name_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.filter_by_not_found_function_name.functions) == 0
}

# Filter by function runtime.
locals {
  function_runtime = huaweicloud_fgs_function.test.runtime
}

data "huaweicloud_fgs_functions" "filter_by_function_runtime" {
  // The behavior of parameter 'runtime' of the resource is 'Required', means this parameter does not 
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  runtime = local.function_runtime
}

data "huaweicloud_fgs_functions" "filter_by_not_found_function_runtime" {
  // Query functions using a not exist runtime after function resource create.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  runtime = "function_runtime_not_found"
}

locals {
  function_runtime_filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_function_runtime.functions[*].runtime :
    v == local.function_runtime]
}

output "is_function_runtime_filter_useful" {
  value = length(local.function_runtime_filter_result) > 0 && alltrue(local.function_runtime_filter_result)
}

output "function_runtime_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.filter_by_not_found_function_runtime.functions) == 0
}

# Filter by function enterprise project ID.
locals {
  enterprise_project_id = huaweicloud_fgs_function.test.enterprise_project_id
}

data "huaweicloud_fgs_functions" "filter_by_eps_id" {
  enterprise_project_id = local.enterprise_project_id
}

data "huaweicloud_fgs_functions" "filter_by_not_found_eps_id" {
  // Query functions using a not exist enterprise project ID after function resource create.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  enterprise_project_id = "eps_id_not_found"
}

locals {
  eps_id_filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_function_urn.functions[*].enterprise_project_id :
    v == local.enterprise_project_id]
}

output "is_eps_id_filter_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}

output "eps_id_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.filter_by_not_found_eps_id.functions) == 0
}
`, testAccDataFunctions_base(randName))
}
