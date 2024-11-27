package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccGaussDBRestore_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBRestoreonfig_basic(name),
			},
		},
	})
}

func TestAccGaussDBRestore_withTimestamp(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBRestoreonfig_withTimestamp(name),
			},
		},
	})
}

func testAccGaussDBRestoreonfig_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_gaussdb_mysql_instance" "instance" {
  count                    = 2
  name                     = "%[2]s_${count.index}"
  password                 = "Test@12345678"
  flavor                   = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  enterprise_project_id    = "0"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  availability_zone_mode   = "multi"
}

resource "huaweicloud_gaussdb_mysql_backup" "backup" {
  instance_id = huaweicloud_gaussdb_mysql_instance.instance[0].id
  name        = "%[2]s_backup"
}

resource "huaweicloud_gaussdb_mysql_restore" "test" {
  target_instance_id = huaweicloud_gaussdb_mysql_instance.instance[1].id
  source_instance_id = huaweicloud_gaussdb_mysql_instance.instance[0].id
  type               = "backup"
  backup_id          = huaweicloud_gaussdb_mysql_backup.backup.id
}
`, common.TestBaseNetwork(name), name)
}

func testAccGaussDBRestoreonfig_withTimestamp(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_gaussdb_mysql_instance" "instance" {
  count                    = 2
  name                     = "%[2]s_${count.index}"
  password                 = "Test@12345678"
  flavor                   = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  enterprise_project_id    = "0"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  availability_zone_mode   = "multi"
}

resource "huaweicloud_gaussdb_mysql_backup" "backup" {
  instance_id = huaweicloud_gaussdb_mysql_instance.instance[0].id
  name        = "%[2]s_backup"
}

data "huaweicloud_gaussdb_mysql_restore_time_ranges" "restore_times" {
  depends_on  = [huaweicloud_gaussdb_mysql_backup.backup]
  instance_id = huaweicloud_gaussdb_mysql_instance.instance[0].id
}

resource "huaweicloud_gaussdb_mysql_restore" "test" {
  target_instance_id = huaweicloud_gaussdb_mysql_instance.instance[1].id
  source_instance_id = huaweicloud_gaussdb_mysql_instance.instance[0].id
  type               = "timestamp"
  restore_time       = data.huaweicloud_gaussdb_mysql_restore_time_ranges.restore_times.restore_times[0].start_time
}
`, common.TestBaseNetwork(name), name)
}
