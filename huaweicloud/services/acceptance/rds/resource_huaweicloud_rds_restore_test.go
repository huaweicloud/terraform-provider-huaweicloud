package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccRdsInstanceRestore_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceRestoreConfig_basic(name),
			},
		},
	})
}

func TestAccRdsInstanceRestore_withTimestamp(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceRestoreConfig_withTimestamp(name),
			},
		},
	})
}

func testAccRdsInstanceRestoreConfig_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "instance" {
  count                    = 2
  name                   = "%[2]s__${count.index}"
  flavor                 = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id      = huaweicloud_networking_secgroup.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  vpc_id                 = huaweicloud_vpc.test.id
  availability_zone      = [data.huaweicloud_availability_zones.test.names[0]]

  db {
    type    = "MySQL"
    version = "8.0"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_backup" "backup" {
  instance_id = huaweicloud_rds_instance.instance[0].id
  name        = "%[2]s_backup"
}

resource "huaweicloud_rds_restore" "test" {
  target_instance_id = huaweicloud_rds_instance.instance[1].id
  source_instance_id = huaweicloud_rds_instance.instance[0].id
  type               = "backup"
  backup_id          = huaweicloud_rds_backup.backup.id
}
`, common.TestBaseNetwork(name), name)
}

func testAccRdsInstanceRestoreConfig_withTimestamp(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "instance" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]

  db {
    type    = "MySQL"
    version = "8.0"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = "%[3]s"
}

resource "huaweicloud_rds_restore" "test" {
  target_instance_id = huaweicloud_rds_instance.instance.id
  source_instance_id = "%[3]s"
  type               = "timestamp"
  restore_time       = data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].start_time
}
`, common.TestBaseNetwork(name), name, acceptance.HW_RDS_INSTANCE_ID)
}
