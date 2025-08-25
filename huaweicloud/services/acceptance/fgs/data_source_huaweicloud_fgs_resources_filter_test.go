package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataResourcesFilter_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all               = "data.huaweicloud_fgs_resources_filter.all"
		dcForAllResources = acceptance.InitDataSourceCheck(all)

		byMatches   = "data.huaweicloud_fgs_resources_filter.filter_by_matches"
		dcByMatches = acceptance.InitDataSourceCheck(byMatches)

		byTags   = "data.huaweicloud_fgs_resources_filter.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataResourcesFilter_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcForAllResources.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// Filter by matches.
					dcByMatches.CheckResourceExists(),
					resource.TestMatchResourceAttr(byMatches, "resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_matches_filter_useful", "true"),

					resource.TestCheckResourceAttrPair(byMatches, "resources.0.id", "huaweicloud_fgs_function.with_tags", "id"),
					resource.TestCheckResourceAttrPair(byMatches, "resources.0.name", "huaweicloud_fgs_function.with_tags", "name"),
					resource.TestCheckResourceAttrSet(byMatches, "resources.0.detail"),
					resource.TestMatchResourceAttr(byMatches, "resources.0.tags.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byMatches, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(byMatches, "resources.0.tags.0.value"),

					// Filter by tags.
					dcByTags.CheckResourceExists(),
					resource.TestMatchResourceAttr(byTags, "resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataResourcesFilter_base(name string) string {
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

resource "huaweicloud_fgs_function" "with_tags" {
  name        = "%[1]s_with_tags"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.function_code_content)

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_fgs_function" "without_tags" {
  name        = "%[1]s_without_tags"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.function_code_content)
}
`, name)
}

func testAccDataResourcesFilter_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_fgs_resources_filter" "all" {
  depends_on = [
    huaweicloud_fgs_function.with_tags,
    huaweicloud_fgs_function.without_tags,
  ]

  resource_type = "functions"
}

# Filter by matches
locals {
  function_name = huaweicloud_fgs_function.with_tags.name
}

data "huaweicloud_fgs_resources_filter" "filter_by_matches" {
  depends_on = [
    huaweicloud_fgs_function.with_tags,
    huaweicloud_fgs_function.without_tags,
  ]

  resource_type = "functions"

  matches {
    key   = "resource_name"
    value = local.function_name
  }
}

locals {
  matches_filter_result = [
    for v in data.huaweicloud_fgs_resources_filter.filter_by_matches.resources[*].name : strcontains(v, local.function_name)
  ]
}

output "is_matches_filter_useful" {
  value = length(local.matches_filter_result) > 0 && alltrue(local.matches_filter_result)
}

# Filter by tags
data "huaweicloud_fgs_resources_filter" "filter_by_tags" {
  depends_on = [
    huaweicloud_fgs_function.with_tags,
    huaweicloud_fgs_function.without_tags,
  ]

  resource_type = "functions"

  tags {
    key   = "foo"
    values = ["bar"]
  }
}

locals {
  tags_filter_result = [
    for v in data.huaweicloud_fgs_resources_filter.filter_by_tags.resources[*].tags : contains(v[*].key, "foo") && alltrue(
	  [for vv in v : vv.value == "bar" if vv.key == "foo"])
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}
`, testAccDataResourcesFilter_base(name), name)
}
