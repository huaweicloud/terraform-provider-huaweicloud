package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
)

func TestAccCCEClusterV3_basic(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster_v3.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "x509"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCCEClusterV3_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "new description"),
				),
			},
		},
	})
}

func TestAccCCEClusterV3_withEip(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster_v3.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_withEip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"eip",
				},
			},
		},
	})
}

func testAccCheckCCEClusterV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cce_cluster_v3" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckCCEClusterV3Exists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCCEClusterV3_Base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet_v1" "test" {
  name          = "%s"
  cidr          = "192.168.0.0/16"
  gateway_ip    = "192.168.0.1"

  //dns is required for cce node installing
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = huaweicloud_vpc_v1.test.id
}
`, rName, rName)
}

func testAccCCEClusterV3_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster_v3" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc_v1.test.id
  subnet_id              = huaweicloud_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster_v3" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc_v1.test.id
  subnet_id              = huaweicloud_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
  description            = "new description"
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_withEip(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip_v1" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cce_cluster_v3" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc_v1.test.id
  subnet_id              = huaweicloud_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = huaweicloud_vpc_eip_v1.test.address
}
`, testAccCCEClusterV3_Base(rName), rName)
}
