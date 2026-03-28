package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRemoteDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_remote_databases.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRemoteDatabases_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.character_set"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.state"),
				),
			},
		},
	})
}

func testDataSourceRemoteDatabases_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = data.huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2022_SE"
  instance_mode = "single"
}

resource "huaweicloud_rds_instance" "test" {
  depends_on = [huaweicloud_networking_secgroup_rule.ingress]

  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id

  db {
    password = "Terraform145@"
    type     = "SQLServer"
    version  = "2022_SE"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testDataSourceRemoteDatabases_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_remote_databases" "test" {
  instance_id         = huaweicloud_rds_instance.test.id
  server_ip           = huaweicloud_rds_instance.test.fixed_ip
  server_port         = huaweicloud_rds_instance.test.db[0].port
  login_user_name     = "rdsuser"
  login_user_password = "Terraform145@"
}
`, testDataSourceRemoteDatabases_base(name))
}
