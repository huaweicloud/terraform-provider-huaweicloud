package huaweicloud

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingNetworkV2DataSource_basic(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	network := fmt.Sprintf("acc_test_network-%06x", rand.Int31n(1000000))
	cidr := fmt.Sprintf("192.168.%d.0/24", rand.Intn(200))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeprecated(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkV2DataSource_basic(network, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloud_networking_network_v2.net_by_name"),
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloud_networking_network_v2.net_by_id"),
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloud_networking_network_v2.net_by_cidr"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_network_v2.net_by_name", "name", network),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_network_v2.net_by_id", "name", network),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_network_v2.net_by_cidr", "name", network),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_network_v2.net_by_name", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_network_v2.net_by_id", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_network_v2.net_by_cidr", "matching_subnet_cidr", cidr),
				),
			},
		},
	})
}

func testAccCheckNetworkingNetworkV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find network data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Network data source ID not set")
		}

		return nil
	}
}

func testAccNetworkingNetworkV2DataSource_basic(name, cidr string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_network_v2" "net" {
  name = "%s"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet" {
  name = "huaweicloud_test_subnet"
  cidr = "%s"
  no_gateway = true
  network_id = huaweicloud_networking_network_v2.net.id
}

data "huaweicloud_networking_network_v2" "net_by_name" {
	name = huaweicloud_networking_network_v2.net.name
}

data "huaweicloud_networking_network_v2" "net_by_id" {
	network_id = huaweicloud_networking_network_v2.net.id
}

data "huaweicloud_networking_network_v2" "net_by_cidr" {
	matching_subnet_cidr = huaweicloud_networking_subnet_v2.subnet.cidr
}
`, name, cidr)
}
