package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDdsInstanceFlavorUpdate_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDdsInstanceFlavorUpdate_basic(rName),
			},
		},
	})
}

func TestAccDdsInstanceFlavorUpdate_prepaid(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDdsInstanceFlavorUpdate_prpaid(rName),
			},
		},
	})
}

func testAccDdsInstanceFlavorUpdate_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "test" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"
  port              = 8800

  datastore {
    type           = "DDS-Community"
    version        = "5.0"
    storage_engine = "rocksDB"
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
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDdsInstanceFlavorUpdate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance_flavor_update" "test_update_mongos_flavor" {
  instance_id      = huaweicloud_dds_instance.test.id
  target_spec_code = "dds.mongodb.c6.2xlarge.4.mongos"
  target_type      = "mongos"
  target_id        = [for v in huaweicloud_dds_instance.test.groups : [for vv in v.nodes : vv.id][0] if v.type == "mongos"][0]
}

resource "huaweicloud_dds_instance_flavor_update" "test_update_shard_flavor" {
  instance_id      = huaweicloud_dds_instance.test.id
  target_spec_code = "dds.mongodb.s6.medium.4.shard"
  target_type      = "shard"
  target_id        = [for v in huaweicloud_dds_instance.test.groups : v if v.type == "shard"][0].id
}
`, testAccDdsInstanceFlavorUpdate_base(rName))
}

func testAccDdsInstanceFlavorUpdate_prpaid_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "test" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"
  port              = 8800

  datastore {
    type           = "DDS-Community"
    version        = "5.0"
    storage_engine = "rocksDB"
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

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDdsInstanceFlavorUpdate_prpaid(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance_flavor_update" "test_update_mongos_flavor" {
  instance_id      = huaweicloud_dds_instance.test.id
  target_spec_code = "dds.mongodb.s6.large.4.mongos"
  target_type      = "mongos"
  target_id        = [for v in huaweicloud_dds_instance.test.groups : [for vv in v.nodes : vv.id][0] if v.type == "mongos"][0]
}

resource "huaweicloud_dds_instance_flavor_update" "test_update_shard_flavor" {
  instance_id      = huaweicloud_dds_instance.test.id
  target_spec_code = "dds.mongodb.s6.medium.4.shard"
  target_type      = "shard"
  target_id        = [for v in huaweicloud_dds_instance.test.groups : v if v.type == "shard"][0].id
}
`, testAccDdsInstanceFlavorUpdate_prpaid_base(rName))
}
