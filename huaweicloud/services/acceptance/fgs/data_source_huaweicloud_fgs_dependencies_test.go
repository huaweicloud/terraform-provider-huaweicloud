package fgs

import (
	"fmt"
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
				Config: testAccFunctionGraphDependencies_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "packages.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

func TestAccFunctionGraphDependencies_name(t *testing.T) {
	dataSourceName := "data.huaweicloud_fgs_dependencies.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphDependencies_name(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "type", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "obssdk-3.0.2"),
					resource.TestCheckResourceAttr(dataSourceName, "packages.#", "1"),
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
				Config: testAccFunctionGraphDependencies_runtime(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "type", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "runtime", "Python2.7"),
					resource.TestMatchResourceAttr(dataSourceName, "packages.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

func testAccFunctionGraphDependencies_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_fgs_dependencies" "test" {}
`)
}

func testAccFunctionGraphDependencies_name() string {
	return fmt.Sprintf(`
data "huaweicloud_fgs_dependencies" "test" {
  type = "public"
  name = "obssdk-3.0.2"
}
`)
}

func testAccFunctionGraphDependencies_runtime() string {
	return fmt.Sprintf(`
data "huaweicloud_fgs_dependencies" "test" {
  type    = "public"
  runtime = "Python2.7"
}
`)
}
