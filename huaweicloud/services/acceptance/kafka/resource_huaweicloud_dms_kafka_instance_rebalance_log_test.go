package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
)

func getInstanceRebalanceLogResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dmsv2", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return kafka.GetInstanceRebalanceLog(client, state.Primary.ID)
}

func TestAccInstanceRebalanceLog_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dms_kafka_instance_rebalance_log.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getInstanceRebalanceLogResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceRebalanceLog_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_dms_kafka_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "log_group_id"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_id"),
					resource.TestCheckResourceAttrSet(rName, "dashboard_id"),
					resource.TestCheckResourceAttrSet(rName, "log_type"),
					resource.TestCheckResourceAttrSet(rName, "log_file_name"),
					resource.TestCheckResourceAttr(rName, "status", "OPEN"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccInstanceRebalanceLog_logExists(name),
				ExpectError: regexp.MustCompile(`log already exists`),
			},
		},
	})
}

func testAccInstanceRebalanceLog_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type = "cluster.small"
}

locals {
  flavor = try(data.huaweicloud_dms_kafka_flavors.test.flavors[0], {})
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%[2]s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = try(local.flavor.id, null)
  storage_spec_code  = try(local.flavor.ios[0].storage_spec_code, null)
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  engine_version     = "3.x"
  storage_space      = try(local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node, null)
  broker_num         = 3
}`, common.TestBaseNetwork(name), name)
}

func testAccInstanceRebalanceLog_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_instance_rebalance_log" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
}
`, testAccInstanceRebalanceLog_base(name))
}

func testAccInstanceRebalanceLog_logExists(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_instance_rebalance_log" "expect_error" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
}
`, testAccInstanceRebalanceLog_basic(name))
}
