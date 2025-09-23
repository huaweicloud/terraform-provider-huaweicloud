package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccGaussDBMysqlInstanceScheduledTaskCancel_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_gaussdb_mysql_scheduled_task_cancel.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBMysqlInstanceScheduledTaskCancel_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "job_status", "Canceled"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_status"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "job_name"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "end_time"),
					resource.TestCheckResourceAttrSet(resourceName, "datastore_type"),
				),
			},
		},
	})
}

func testAccGaussDBMysqlInstanceScheduledTaskCancel_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_mysql_flavors" "test" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = "multi"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name                     = "%[2]s"
  password                 = "Test@12345678"
  flavor                   = data.huaweicloud_gaussdb_mysql_flavors.test.flavors[0].name
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  enterprise_project_id    = "0"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  availability_zone_mode   = "multi"
  read_replicas            = 2
}

resource "huaweicloud_gaussdb_mysql_instance_restart" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  delay       = true
}

data "huaweicloud_gaussdb_mysql_scheduled_tasks" "test" {
  depends_on = [huaweicloud_gaussdb_mysql_instance_restart.test]
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccGaussDBMysqlInstanceScheduledTaskCancel_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_mysql_scheduled_task_cancel" "test" {
  job_id = data.huaweicloud_gaussdb_mysql_scheduled_tasks.test.tasks[0].job_id
}`, testAccGaussDBMysqlInstanceScheduledTaskCancel_base(rName))
}
