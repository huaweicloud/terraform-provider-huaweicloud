package elb

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/elb/v3/logtanks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func getELBLogTankResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}
	return logtanks.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3LogTanks_basic(t *testing.T) {
	var logTanks logtanks.LogTank
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_logtanks.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&logTanks,
		getELBLogTankResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3LogTanksDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LogTanksConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_topic_id",
						"huaweicloud_lts_stream.test", "id"),
				),
			},
			{
				Config: testAccElbV3LogTanksConfig_update(rNameUpdate),
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

func testAccCheckElbV3LogTanksDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_elb_logtanks" {
			continue
		}

		_, err := logtanks.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("LogStanks still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccElbV3LogTanksConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}
data "huaweicloud_availability_zones" "test" {}
resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.subnet_id
  ipv6_network_id = data.huaweicloud_vpc_subnet.test.id
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
resource "huaweicloud_elb_logtanks" "test" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  log_group_id    = huaweicloud_lts_group.test.id
  log_topic_id    = huaweicloud_lts_stream.test.id
}
`, rName)
}

func testAccElbV3LogTanksConfig_update(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}
data "huaweicloud_availability_zones" "test" {}
resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.subnet_id
  ipv6_network_id = data.huaweicloud_vpc_subnet.test.id
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
resource "huaweicloud_elb_logtanks" "test" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  log_group_id    = huaweicloud_lts_group.test_update.id
  log_topic_id    = huaweicloud_lts_stream.test_update.id
}
`, rName)
}
