package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFunctionGraphFunctions_basic(t *testing.T) {
	var (
		baseRes       = "huaweicloud_fgs_function.test"
		byPackageName = "data.huaweicloud_fgs_functions.filter_by_package_name"
		NotFound      = "data.huaweicloud_fgs_functions.not_found"

		obj             function.Function
		rc              = acceptance.InitResourceCheck(baseRes, &obj, getResourceObj)
		dcByPackageName = acceptance.InitDataSourceCheck(byPackageName)
		dcNotFound      = acceptance.InitDataSourceCheck(NotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphFunctions_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcByPackageName.CheckResourceExists(),
					resource.TestCheckOutput("is_package_name_filter_useful", "true"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccFunctionGraphFunctions_base(rName string) string {
	//nolint:revive
	return fmt.Sprintf(`
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
  func_code             = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccFunctionGraphFunctions_basic() string {
	randName := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_fgs_functions" "filter_by_package_name" {
  // The behavior of parameter 'package name' of the resource is 'Required', means this parameter does not 
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  package_name = huaweicloud_fgs_function.test.package
}

data "huaweicloud_fgs_functions" "not_found" {
  // Since a specified package name is used, there is no dependency relationship with resource attachment, 
  // and the dependency needs to be manually set.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  package_name = "resource_not_found"
}

locals {
  filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_package_name.functions[*].name :
                   v == huaweicloud_fgs_function.test.name]
}

output "is_package_name_filter_useful" {
  value = length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.not_found.functions) == 0
}
`, testAccFunctionGraphFunctions_base(randName))
}

func TestAccFunctionGraphFunctions_filterByName(t *testing.T) {
	var (
		baseRes  = "huaweicloud_fgs_function.test"
		byName   = "data.huaweicloud_fgs_functions.filter_by_name"
		NotFound = "data.huaweicloud_fgs_functions.not_found"

		obj        function.Function
		rc         = acceptance.InitResourceCheck(baseRes, &obj, getResourceObj)
		dcByName   = acceptance.InitDataSourceCheck(byName)
		dcNotFound = acceptance.InitDataSourceCheck(NotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphFunctions_filterByName(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccFunctionGraphFunctions_filterByName() string {
	randName := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_fgs_functions" "filter_by_name" {
  // The behavior of parameter 'name' of the resource is 'Required', means this parameter does not 
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  name = huaweicloud_fgs_function.test.name
}

data "huaweicloud_fgs_functions" "not_found" {
  // Since a specified name is used, there is no dependency relationship with resource attachment, 
  // and the dependency needs to be manually set.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  name = "resource_not_found"
}

locals {
  filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_name.functions[*].name :
                   v == huaweicloud_fgs_function.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.not_found.functions) == 0
}
`, testAccFunctionGraphFunctions_base(randName))
}

func TestAccFunctionGraphFunctions_filterByUrn(t *testing.T) {
	var (
		baseRes  = "huaweicloud_fgs_function.test"
		byUrn    = "data.huaweicloud_fgs_functions.filter_by_urn"
		NotFound = "data.huaweicloud_fgs_functions.not_found"

		obj        function.Function
		rc         = acceptance.InitResourceCheck(baseRes, &obj, getResourceObj)
		dcByUrn    = acceptance.InitDataSourceCheck(byUrn)
		dcNotFound = acceptance.InitDataSourceCheck(NotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphFunctions_filterByUrn(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcByUrn.CheckResourceExists(),
					resource.TestCheckOutput("is_urn_filter_useful", "true"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccFunctionGraphFunctions_filterByUrn() string {
	randName := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_fgs_functions" "filter_by_urn" {
  // The behavior of parameter 'urn' of the resource is 'Required', means this parameter does not 
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  urn = huaweicloud_fgs_function.test.urn
}

data "huaweicloud_fgs_functions" "not_found" {
  // Since a specified urn is used, there is no dependency relationship with resource attachment, 
  // and the dependency needs to be manually set.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  urn = "resource_not_found"
}

locals {
  filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_urn.functions[*].name :
                   v == huaweicloud_fgs_function.test.name]
}

output "is_urn_filter_useful" {
  value = length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.not_found.functions) == 0
}
`, testAccFunctionGraphFunctions_base(randName))
}

func TestAccFunctionGraphFunctions_filterByRuntime(t *testing.T) {
	var (
		baseRes   = "huaweicloud_fgs_function.test"
		byRuntime = "data.huaweicloud_fgs_functions.filter_by_runtime"
		NotFound  = "data.huaweicloud_fgs_functions.not_found"

		obj         function.Function
		rc          = acceptance.InitResourceCheck(baseRes, &obj, getResourceObj)
		dcByRuntime = acceptance.InitDataSourceCheck(byRuntime)
		dcNotFound  = acceptance.InitDataSourceCheck(NotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphFunctions_filterByRuntime(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcByRuntime.CheckResourceExists(),
					resource.TestCheckOutput("is_runtime_filter_useful", "true"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccFunctionGraphFunctions_filterByRuntime() string {
	randName := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_fgs_functions" "filter_by_runtime" {
  // The behavior of parameter 'runtime' of the resource is 'Required', means this parameter does not 
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  runtime = huaweicloud_fgs_function.test.runtime
}

data "huaweicloud_fgs_functions" "not_found" {
  // Since a specified runtime is used, there is no dependency relationship with resource attachment, 
  // and the dependency needs to be manually set.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  runtime = "resource_not_found"
}

locals {
  filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_runtime.functions[*].name :
                   v == huaweicloud_fgs_function.test.name]
}

output "is_runtime_filter_useful" {
  value = length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.not_found.functions) == 0
}
`, testAccFunctionGraphFunctions_base(randName))
}

func TestAccFunctionGraphFunctions_filterByEpsId(t *testing.T) {
	var (
		baseRes  = "huaweicloud_fgs_function.test"
		byEpsId  = "data.huaweicloud_fgs_functions.filter_by_eps_id"
		NotFound = "data.huaweicloud_fgs_functions.not_found"

		obj        function.Function
		rc         = acceptance.InitResourceCheck(baseRes, &obj, getResourceObj)
		dcByEpsId  = acceptance.InitDataSourceCheck(byEpsId)
		dcNotFound = acceptance.InitDataSourceCheck(NotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphFunctions_filterByEpsId(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccFunctionGraphFunctions_filterByEpsId() string {
	randName := acceptance.RandomAccResourceName()
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_fgs_functions" "filter_by_eps_id" {
  // The behavior of parameter 'eps_id' of the resource is 'Required', means this parameter does not 
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  enterprise_project_id = "%[2]s"
}

data "huaweicloud_fgs_functions" "not_found" {
  // Since a specified eps_id is used, there is no dependency relationship with resource attachment, 
  // and the dependency needs to be manually set.
  depends_on = [
    huaweicloud_fgs_function.test,
  ]

  enterprise_project_id = "%[3]s"
}

locals {
  filter_result = [for v in data.huaweicloud_fgs_functions.filter_by_eps_id.functions[*].name :
                   v == huaweicloud_fgs_function.test.name]
}

output "is_eps_id_filter_useful" {
  value = length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_functions.not_found.functions) == 0
}
`, testAccFunctionGraphFunctions_base(randName), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, randUUID)
}
