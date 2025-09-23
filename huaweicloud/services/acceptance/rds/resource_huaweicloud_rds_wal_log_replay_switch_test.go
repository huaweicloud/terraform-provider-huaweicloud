package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceRdsWalLogReplaySwitch_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testResourceRdsWalLogReplaySwitch_basic(rName),
			},
		},
	})
}

func testResourceRdsWalLogReplaySwitch_base(name string) string {
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
}
`, testAccRdsInstance_base(), name)
}

func testResourceRdsWalLogReplaySwitch_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_wal_log_replay_switch" "test" {
  instance_id      = huaweicloud_rds_read_replica_instance.test.id
  pause_log_replay = "true"
}
`, testResourceRdsWalLogReplaySwitch_base(name))
}
