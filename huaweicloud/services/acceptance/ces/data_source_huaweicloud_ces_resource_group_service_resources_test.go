package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceCesGroupServiceResources_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_resource_group_service_resources.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesGroupServiceResources_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.dimensions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.dimensions.0.value"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_dim_name_filter_useful", "true"),
					resource.TestCheckOutput("is_dim_value_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesGroupServiceResources_basic() string {
	return fmt.Sprintf(`
%[1]s

locals {
  dim_name  = "kafka_instance_id"
  dim_value = huaweicloud_dms_rocketmq_instance.test.id
}

data "huaweicloud_ces_resource_group_service_resources" "test" {
  group_id = huaweicloud_ces_resource_group.test.id
  service  = "SYS.DMS" 

  depends_on = [
    huaweicloud_ces_resource_group.test,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_ces_resource_group_service_resources.test.resources) >= 1
}

data "huaweicloud_ces_resource_group_service_resources" "filter_by_dim_name" {
  group_id = huaweicloud_ces_resource_group.test.id
  service  = "SYS.DMS" 
  dim_name = local.dim_name

  depends_on = [
    huaweicloud_ces_resource_group.test,
  ]
}

output "is_dim_name_filter_useful" {
  value = length(data.huaweicloud_ces_resource_group_service_resources.filter_by_dim_name.resources) >= 1 && alltrue(
    [for r in data.huaweicloud_ces_resource_group_service_resources.filter_by_dim_name.resources[*]: r.dimensions[0].name == local.dim_name]
  )
}

data "huaweicloud_ces_resource_group_service_resources" "filter_by_dim_value" {
  group_id  = huaweicloud_ces_resource_group.test.id
  service   = "SYS.DMS" 
  dim_value = local.dim_value

  depends_on = [
    huaweicloud_ces_resource_group.test,
  ]
}

output "is_dim_value_filter_useful" {
  value = length(data.huaweicloud_ces_resource_group_service_resources.filter_by_dim_value.resources) >= 1 && alltrue(
    [for r in data.huaweicloud_ces_resource_group_service_resources.filter_by_dim_value.resources[*]: r.dimensions[0].value == local.dim_value]
  )
}
`, testDataSourceCesGroupServices_base())
}

func testDataSourceCesGroupServices_base() string {
	name := acceptance.RandomAccResourceName()
	baseNetwork := common.TestBaseNetwork(name)
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

%[2]s

%[3]s

resource "huaweicloud_ces_resource_group" "test" {
  name = "%[4]s"
  
  resources {
    namespace = "SYS.DMS"
    dimensions {
      name  = "kafka_instance_id"
      value = huaweicloud_dms_kafka_instance.test.id
    }
  }

  resources {
    namespace = "SYS.DMS"
    dimensions {
      name  = "reliablemq_instance_id"
      value = huaweicloud_dms_rocketmq_instance.test.id
    }
  }

  depends_on = [
    huaweicloud_dms_kafka_instance.test,
    huaweicloud_dms_rocketmq_instance.test,
  ]
}
`, baseNetwork, testKafkaInstance(name), testRocketMQInstance(name), name)
}

func testKafkaInstance(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  kafka_flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%[1]s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  flavor_id         = local.kafka_flavor.id
  storage_spec_code = local.kafka_flavor.ios[0].storage_spec_code
  
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  engine_version = "2.7"
  storage_space  = local.kafka_flavor.properties[0].min_broker * local.kafka_flavor.properties[0].min_storage_per_node
  broker_num     = 3

  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, name)
}

func testRocketMQInstance(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.huaweicloud_dms_rocketmq_flavors.test
  flavor        = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = local.query_results.versions[0]
  storage_space     = 300
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2],
  ]

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 1
  enable_acl        = true

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`, name)
}
