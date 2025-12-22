package rabbitmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceRabbitMQInstances_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dms_rabbitmq_instances.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byInstanceId   = "data.huaweicloud_dms_rabbitmq_instances.filter_by_instance_id"
		dcByInstanceId = acceptance.InitDataSourceCheck(byInstanceId)

		byEngineVersion   = "data.huaweicloud_dms_rabbitmq_instances.filter_by_engine_version"
		dcByEngineVersion = acceptance.InitDataSourceCheck(byEngineVersion)

		byEnterpriseProjectId   = "data.huaweicloud_dms_rabbitmq_instances.filter_by_enterprise_project_id"
		dcByEnterpriseProjectId = acceptance.InitDataSourceCheck(byEnterpriseProjectId)

		byFlavorId   = "data.huaweicloud_dms_rabbitmq_instances.filter_by_flavor_id"
		dcByFlavorId = acceptance.InitDataSourceCheck(byFlavorId)

		byName   = "data.huaweicloud_dms_rabbitmq_instances.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_dms_rabbitmq_instances.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byType   = "data.huaweicloud_dms_rabbitmq_instances.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceRabbitMQInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'instance_id' parameter.
					dcByInstanceId.CheckResourceExists(),
					resource.TestCheckOutput("instance_id_filter_result", "true"),
					// Filter by 'engine_version' parameter.
					dcByEngineVersion.CheckResourceExists(),
					resource.TestCheckOutput("engine_version_filter_result", "true"),
					// Filter by 'enterprise_project_id' parameter.
					dcByEnterpriseProjectId.CheckResourceExists(),
					resource.TestCheckOutput("enterprise_project_id_filter_result", "true"),
					// Filter by 'flavor_id' parameter.
					dcByFlavorId.CheckResourceExists(),
					resource.TestCheckOutput("flavor_id_filter_result", "true"),
					// Filter by 'name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_result", "true"),
					// Filter by 'status' parameter.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_result", "true"),
					// Filter by 'type' parameter.
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_result", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.access_user"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.availability_zones.#"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.broker_num"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.charging_mode"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.connect_address"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.description"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.engine"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.engine_version"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.extend_times"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.is_logical_volume"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.maintain_begin"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.maintain_end"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.management_connect_address"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.name"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.port"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.security_group_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.security_group_name"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.specification"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.ssl_enable"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.status"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.storage_resource_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.storage_space"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.storage_spec_code"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.tags.%"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.type"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.user_name"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.used_storage_space"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.vpc_name"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.disk_encrypted_enable"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.disk_encrypted_key"),
				),
			},
		},
	})
}

func testAccDatasourceRabbitMQInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_dms_rabbitmq_instances" "all" {
  depends_on = [huaweicloud_dms_rabbitmq_instance.test]
}

# Filter by 'instance_id' parameter.
locals {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
}

data "huaweicloud_dms_rabbitmq_instances" "filter_by_instance_id" {
  instance_id = local.instance_id
}

locals {
  instance_id_filter_result = [for v in data.huaweicloud_dms_rabbitmq_instances.filter_by_instance_id.instances[*].id :
    v == local.instance_id
  ]
}

output "instance_id_filter_result" {
  value = length(local.instance_id_filter_result) > 0 && alltrue(local.instance_id_filter_result)
}

# Filter by 'engine_version' parameter.
locals {
  engine_version = huaweicloud_dms_rabbitmq_instance.test.engine_version
}

data "huaweicloud_dms_rabbitmq_instances" "filter_by_engine_version" {
  engine_version = local.engine_version
  depends_on     = [huaweicloud_dms_rabbitmq_instance.test]
}

locals {
  engine_version_filter_result = [for v in data.huaweicloud_dms_rabbitmq_instances.filter_by_engine_version.instances[*].engine_version :
    v == local.engine_version
  ]
}

output "engine_version_filter_result" {
  value = length(local.engine_version_filter_result) > 0 && alltrue(local.engine_version_filter_result)
}

# Filter by 'enterprise_project_id' parameter.
locals {
  enterprise_project_id = huaweicloud_dms_rabbitmq_instance.test.enterprise_project_id
}

data "huaweicloud_dms_rabbitmq_instances" "filter_by_enterprise_project_id" {
  enterprise_project_id = huaweicloud_dms_rabbitmq_instance.test.enterprise_project_id
}

locals {
  enterprise_project_id_filter_result = [
    for v in data.huaweicloud_dms_rabbitmq_instances.filter_by_enterprise_project_id.instances[*].enterprise_project_id :
    v == local.enterprise_project_id
  ]
}

output "enterprise_project_id_filter_result" {
  value = length(local.enterprise_project_id_filter_result) > 0 && alltrue(local.enterprise_project_id_filter_result)
}

# Filter by 'flavor_id' parameter.
locals {
  flavor_id = huaweicloud_dms_rabbitmq_instance.test.flavor_id
}

data "huaweicloud_dms_rabbitmq_instances" "filter_by_flavor_id" {
  flavor_id = huaweicloud_dms_rabbitmq_instance.test.flavor_id
}

locals {
  flavor_id_filter_result = [for v in data.huaweicloud_dms_rabbitmq_instances.filter_by_flavor_id.instances[*].flavor_id :
    v == local.flavor_id
  ]
}

output "flavor_id_filter_result" {
  value = length(local.flavor_id_filter_result) > 0 && alltrue(local.flavor_id_filter_result)
}

# Filter by 'name' parameter.
locals {
  name = huaweicloud_dms_rabbitmq_instance.test.name
}

data "huaweicloud_dms_rabbitmq_instances" "filter_by_name" {
  name             = huaweicloud_dms_rabbitmq_instance.test.name
  exact_match_name = true

  depends_on = [huaweicloud_dms_rabbitmq_instance.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_dms_rabbitmq_instances.filter_by_name.instances[*].name :
    v == local.name
  ]
}

output "name_filter_result" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'status' parameter.
locals {
  status = huaweicloud_dms_rabbitmq_instance.test.status
}

data "huaweicloud_dms_rabbitmq_instances" "filter_by_status" {
  status = huaweicloud_dms_rabbitmq_instance.test.status
}

locals {
  status_filter_result = [for v in data.huaweicloud_dms_rabbitmq_instances.filter_by_status.instances[*].status :
    v == local.status
  ]
}

output "status_filter_result" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by 'type' parameter.
locals {
  type = huaweicloud_dms_rabbitmq_instance.test.type
}

data "huaweicloud_dms_rabbitmq_instances" "filter_by_type" {
  type = huaweicloud_dms_rabbitmq_instance.test.type
}

locals {
  type_filter_result = [for v in data.huaweicloud_dms_rabbitmq_instances.filter_by_type.instances[*].type :
    v == local.type
  ]
}

output "type_filter_result" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}
`, testAccDmsRabbitmqInstance_newFormat_cluster(name))
}
