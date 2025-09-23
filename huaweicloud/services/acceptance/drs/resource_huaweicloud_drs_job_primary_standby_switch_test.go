package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDRSPrimaryStandbySwitch_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dbName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_drs_job.test"
	pwd := "TestDrs@123"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDRSPrimaryStandbySwitch_base(rName, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					// when resource complete, job status can be one of FULL_TRANSFER_STARTED, FULL_TRANSFER_COMPLETE, INCRE_TRANSFER_STARTED
					// primary standby switch can only be updated when job status is INCRE_TRANSFER_STARTED or INCRE_TRANSFER_FAILED
					// wait for job status to be INCRE_TRANSFER_STARTED
					waitForJobStatus(resourceName),
				),
			},
			{
				Config: testAccDRSPrimaryStandbySwitch_basic(rName, dbName, pwd),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccDRSPrimaryStandbySwitch_basic(rName, dbName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_drs_job_primary_standby_switch" "test" {
  job_id = huaweicloud_drs_job.test.id
}`, testAccDRSPrimaryStandbySwitch_base(rName, dbName, pwd))
}

func testAccDRSPrimaryStandbySwitch_base(name, dbName, pwd string) string {
	netConfig := common.TestBaseNetwork(name)
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")
	database1 := testAccRdsMysqlDatabse(dbName, 1)

	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_drs_node_types" "test" {
  engine_type = "cloudDataGuard-mysql"
  type        = "cloudDataGuard"
  direction   = "up"
}

%[3]s

%[4]s

%[5]s

resource "huaweicloud_drs_job" "test" {
  name           = "%[6]s"
  type           = "cloudDataGuard"
  engine_type    = "cloudDataGuard-mysql"
  direction      = "up"
  node_type      = data.huaweicloud_drs_node_types.test.node_types[0]
  net_type       = "eip"
  migration_type = "FULL_INCR_TRANS"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "%[7]s"
    vpc_id      = huaweicloud_rds_instance.test1.vpc_id
    subnet_id   = huaweicloud_rds_instance.test1.subnet_id
  }

  destination_db {
    region      = huaweicloud_rds_instance.test2.region
    ip          = huaweicloud_rds_instance.test2.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%[7]s"
    instance_id = huaweicloud_rds_instance.test2.id
    subnet_id   = huaweicloud_rds_instance.test2.subnet_id
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy, destination_db_readnoly, direction
    ]
  }
}
`, netConfig, testAccSecgroupRule, sourceDb, destDb, database1, name, pwd)
}
