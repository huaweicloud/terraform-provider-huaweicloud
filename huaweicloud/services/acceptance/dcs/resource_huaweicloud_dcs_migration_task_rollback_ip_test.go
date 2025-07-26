package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsMigrationTaskRollbackIp_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDcsMigrationTaskRollbackIp_basic(name),
			},
		},
	})
}

func testDcsMigrationTaskRollbackIp_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "6.0"
  capacity       = 1
  name           = "redis.ha.au1.large.r4.1"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_dcs_instance" "test" {
  count = 2

  name               = "%[1]s_${count.index}"
  engine_version     = "6.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}

resource "huaweicloud_dcs_online_data_migration_task" "test" {
  task_name          = "%[1]s"
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  migration_method   = "incremental_migration"
  resume_mode        = "auto"
  bandwidth_limit_mb = "100"

  source_instance {
    id       = huaweicloud_dcs_instance.test[0].id
    password = "Huawei_test"
  }

  target_instance {
    id       = huaweicloud_dcs_instance.test[1].id
    password = "Huawei_test"
  }

  lifecycle {
    ignore_changes = [
      source_instance.0.addrs, target_instance.0.addrs,
    ]
  }
}

resource "huaweicloud_dcs_migration_task_exchange_ip" "test" {
  task_id = huaweicloud_dcs_online_data_migration_task.test.id
}
`, name)
}

func testDcsMigrationTaskRollbackIp_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_migration_task_rollback_ip" "test" {
  depends_on = [huaweicloud_dcs_migration_task_exchange_ip.test]

  task_id = huaweicloud_dcs_online_data_migration_task.test.id
}
`, testDcsMigrationTaskRollbackIp_base(name))
}
