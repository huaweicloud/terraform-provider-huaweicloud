package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/logtanks"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getELBLogTankResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud ELB v3 client: %s", err)
	}
	return logtanks.Get(client, state.Primary.ID).Extract()
}

func TestAccElbLogTank_basic(t *testing.T) {
	var logTanks logtanks.LogTank
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_logtank.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&logTanks,
		getELBLogTankResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbLogTankConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_topic_id",
						"huaweicloud_lts_stream.test", "id"),
				),
			},
			{
				Config: testAccElbLogTankConfig_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id",
						"huaweicloud_lts_group.test_update", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_topic_id",
						"huaweicloud_lts_stream.test_update", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccElbLogTankConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}
resource "huaweicloud_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}
resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id = huaweicloud_vpc_subnet.test.id
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}
resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}
resource "huaweicloud_elb_logtank" "test" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  log_group_id    = huaweicloud_lts_group.test.id
  log_topic_id    = huaweicloud_lts_stream.test.id
}
`, rName)
}

func testAccElbLogTankConfig_update(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}
resource "huaweicloud_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}
resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id = huaweicloud_vpc_subnet.test.id
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}
resource "huaweicloud_lts_group" "test_update" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}
resource "huaweicloud_lts_stream" "test_update" {
  group_id    = huaweicloud_lts_group.test_update.id
  stream_name = "%[1]s"
}
resource "huaweicloud_elb_logtank" "test" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  log_group_id    = huaweicloud_lts_group.test_update.id
  log_topic_id    = huaweicloud_lts_stream.test_update.id
}
`, rName)
}
