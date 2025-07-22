package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsInstanceNodeIpRemove_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDcsInstanceNodeIpRemove_baic(name),
			},
		},
	})
}

func testDcsInstanceNodeIpRemove_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 4
  name           = "redis.cluster.xu1.large.r4.4"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 4
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}

data "huaweicloud_dcs_instance_shards" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}
`, name)
}

func testDcsInstanceNodeIpRemove_baic(name string) string {
	return fmt.Sprintf(`
%s

locals {
  replication_list = data.huaweicloud_dcs_instance_shards.test.group_list[0].replication_list
}

resource "huaweicloud_dcs_instance_node_ip_remove" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
  group_id    = data.huaweicloud_dcs_instance_shards.test.group_list[0].group_id
  node_id     = [for v in local.replication_list : v.node_id if v.replication_role == "slave"][0]
}
`, testDcsInstanceNodeIpRemove_base(name))
}
