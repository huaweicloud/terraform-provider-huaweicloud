package fgs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFunctionGraphDependencies_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_fgs_dependencies.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphDependenciesBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "packages.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

func TestAccFunctionGraphDependencies_filterByName(t *testing.T) {
	var (
		byName   = "data.huaweicloud_fgs_dependencies.filter_by_name"
		notFound = "data.huaweicloud_fgs_dependencies.not_found"

		dcByName   = acceptance.InitDataSourceCheck(byName)
		dcNotFound = acceptance.InitDataSourceCheck(notFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphDependenciesName,
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestMatchResourceAttr(byName, "packages.0.versions.#", regexp.MustCompile(`[1-9][0-9]*`)),
					resource.TestCheckResourceAttrSet(byName, "packages.0.versions.0.id"),
					resource.TestCheckResourceAttrSet(byName, "packages.0.versions.0.version"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func TestAccFunctionGraphDependencies_runtime(t *testing.T) {
	dataSourceName := "data.huaweicloud_fgs_dependencies.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphDependenciesRuntime,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "type", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "runtime", "Python2.7"),
					resource.TestMatchResourceAttr(dataSourceName, "packages.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

const testAccFunctionGraphDependenciesBasic = `data "huaweicloud_fgs_dependencies" "test" {}`

const testAccFunctionGraphDependenciesName = `
data "huaweicloud_fgs_dependencies" "filter_by_name" {
  type = "public"
  name = "obssdk-3.0.2"
}

locals {
  filter_result = [for v in data.huaweicloud_fgs_dependencies.filter_by_name.packages[*].name : v == "obssdk-3.0.2"]
}

output "is_name_filter_useful" {
  value = length(local.filter_result) > 0
}

data "huaweicloud_fgs_dependencies" "not_found" {
  type = "public"
  name = "not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_dependencies.not_found.packages) == 0
}
`

const testAccFunctionGraphDependenciesRuntime = `data "huaweicloud_fgs_dependencies" "test" {
  type    = "public"
  runtime = "Python2.7"
}`
