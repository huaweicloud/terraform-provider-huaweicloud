package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV3Components_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_servicestagev3_components.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byApplicationId           = "data.huaweicloud_servicestagev3_components.by_application_id"
		dcByApplicationId         = acceptance.InitDataSourceCheck(byApplicationId)
		byNotFoundApplicationId   = "data.huaweicloud_servicestagev3_components.by_not_found_application_id"
		dcByNotFoundApplicationId = acceptance.InitDataSourceCheck(byNotFoundApplicationId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure at least one of node exist.
			acceptance.TestAccPreCheckCceClusterId(t)
			// Make sure the networks of the CCE cluster and the CSE engine and the ELB loadbalancer are same.
			acceptance.TestAccPreCheckCSEMicroserviceEngineID(t)
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
			// At least one of JAR package need to be provided.
			acceptance.TestAccPreCheckServiceStageJarPkgStorageURLs(t, 1)
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
					resource.TestCheckOutput("is_component_tags_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_status_set", "true"),
					dcByApplicationId.CheckResourceExists(),
					resource.TestMatchResourceAttr(byApplicationId, "components.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_application_id_filter_param_useful", "true"),
					dcByNotFoundApplicationId.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_application_id_filter_param_useful", "true"),
				),
			},
		},
	})
}

func testAccDataV3Components_basic(name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_servicestagev3_components" "all" {
  depends_on = [
    huaweicloud_servicestagev3_component.test
  ]
}

data "huaweicloud_servicestagev3_components" "by_application_id" {
  depends_on = [
    huaweicloud_servicestagev3_component.test
  ]

  application_id = huaweicloud_servicestagev3_application.test.id
}

data "huaweicloud_servicestagev3_components" "by_not_found_application_id" {
  depends_on = [
    huaweicloud_servicestagev3_component.test
  ]

  application_id = "%[2]s"
}

locals {
  component_id                    = huaweicloud_servicestagev3_component.test.id
  manully_filter_component_result = try([
    for v in data.huaweicloud_servicestagev3_components.all.components : v if v.id == local.component_id
  ][0], null)
  application_id_filter_result    = try([
    for v in data.huaweicloud_servicestagev3_components.by_application_id.components : v if v.id == local.component_id
  ][0], null)
}

output "is_component_id_set_and_valid" {
  value = local.manully_filter_component_result.id != null
}

output "is_environment_id_set_and_valid" {
  value = try(local.manully_filter_component_result.environment_id == huaweicloud_servicestagev3_environment.test.id, false)
}

output "is_component_name_set_and_valid" {
  value = try(local.manully_filter_component_result.name == huaweicloud_servicestagev3_component.test.name, false)
}

output "is_component_runtime_stack_set_and_valid" {
  value = try(length(local.manully_filter_component_result.runtime_stack) > 0 && alltrue([
    local.manully_filter_component_result.runtime_stack[0].name == huaweicloud_servicestagev3_component.test.runtime_stack[0].name,
    local.manully_filter_component_result.runtime_stack[0].type == huaweicloud_servicestagev3_component.test.runtime_stack[0].type,
    local.manully_filter_component_result.runtime_stack[0].deploy_mode == huaweicloud_servicestagev3_component.test.runtime_stack[0].deploy_mode,
    local.manully_filter_component_result.runtime_stack[0].version == huaweicloud_servicestagev3_component.test.runtime_stack[0].version,
  ]))
}

output "is_component_source_set_and_valid" {
  value = try(local.manully_filter_component_result.source == huaweicloud_servicestagev3_component.test.source, false)
}

output "is_component_version_set_and_valid" {
  value = try(local.manully_filter_component_result.version == huaweicloud_servicestagev3_component.test.version, false)
}

output "is_component_refer_resources_set_and_valid" {
  value = try(length(local.manully_filter_component_result.refer_resources) == length(huaweicloud_servicestagev3_component.test.refer_resources),
    false)
}

output "is_component_tags_set_and_valid" {
  value = try(length(setintersection(matchkeys(values(huaweicloud_servicestagev3_component.test.tags),
    keys(huaweicloud_servicestagev3_component.test.tags), keys(local.manully_filter_component_result.tags)),
    values(local.manully_filter_component_result.tags))) == length(local.manully_filter_component_result.tags), false)
}

output "is_component_status_set" {
  value = local.manully_filter_component_result.status != ""
}

output "is_application_id_filter_param_useful" {
  value = length(local.application_id_filter_result) > 0
}

output "is_not_found_application_id_filter_param_useful" {
  value = length(data.huaweicloud_servicestagev3_components.by_not_found_application_id.components) == 0
}
`, testAccV3Component_basic_step1(name), randUUID)
}
