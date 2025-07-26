package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsMigrationTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_migration_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcsMigrationTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.task_name"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.migration_type"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.migration_method"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.ecs_tenant_private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.data_source"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.source_instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.source_instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.source_instance_subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.source_instance_spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.source_instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.source_instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.target_instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.target_instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.target_instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.target_instance_subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.target_instance_addrs"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.target_instance_spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.resume_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.supported_features.#"),
					resource.TestCheckResourceAttrSet(dataSource, "migration_tasks.0.created_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcsMigrationTasks_base(name string) string {
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
  description        = "test migration task"

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
`, name)
}

func testDataSourceDcsMigrationTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dcs_migration_tasks" "test" {
  depends_on = [huaweicloud_dcs_online_data_migration_task.test]
}

data "huaweicloud_dcs_migration_tasks" "name_filter" {
  depends_on = [huaweicloud_dcs_online_data_migration_task.test]

  name = "%[2]s"
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dcs_migration_tasks.name_filter.migration_tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_migration_tasks.name_filter.migration_tasks[*].task_name : v == "%[2]s"]  
  )
}
`, testDataSourceDcsMigrationTasks_base(name), name)
}
