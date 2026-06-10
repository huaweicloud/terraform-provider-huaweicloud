package das

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccBatchSetSqlSwitch_basic(t *testing.T) {
	var (
		rName = "huaweicloud_das_batch_set_sql_switch.test"

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchSetSqlSwitch_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "engine_type", "MySQL"),
					resource.TestCheckResourceAttr(rName, "switch_on", "true"),
					resource.TestCheckResourceAttr(rName, "switch_type", "slowsql"),
					resource.TestCheckResourceAttr(rName, "retention_hours", "168"),
				),
			},
			{
				Config: testAccBatchSetSqlSwitch_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "engine_type", "MySQL"),
					resource.TestCheckResourceAttr(rName, "switch_on", "false"),
					resource.TestCheckResourceAttr(rName, "switch_type", "slowsql"),
					resource.TestCheckResourceAttr(rName, "retention_hours", "300"),
				),
			},
		},
	})
}

func TestAccBatchSetSqlSwitch_sqlserver(t *testing.T) {
	var (
		rName = "huaweicloud_das_batch_set_sql_switch.test"

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchSetSqlSwitch_sqlserver(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "engine_type", "SQLServer"),
					resource.TestCheckResourceAttr(rName, "switch_on", "true"),
					resource.TestCheckResourceAttr(rName, "switch_type", "fullsql"),
					resource.TestCheckResourceAttr(rName, "retention_hours", "168"),
				),
			},
			{
				Config: testAccBatchSetSqlSwitch_sqlserver_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "engine_type", "SQLServer"),
					resource.TestCheckResourceAttr(rName, "switch_on", "false"),
					resource.TestCheckResourceAttr(rName, "switch_type", "slowsql"),
					resource.TestCheckResourceAttr(rName, "retention_hours", "300"),
				),
			},
		},
	})
}

func testAccBatchSetSqlSwitch_base(name string) string {
	return fmt.Sprintf(`
%[1]s

# Create a RDS MySQL instance
data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  count = 2  

  name              = "%[2]s-${count.index}"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    type    = "MySQL"
    version = "8.0"
    port    = 3306
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccBatchSetSqlSwitch_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_batch_set_sql_switch" "test" {
  engine_type     = "MySQL"
  switch_on       = true
  switch_type     = "slowsql"
  retention_hours = 168
  instance_ids    = [
    huaweicloud_rds_instance.test[0].id,
	huaweicloud_rds_instance.test[1].id,
  ]
}
`, testAccBatchSetSqlSwitch_base(name))
}

func testAccBatchSetSqlSwitch_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_batch_set_sql_switch" "test" {
  engine_type     = "MySQL"
  switch_on       = false
  switch_type     = "slowsql"
  retention_hours = 300
  instance_ids    = [
    huaweicloud_rds_instance.test[0].id,
	huaweicloud_rds_instance.test[1].id,
  ]

  enable_force_new = "true"
}
`, testAccBatchSetSqlSwitch_base(name))
}

func testAccBatchSetSqlSwitch_sqlserver_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  ports             = 1433  # The port of SQLServer database 
  protocol          = "tcp"
  remote_group_id   = huaweicloud_networking_secgroup.test.id
  action            = "allow"
  priority          = 1
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "sqlserver"
  db_version    = "2017_SE"
  instance_mode = "single"
  group_type    = "general"
  vcpus         = 2
  memory        = 4
}

resource "huaweicloud_rds_instance" "test" {
  count = 2

  name              = "%[2]s-${count.index}"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    type    = "SQLServer"
    version = "2017_SE"
    port    = 1433
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccBatchSetSqlSwitch_sqlserver(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_batch_set_sql_switch" "test" {
  engine_type     = "SQLServer"
  switch_on       = true
  switch_type     = "fullsql"
  retention_hours = 168
  instance_ids    = [
    huaweicloud_rds_instance.test[0].id,
	huaweicloud_rds_instance.test[1].id,
  ]
}
`, testAccBatchSetSqlSwitch_sqlserver_base(name))
}

func testAccBatchSetSqlSwitch_sqlserver_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_batch_set_sql_switch" "test" {
  engine_type     = "SQLServer"
  switch_on       = false
  switch_type     = "slowsql"
  retention_hours = 300
  instance_ids    = [
    huaweicloud_rds_instance.test[0].id,
	huaweicloud_rds_instance.test[1].id,
  ]

  enable_force_new = "true"
}
`, testAccBatchSetSqlSwitch_sqlserver_base(name))
}
