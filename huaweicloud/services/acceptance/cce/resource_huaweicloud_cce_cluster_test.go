package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccCluster_basic(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_users.0.client_certificate_data",
					"certificate_users.0.client_key_data",
					"kube_config_raw",
				},
			},
			{
				Config: testAccCluster_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "new description"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
		},
	})
}

func TestAccCluster_prePaid(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_prePaid(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccCluster_prePaid(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccCluster_withEip(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_withEip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				Config: testAccCluster_withEipUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"eip", "certificate_users.0.client_certificate_data", "kube_config_raw",
				},
			},
		},
	})
}

func TestAccCluster_withEpsId(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccCluster_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccCluster_turbo(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_turbo(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "eni"),
					resource.TestCheckResourceAttr(resourceName, "enable_dist_mgt", "true"),
					resource.TestCheckOutput("is_eni_subnet_id_different", "false"),
				),
			},
			{
				Config: testAccCluster_turbo(rName, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "eni"),
					resource.TestCheckOutput("is_eni_subnet_id_different", "false"),
				),
			},
		},
	})
}

func TestAccCluster_hibernate(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
			{
				Config: testAccCluster_hibernate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Hibernation"),
				),
			},
			{
				Config: testAccCluster_awake(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
		},
	})
}

func TestAccCluster_multiContainerNetworkCidrs(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_multiContainerNetworkCidrs(rName, "172.16.0.0/24,172.16.1.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "vpc-router"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "container_network_cidr", "172.16.0.0/24,172.16.1.0/24"),
				),
			},
			{
				Config: testAccCluster_multiContainerNetworkCidrs(rName, "172.16.0.0/24,172.16.1.0/24,172.16.2.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "vpc-router"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "container_network_cidr", "172.16.0.0/24,172.16.1.0/24,172.16.2.0/24"),
				),
			},
		},
	})
}

func TestAccCluster_secGroup(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_secGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test1", "id"),
				),
			},
			{
				Config: testAccCluster_secGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test2", "id"),
				),
			},
		},
	})
}

func TestAccCluster_resize(t *testing.T) {
	var cluster clusters.Clusters

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_resize(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
				),
			},
			{
				Config: testAccCluster_resize_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.medium"),
				),
			},
		},
	})
}

func TestAccCluster_resizePeriod(t *testing.T) {
	var cluster clusters.Clusters

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_resizePeriod(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
				),
			},
			{
				Config: testAccCluster_resize_updatePeriod(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.medium"),
				),
			},
		},
	})
}

func testAccCheckClusterDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	cceClient, err := config.CceV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating CCE v3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cce_cluster" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("cluster still exists")
		}
	}

	return nil
}

func testAccCheckClusterExists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		cceClient, err := config.CceV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CCE v3 client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCluster_prePaid(rName string, isAutoRenew bool) string {
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
`, common.TestVpc(rName), rName, isAutoRenew)
}

func testAccCluster_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  timezone               = "Asia/Shanghai"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_update(rName string) string {
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
  timezone               = "Asia/Shanghai"

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_withEip(rName string) string {
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
`, common.TestVpc(rName), rName)
}

func testAccCluster_withEipUpdate(rName string) string {
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
`, common.TestVpc(rName), rName)
}

func testAccCluster_withEpsId(rName string) string {
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

`, common.TestVpc(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCluster_turbo(rName string, eniNum int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "eni_test" {
  count      = %[3]d

  name       = "%[2]s-eni-${count.index}"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index + 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index + 1), 1)
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "eni"
  enable_dist_mgt        = true
  eni_subnet_id          = join(",", huaweicloud_vpc_subnet.eni_test[*].ipv4_subnet_id)
}

output "is_eni_subnet_id_different" {
  value = length(setsubtract(split(",", huaweicloud_cce_cluster.test.eni_subnet_id),
  huaweicloud_vpc_subnet.eni_test[*].ipv4_subnet_id)) != 0
}
`, common.TestVpc(rName), rName, eniNum)
}

func testAccCluster_hibernate(rName string) string {
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
`, common.TestVpc(rName), rName)
}

func testAccCluster_awake(rName string) string {
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
`, common.TestVpc(rName), rName)
}

func testAccCluster_multiContainerNetworkCidrs(rName, containerNetworkCidr string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "vpc-router"
  container_network_cidr = "%s"
  service_network_cidr   = "10.248.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName, containerNetworkCidr)
}

func testAccCluster_secGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "test1" {
  name = "%[2]s-secgroup-1"
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform"
  security_group_id      = huaweicloud_networking_secgroup.test1.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_secGroup_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "test1" {
  name = "%[2]s-secgroup-1"
}

resource "huaweicloud_networking_secgroup" "test2" {
  name = "%[2]s-secgroup-2"
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform update"
  security_group_id      = huaweicloud_networking_secgroup.test2.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_resize(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_resize_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.medium"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform update"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_resizePeriod(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_resize_updatePeriod(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.medium"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform update"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}
