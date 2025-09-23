package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDmsRabbitMQInstances_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_dms_rabbitmq_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDmsRabbitMQInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.access_user"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.availability_zones.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.broker_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.connect_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.engine"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.extend_times"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.is_logical_volume"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.maintain_begin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.maintain_end"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.management_connect_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.security_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.security_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.specification"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.ssl_enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.storage_resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.storage_space"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.storage_spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.tags.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.used_storage_space"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.vpc_name"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("engine_version_filter_is_useful", "true"),
					resource.TestCheckOutput("flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDmsRabbitMQInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rabbitmq_instances" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_instance.test]
}

locals {
  name                  = huaweicloud_dms_rabbitmq_instance.test.name
  id                    = huaweicloud_dms_rabbitmq_instance.test.id
  status                = huaweicloud_dms_rabbitmq_instance.test.status  
  enterprise_project_id = huaweicloud_dms_rabbitmq_instance.test.enterprise_project_id
  engine_version        = huaweicloud_dms_rabbitmq_instance.test.engine_version
  flavor_id             = huaweicloud_dms_rabbitmq_instance.test.flavor_id
  type                  = huaweicloud_dms_rabbitmq_instance.test.type
}

data "huaweicloud_dms_rabbitmq_instances" "name_filter" {
  depends_on       = [huaweicloud_dms_rabbitmq_instance.test]
  name             = local.name
  exact_match_name = "true"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_instances.name_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_instances.name_filter.instances[*].name : v == local.name]
  )
}

data "huaweicloud_dms_rabbitmq_instances" "id_filter" {
  depends_on  = [huaweicloud_dms_rabbitmq_instance.test]
  instance_id = local.id
}
  
output "id_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_instances.id_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_instances.id_filter.instances[*].id : v == local.id]
  )
}

data "huaweicloud_dms_rabbitmq_instances" "status_filter" {
  depends_on = [huaweicloud_dms_rabbitmq_instance.test]
  status     = local.status
}
	
output "status_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_instances.status_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_instances.status_filter.instances[*].status : v == local.status]
  )
}

data "huaweicloud_dms_rabbitmq_instances" "enterprise_project_id_filter" {
  depends_on            = [huaweicloud_dms_rabbitmq_instance.test]
  enterprise_project_id = local.enterprise_project_id
}
	  
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_instances.enterprise_project_id_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_instances.enterprise_project_id_filter.instances[*].enterprise_project_id :
  v == local.enterprise_project_id]
  )
}

data "huaweicloud_dms_rabbitmq_instances" "engine_version_filter" {
  depends_on     = [huaweicloud_dms_rabbitmq_instance.test]
  engine_version = local.engine_version
}
		
output "engine_version_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_instances.engine_version_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_instances.engine_version_filter.instances[*].engine_version : v == local.engine_version]
  )
}

data "huaweicloud_dms_rabbitmq_instances" "flavor_id_filter" {
  depends_on = [huaweicloud_dms_rabbitmq_instance.test]
  flavor_id  = local.flavor_id
}
		  
output "flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_instances.flavor_id_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_instances.flavor_id_filter.instances[*].flavor_id : v == local.flavor_id]
  )
}

data "huaweicloud_dms_rabbitmq_instances" "type_filter" {
  depends_on = [huaweicloud_dms_rabbitmq_instance.test]
  type       = local.type
}
			
output "type_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_instances.type_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_instances.type_filter.instances[*].type : v == local.type]
  )
}

`, testAccDmsRabbitmqInstance_newFormat_cluster(name))
}
