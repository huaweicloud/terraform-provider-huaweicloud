package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerClustersPolicyTemplates_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_clusters_policy_templates.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContainerClustersPolicyTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.template_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.template_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.target_kind"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.tag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.constraint_template"),

					resource.TestCheckOutput("is_template_name_filter_useful", "true"),
					resource.TestCheckOutput("is_template_type_filter_useful", "true"),
					resource.TestCheckOutput("is_target_kind_filter_useful", "true"),
					resource.TestCheckOutput("is_tag_filter_useful", "true"),
					resource.TestCheckOutput("is_level_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceContainerClustersPolicyTemplates_basic() string {
	return `
data "huaweicloud_hss_container_clusters_policy_templates" "test" {}

# Filter using template_name.
locals {
  template_name = data.huaweicloud_hss_container_clusters_policy_templates.test.data_list[0].template_name
}

data "huaweicloud_hss_container_clusters_policy_templates" "template_name_filter" {
  template_name = local.template_name
}

output "is_template_name_filter_useful" {
  value = length(data.huaweicloud_hss_container_clusters_policy_templates.template_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_clusters_policy_templates.template_name_filter.data_list[*].template_name : v == local.template_name]
  )
}

# Filter using template_type.
locals {
  template_type = data.huaweicloud_hss_container_clusters_policy_templates.test.data_list[0].template_type
}

data "huaweicloud_hss_container_clusters_policy_templates" "template_type_filter" {
  template_type = local.template_type
}

output "is_template_type_filter_useful" {
  value = length(data.huaweicloud_hss_container_clusters_policy_templates.template_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_clusters_policy_templates.template_type_filter.data_list[*].template_type : v == local.template_type]
  )
}

# Filter using target_kind.
locals {
  target_kind = data.huaweicloud_hss_container_clusters_policy_templates.test.data_list[0].target_kind
}

data "huaweicloud_hss_container_clusters_policy_templates" "target_kind_filter" {
  target_kind = local.target_kind
}

output "is_target_kind_filter_useful" {
  value = length(data.huaweicloud_hss_container_clusters_policy_templates.target_kind_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_clusters_policy_templates.target_kind_filter.data_list[*].target_kind : v == local.target_kind]
  )
}

# Filter using tag.
locals {
  tag = data.huaweicloud_hss_container_clusters_policy_templates.test.data_list[0].tag
}

data "huaweicloud_hss_container_clusters_policy_templates" "tag_filter" {
  tag = local.tag
}

output "is_tag_filter_useful" {
  value = length(data.huaweicloud_hss_container_clusters_policy_templates.tag_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_clusters_policy_templates.tag_filter.data_list[*].tag : v == local.tag]
  )
}

# Filter using level.
locals {
  level = data.huaweicloud_hss_container_clusters_policy_templates.test.data_list[0].level
}

data "huaweicloud_hss_container_clusters_policy_templates" "level_filter" {
  level = local.level
}

output "is_level_filter_useful" {
  value = length(data.huaweicloud_hss_container_clusters_policy_templates.level_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_clusters_policy_templates.level_filter.data_list[*].level : v == local.level]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_container_clusters_policy_templates" "enterprise_project_id_filter" {
  enterprise_project_id = "0"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_container_clusters_policy_templates.enterprise_project_id_filter.data_list) > 0
}

# Filter using non existent template_name.
data "huaweicloud_hss_container_clusters_policy_templates" "not_found" {
  template_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_container_clusters_policy_templates.not_found.data_list) == 0
}
`
}
