package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSQLServerDatabaseCopy_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSQLServerDatabaseCopy_basic(name),
			},
		},
	})
}

func testSQLServerDatabaseCopy_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_EE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id

  db {
    type    = "SQLServer"
    version = "2019_EE"
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "huaweicloud_rds_sqlserver_database" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s_source"
}
`, testAccRdsInstance_base(), name)
}

func testSQLServerDatabaseCopy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_sqlserver_database_copy" "test" {
  depends_on = [huaweicloud_rds_sqlserver_database.test]

  instance_id    = huaweicloud_rds_instance.test.id
  procedure_name = "copy_database"
  db_name_source = "%[2]s_source"
  db_name_target = "%[2]s_target"
}
`, testSQLServerDatabaseCopy_base(name), name)
}
