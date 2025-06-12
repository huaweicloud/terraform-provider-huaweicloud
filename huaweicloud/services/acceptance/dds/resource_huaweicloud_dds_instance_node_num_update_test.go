package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDdsInstanceNodeNumUpdate_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDdsInstanceNodeNumUpdate_basic(rName),
			},
		},
	})
}

func TestAccDdsInstanceNodeNumUpdate_prepaid(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDdsInstanceNodeNumUpdate_prpaid(rName),
			},
		},
	})
}

func testAccDdsInstanceNodeNumUpdate_base(rName string) string {
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

func testAccDdsInstanceNodeNumUpdate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance_node_num_update" "test" {
  instance_id = huaweicloud_dds_instance.test.id
  type        = "mongos"
  spec_code   = "dds.mongodb.s6.large.4.mongos"
  num         = "2"
}
`, testAccDdsInstanceNodeNumUpdate_base(rName))
}

func testAccDdsInstanceNodeNumUpdate_prpaid_base(rName string) string {
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

func testAccDdsInstanceNodeNumUpdate_prpaid(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance_node_num_update" "test" {
  instance_id = huaweicloud_dds_instance.test.id
  type        = "shard"
  spec_code   = "dds.mongodb.s6.medium.4.shard"
  num         = "1"

  volume {
    size = "20"
  }
}
`, testAccDdsInstanceNodeNumUpdate_prpaid_base(rName))
}
