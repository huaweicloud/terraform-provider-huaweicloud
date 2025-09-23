package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsDbLogsShrinking_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsDbLogsShrinking_basic(name),
			},
		},
	})
}

func testAccRdsDbLogsShrinking_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mssql.spec.se.s6.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  charging_mode     = "postPaid"

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2022_SE"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

resource "huaweicloud_rds_sqlserver_database" "test" {
  depends_on  = [huaweicloud_rds_instance.test]
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[3]s"
}
`, testAccRdsInstance_base(), name, name)
}

func testAccRdsDbLogsShrinking_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_database_logs_shrinking" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_sqlserver_database.test.name
}
`, testAccRdsDbLogsShrinking_base(name))
}
