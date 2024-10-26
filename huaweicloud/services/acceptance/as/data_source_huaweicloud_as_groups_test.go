package as

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceASGroup_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_groups.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_as_groups.name_filter"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byScalingConfigurationID   = "data.huaweicloud_as_groups.scaling_configuration_id_filter"
		dcByScalingConfigurationID = acceptance.InitDataSourceCheck(byScalingConfigurationID)

		byStatus   = "data.huaweicloud_as_groups.status_filter"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEnterpriseProjectID   = "data.huaweicloud_as_groups.enterprise_project_id_filter"
		dcByEnterpriseProjectID = acceptance.InitDataSourceCheck(byEnterpriseProjectID)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare AS group in advance and configure the group ID into the environment variable.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceASGroup_conf,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.availability_zones.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.cool_down_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.current_instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.delete_publicip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.delete_volume"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.desire_instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.health_periodic_audit_grace_period"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.health_periodic_audit_method"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.health_periodic_audit_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.instance_terminate_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.instances.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.is_scaling"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.max_instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.min_instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.multi_az_scaling_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.networks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.networks.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.networks.0.ipv6_enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.networks.0.source_dest_check"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.scaling_configuration_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.scaling_configuration_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.scaling_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.scaling_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.vpc_id"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					dcByScalingConfigurationID.CheckResourceExists(),
					resource.TestCheckOutput("is_scaling_configuration_id_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					dcByEnterpriseProjectID.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceASGroup_conf = `
data "huaweicloud_as_groups" "test" {
}

# Filter by name
locals {
  name = data.huaweicloud_as_groups.test.groups[0].scaling_group_name
}

data "huaweicloud_as_groups" "name_filter" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_as_groups.name_filter.groups[*].scaling_group_name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by scaling_configuration_id
locals {
  scaling_configuration_id = data.huaweicloud_as_groups.test.groups[0].scaling_configuration_id
}

data "huaweicloud_as_groups" "scaling_configuration_id_filter" {
  scaling_configuration_id = local.scaling_configuration_id
}

locals {
  scaling_configuration_id_filter_result = [
    for v in data.huaweicloud_as_groups.scaling_configuration_id_filter.groups[*].scaling_configuration_id : v == local.scaling_configuration_id
  ]
}

output "is_scaling_configuration_id_filter_useful" {
  value = alltrue(local.scaling_configuration_id_filter_result) && length(local.scaling_configuration_id_filter_result) > 0
}

# Filter by status
locals {
  status = data.huaweicloud_as_groups.test.groups[0].status
}

data "huaweicloud_as_groups" "status_filter" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_as_groups.status_filter.groups[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

# Filter by enterprise_project_id
locals {
  enterprise_project_id = data.huaweicloud_as_groups.test.groups[0].enterprise_project_id
}

data "huaweicloud_as_groups" "enterprise_project_id_filter" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  enterprise_project_id_filter_result = [
    for v in data.huaweicloud_as_groups.enterprise_project_id_filter.groups[*].enterprise_project_id : v == local.enterprise_project_id
  ]
}

output "is_enterprise_project_id_filter_useful" {
  value = alltrue(local.enterprise_project_id_filter_result) && length(local.enterprise_project_id_filter_result) > 0
}
`
