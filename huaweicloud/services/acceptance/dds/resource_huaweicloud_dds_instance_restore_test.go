package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDDSV3InstanceRestore_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	// Avoid CheckDestroy
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceRestore_basic(rName),
			},
		},
	})
}

func testAccDDSInstanceRestore_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%s

%s

%s

resource "huaweicloud_dds_backup" "test" {
  instance_id = huaweicloud_dds_instance.instance2.id
  name        = "%s"
}

resource "huaweicloud_dds_instance_restore" "test" {
  depends_on = [huaweicloud_dds_instance.instance1, huaweicloud_dds_backup.test]

  target_id = huaweicloud_dds_instance.instance1.id
  source_id = huaweicloud_dds_instance.instance2.id
  backup_id = huaweicloud_dds_backup.test.id
}`, common.TestBaseNetwork(rName), testAccDDSInstanceRestoreBase(rName, 1), testAccDDSInstanceRestoreBase(rName, 2), rName)
}

func testAccDDSInstanceRestoreBase(rName string, i int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_instance" "instance%[1]v" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }

  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.shard"
  }

  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }
}`, i, rName)
}
