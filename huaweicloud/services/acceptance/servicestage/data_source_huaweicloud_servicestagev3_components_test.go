package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV3Components_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_servicestagev3_components.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure at least one of node exist.
			acceptance.TestAccPreCheckCceClusterId(t)
			// Make sure the networks of the CCE cluster and the CSE engine are same.
			acceptance.TestAccPreCheckCSEMicroserviceEngineID(t)
			acceptance.TestAccPreCheckImsImageUrl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV3Components_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "components.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_component_id_set_and_valid", "true"),
					resource.TestCheckOutput("is_environment_id_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_name_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_runtime_stack_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_source_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_version_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_refer_resources_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_build_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_tags_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_status_set", "true"),
					resource.TestCheckOutput("is_component_created_at_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_updated_at_set_and_valid", "true"),
				),
			},
		},
	})
}

func testAccDataV3Components_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_servicestagev3_components" "test" {
  depends_on = [
    huaweicloud_servicestagev3_component.test
  ]

  application_id = huaweicloud_servicestagev3_application.test.id
}

locals {
  component_id            = huaweicloud_servicestagev3_component.test.id
  component_filter_result = try([
    for v in data.huaweicloud_servicestagev3_components.test.components : v if v.id == local.component_id
  ][0], null)
}

output "is_component_id_set_and_valid" {
  value = local.component_filter_result != null
}

output "is_environment_id_set_and_valid" {
  value = try(local.component_filter_result.environment_id == huaweicloud_servicestagev3_environment.test.id, false)
}

output "is_component_name_set_and_valid" {
  value = try(local.component_filter_result.name == huaweicloud_servicestagev3_component.test.name, false)
}

output "is_component_runtime_stack_set_and_valid" {
  value = try(length(local.component_filter_result.runtime_stack) > 0 && alltrue([
	  local.component_filter_result.runtime_stack[0].name == huaweicloud_servicestagev3_component.test.runtime_stack[0].name,
	  local.component_filter_result.runtime_stack[0].type == huaweicloud_servicestagev3_component.test.runtime_stack[0].type,
	  local.component_filter_result.runtime_stack[0].deploy_mode == huaweicloud_servicestagev3_component.test.runtime_stack[0].deploy_mode,
	  local.component_filter_result.runtime_stack[0].version == huaweicloud_servicestagev3_component.test.runtime_stack[0].version,
  ]))
}

output "is_component_source_set_and_valid" {
  value = try(local.component_filter_result.source == huaweicloud_servicestagev3_component.test.source, false)
}

output "is_component_version_set_and_valid" {
  value = try(local.component_filter_result.version == huaweicloud_servicestagev3_component.test.version, false)
}

output "is_component_refer_resources_set_and_valid" {
  value = try(length(local.component_filter_result.refer_resources) == length(huaweicloud_servicestagev3_component.test.refer_resources), false)
}

output "is_component_build_set_and_valid" {
  value = try(local.component_filter_result.build == huaweicloud_servicestagev3_component.test.build, false)
}

output "is_component_tags_set_and_valid" {
  value = try(length(setintersection(matchkeys(values(huaweicloud_servicestagev3_component.test.tags),
    keys(huaweicloud_servicestagev3_component.test.tags), keys(local.component_filter_result.tags)),
    values(local.component_filter_result.tags))) == length(local.component_filter_result.tags), false)
}

output "is_component_status_set" {
  value = local.component_filter_result.status != ""
}

output "is_component_created_at_set_and_valid" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.component_filter_result.created_at)) > 0, false)
}

output "is_component_updated_at_set_and_valid" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.component_filter_result.updated_at)) > 0, false)
}
`, testAccV3Component_basic_step1(name))
}
