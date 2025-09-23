package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPgDatabasePrivilege_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsPgDatabasePrivilege_basic(name),
			},
			{
				Config: testAccRdsPgDatabasePrivilege_basic_update(name),
			},
		},
	})
}

func testAccRdsPgDatabasePrivilege_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  description       = "test_description"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    type    = "PostgreSQL"
    version = "16"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_pg_account" "test" {
  count = 4

  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s_${count.index}"
  password    = "Terraform145@"
}

resource "huaweicloud_rds_pg_database" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = "%[2]s"
  owner         = "root"
  character_set = "UTF8"
  template      = "template1"
  lc_collate    = "en_US.UTF-8"
  lc_ctype      = "en_US.UTF-8"
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsPgDatabasePrivilege_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_database_privilege" "test" {
  depends_on = [
    huaweicloud_rds_pg_database.test,
    huaweicloud_rds_pg_account.test[0],
    huaweicloud_rds_pg_account.test[1],
  ]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_pg_database.test.name

  users {
    name        = huaweicloud_rds_pg_account.test[0].name
    readonly    = true
    schema_name = "public"
  }
  users {
    name        = huaweicloud_rds_pg_account.test[1].name
    readonly    = false
    schema_name = "public"
  }
}
`, testAccRdsPgDatabasePrivilege_base(name))
}

func testAccRdsPgDatabasePrivilege_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_database_privilege" "test" {
  depends_on = [
    huaweicloud_rds_pg_database.test,
    huaweicloud_rds_pg_account.test[0],
    huaweicloud_rds_pg_account.test[1],
  ]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_pg_database.test.name

  users {
    name        = huaweicloud_rds_pg_account.test[2].name
    readonly    = false
    schema_name = "public"
  }
  users {
    name        = huaweicloud_rds_pg_account.test[3].name
    readonly    = true
    schema_name = "public"
  }
}
`, testAccRdsPgDatabasePrivilege_base(name))
}
