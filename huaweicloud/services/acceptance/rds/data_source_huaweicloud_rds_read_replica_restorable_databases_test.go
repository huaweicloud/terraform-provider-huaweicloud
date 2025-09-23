package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsReadReplicaRestorableDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_read_replica_restorable_databases.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsReadReplicaRestorableDatabases_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "database_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "total_tables"),
					resource.TestCheckResourceAttrSet(dataSource, "table_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.total_tables"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.schemas.#"),
				),
			},
		},
	})
}

func testDataSourceRdsReadReplicaRestorableDatabases_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.medium.2"
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  charging_mode     = "postPaid"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]

  db {
    type     = "PostgreSQL"
    version  = "12"
    password = "Trerraform125@"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_pg_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "test_account_name"
  password    = "Terraform145@!"
}

resource "huaweicloud_rds_pg_database" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "test_database"

  depends_on = [huaweicloud_rds_pg_account.test]
}

resource "huaweicloud_rds_pg_schema" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_pg_database.test.name
  schema_name = "test_schema"
  owner       = huaweicloud_rds_pg_account.test.name

  depends_on = [huaweicloud_rds_pg_database.test]
}

resource "huaweicloud_rds_read_replica_instance" "test" {
  name                = "%[2]s-read-replica"
  flavor              = "rds.pg.n1.medium.2.rr"
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  security_group_id   = data.huaweicloud_networking_secgroup.test.id

  db {
    port = "5432"
  }

  volume {
    type              = "CLOUDSSD"
    size              = 40
    limit_size        = 200
    trigger_threshold = 10
  }

  depends_on = [huaweicloud_rds_pg_schema.test]
}
`, testAccRdsInstance_base(), name)
}

func testDataSourceRdsReadReplicaRestorableDatabases_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_read_replica_restorable_databases" "test" {
  instance_id = huaweicloud_rds_read_replica_instance.test.id
}
`, testDataSourceRdsReadReplicaRestorableDatabases_base(name))
}
