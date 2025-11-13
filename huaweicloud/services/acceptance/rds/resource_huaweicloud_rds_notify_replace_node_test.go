package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNotifyReplaceNode_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccNotifyReplaceNode_basic(rName),
			},
			{
				Config: testAccNotifyReplaceNode_rollback(rName),
			},
		},
	})
}

func testAccNotifyReplaceNode_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_rds_instance" "test" {
  name                = "%[1]s"
  flavor              = "rds.mssql.spec.x1.ee.large.4.ha"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = data.huaweicloud_vpc_subnet.test.id
  vpc_id              = data.huaweicloud_vpc.test.id
  ha_replication_mode = "sync"
  availability_zone   = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[2],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2022_EE"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_read_replica_instance" "test" {
  name                = "%[1]s"
  flavor              = "rds.mssql.spec.n1.ee.xlarge.2.rr"
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  security_group_id   = huaweicloud_networking_secgroup.test.id

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

data "huaweicloud_rds_instances" "test" {
  depends_on = [huaweicloud_rds_read_replica_instance.test]

  type           = "Replica"
  datastore_type = "SQLServer"
}

locals{
  instances           = data.huaweicloud_rds_instances.test.instances
  replica_instance_id = huaweicloud_rds_read_replica_instance.test.id
  replica_node_id     = [for v in local.instances[*]: v if v.id == local.replica_instance_id][0].nodes[0].id
}
`, rName)
}

func testAccNotifyReplaceNode_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_notify_replace_node" "replace" {
  instance_id    = huaweicloud_rds_read_replica_instance.test.id
  node_id        = local.replica_node_id
  replace_action = "REPLACE"
}
`, testAccNotifyReplaceNode_base(rName))
}

func testAccNotifyReplaceNode_rollback(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_notify_replace_node" "rollback" {
  instance_id    = huaweicloud_rds_read_replica_instance.test.id
  node_id        = local.replica_node_id
  replace_action = "REPLACE_ROLLBACK"
}
`, testAccNotifyReplaceNode_base(rName))
}
