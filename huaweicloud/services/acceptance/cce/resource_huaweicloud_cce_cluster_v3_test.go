package cce

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCCEClusterV3_basic(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
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
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
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

func TestAccCCEClusterV3_prePaid(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_prePaid(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccCCEClusterV3_prePaid(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccCCEClusterV3_withEip(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_withEip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				Config: testAccCCEClusterV3_withEipUpdate(rName),
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

func TestAccCCEClusterV3_withEpsId(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccCCEClusterV3_turbo(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_turbo(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "eni"),
				),
			},
		},
	})
}

func TestAccCCEClusterV3_HibernateAndAwake(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
			{
				Config: testAccCCEClusterV3_hibernate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Hibernation"),
				),
			},
			{
				Config: testAccCCEClusterV3_awake(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	cceClient, err := config.CceV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cce_cluster" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckCCEClusterV3Exists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		cceClient, err := config.CceV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmtp.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCCEClusterV3_Base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name          = "%s"
  cidr          = "192.168.0.0/24"
  gateway_ip    = "192.168.0.1"

  //dns is required for cce node installing
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = huaweicloud_vpc.test.id
}
`, rName, rName)
}

func testAccCCEClusterV3_prePaid(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "%[3]v"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCEClusterV3_Base(rName), rName, isAutoRenew)
}

func testAccCCEClusterV3_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "new description"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_withEip(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test" {
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

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = huaweicloud_vpc_eip.test.address
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_withEipUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test" {
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

resource "huaweicloud_vpc_eip" "update" {
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

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = huaweicloud_vpc_eip.update.address
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  enterprise_project_id  = "%s"
}

`, testAccCCEClusterV3_Base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCCEClusterV3_turbo(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "eni_test" {
  name          = "%s-eni"
  cidr          = "192.168.2.0/24"
  gateway_ip    = "192.168.2.1"
  vpc_id        = huaweicloud_vpc.test.id
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "eni"
  eni_subnet_id          = huaweicloud_vpc_subnet.eni_test.ipv4_subnet_id
  eni_subnet_cidr        = huaweicloud_vpc_subnet.eni_test.cidr
}

`, testAccCCEClusterV3_Base(rName), rName, rName)
}

func testAccCCEClusterV3_hibernate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  hibernate              = true
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_awake(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  hibernate              = false
}
`, testAccCCEClusterV3_Base(rName), rName)
}
