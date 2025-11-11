package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsConfigurableSubscriberInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_configurable_subscriber_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsConfigurableSubscriberInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_name"),
					resource.TestCheckOutput("subscriber_instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("subscriber_instance_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsConfigurableSubscriberInstances_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_rds_instance" "test" {
  count = 3

  name              = "%[1]s_${count.index}"
  flavor            = "rds.mssql.spec.x1.se.xlarge.4"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2022_SE"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, name)
}

func testDataSourceRdsConfigurableSubscriberInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_configurable_subscriber_instances" "test" {
  depends_on = [huaweicloud_rds_instance.test[0], huaweicloud_rds_instance.test[1], huaweicloud_rds_instance.test[2]]

  instance_id = huaweicloud_rds_instance.test[0].id
}

data "huaweicloud_rds_configurable_subscriber_instances" "subscriber_instance_id_filter" {
  depends_on = [huaweicloud_rds_instance.test[0], huaweicloud_rds_instance.test[1], huaweicloud_rds_instance.test[2]]

  instance_id            = huaweicloud_rds_instance.test[0].id
  subscriber_instance_id = huaweicloud_rds_instance.test[1].id
}
locals {
  subscriber_instance_id = huaweicloud_rds_instance.test[1].id
  id_filter_instances    = data.huaweicloud_rds_configurable_subscriber_instances.subscriber_instance_id_filter.instances
}
output "subscriber_instance_id_filter_is_useful" {
  value = length(local.id_filter_instances) > 0 && alltrue(
    [for v in local.id_filter_instances[*].instance_id : v == local.subscriber_instance_id]
  )
}

data "huaweicloud_rds_configurable_subscriber_instances" "subscriber_instance_name_filter" {
  depends_on = [huaweicloud_rds_instance.test[0], huaweicloud_rds_instance.test[1], huaweicloud_rds_instance.test[2]]

  instance_id              = huaweicloud_rds_instance.test[0].id
  subscriber_instance_name = huaweicloud_rds_instance.test[2].name
}
locals {
  subscriber_instance_name = huaweicloud_rds_instance.test[2].name
  name_filter_instances    = data.huaweicloud_rds_configurable_subscriber_instances.subscriber_instance_name_filter.instances
}
output "subscriber_instance_name_filter_is_useful" {
  value = length(local.name_filter_instances) > 0 && alltrue(
    [for v in local.name_filter_instances[*].instance_name : v == local.subscriber_instance_name]
  )
}
`, testDataSourceRdsConfigurableSubscriberInstances_base(name))
}
