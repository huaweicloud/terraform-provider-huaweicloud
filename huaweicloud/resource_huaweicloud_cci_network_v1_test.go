package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cci/v1/networks"
)

func TestAccCCINetworkV1_basic(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCCI(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCINetworkV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCINetworkV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCINetworkV1Exists("huaweicloud_cci_network_v1.net_1", &network),
					resource.TestCheckResourceAttr(
						"huaweicloud_cci_network_v1.net_1", "name", "cci-net"),
				),
			},
		},
	})
}

func testAccCheckCCINetworkV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cciClient, err := config.cciV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cci_network_v1" {
			continue
		}

		_, err := networks.Get(cciClient, "test_ns", rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func testAccCheckCCINetworkV1Exists(n string, network *networks.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		cciClient, err := config.cciV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud CCI client: %s", err)
		}

		found, err := networks.Get(cciClient, "test_ns", rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Name != rs.Primary.ID {
			return fmt.Errorf("Network not found")
		}

		*network = *found

		return nil
	}
}

var testAccCCINetworkV1_basic = fmt.Sprintf(`
resource "huaweicloud_cci_network_v1" "net_1" {
  name = "cci-net"
  namespace = "test-ns"
  security_group = "3b5ceb06-3b8d-43ee-866a-dc0443b85de8"
  project_id = "%s"
  domain_id = "%s"
  vpc_id = "%s"
  network_id = "%s"
  subnet_id = "%s"
  available_zone = "cn-north-1a"
  cidr = "192.168.0.0/24"
}`, OS_TENANT_ID, OS_DOMAIN_ID, OS_VPC_ID, OS_NETWORK_ID, OS_SUBNET_ID)
