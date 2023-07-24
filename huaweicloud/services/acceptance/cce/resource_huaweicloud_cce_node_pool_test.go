package cce

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodepools"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCCENodePool_basic(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "huaweicloud_cce_node_pool.test"
	// clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENodePoolImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"initial_node_count", "extend_params",
				},
			},
			{
				Config: testAccCCENodePool_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "scall_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", "9"),
					resource.TestCheckResourceAttr(resourceName, "scale_down_cooldown_time", "100"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
				),
			},
			{
				Config: testAccCCENodePool_volume_extendParams(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.extend_params.test_key", "test_val"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.0.extend_params.test_key1", "test_val1"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.1.extend_params.test_key2", "test_val2"),
				),
			},
		},
	})
}

func TestAccCCENodePool_tagsLabelsTaints(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node_pool.test"
	// clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_tagsLabelsTaints(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.test1", "val1"),
					resource.TestCheckResourceAttr(resourceName, "tags.test2", "val2"),
					resource.TestCheckResourceAttr(resourceName, "labels.test1", "val1"),
					resource.TestCheckResourceAttr(resourceName, "labels.test2", "val2"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.value", "test_value"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.effect", "NoSchedule"),
				),
			},
			{
				Config: testAccCCENodePool_tagsLabelsTaints_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.test1", "val1_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.test2_update", "val2_update"),
					resource.TestCheckResourceAttr(resourceName, "labels.test1", "val1_update"),
					resource.TestCheckResourceAttr(resourceName, "labels.test2_update", "val2_update"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.value", "test_value_update"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(resourceName, "taints.1.key", "test_key_update"),
					resource.TestCheckResourceAttr(resourceName, "taints.1.value", "test_value_update"),
					resource.TestCheckResourceAttr(resourceName, "taints.1.effect", "NoSchedule"),
				),
			},
		},
	})
}

func TestAccCCENodePool_volume_encryption(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node_pool.test"
	// clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_volume_encryption(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "root_volume.0.kms_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "data_volumes.0.kms_key_id"),
				),
			},
		},
	})
}

func TestAccCCENodePool_prePaid(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node_pool.test"
	// clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "1"),
				),
			},
		},
	})
}

func TestAccCCENodePool_SecurityGroups(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node_pool.test"
	// clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_SecurityGroups(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups.0",
						"huaweicloud_networking_secgroup.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups.1",
						"huaweicloud_networking_secgroup.test.1", "id"),
				),
			},
		},
	})
}

func TestAccCCENodePool_serverGroup(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node_pool.test"
	// clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_serverGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "ecs_group_id",
						"huaweicloud_compute_servergroup.test", "id"),
				),
			},
		},
	})
}

func TestAccCCENodePool_storage(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node_pool.test"
	// clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_storage(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "storage.0.selectors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "storage.0.groups.#"),
				),
			},
		},
	})
}

func testAccCheckCCENodePoolDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	cceClient, err := config.CceV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	var clusterId string
	var nodepollId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "huaweicloud_cce_cluster" {
			clusterId = rs.Primary.ID
		}

		if rs.Type == "huaweicloud_cce_node_pool" {
			nodepollId = rs.Primary.ID
		}

		if clusterId == "" || nodepollId == "" {
			continue
		}

		_, err := nodepools.Get(cceClient, clusterId, nodepollId).Extract()
		if err == nil {
			return fmtp.Errorf("Node still exists")
		}
	}

	return nil
}

func testAccCCENodePoolImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		cluster, ok := s.RootModule().Resources["huaweicloud_cce_cluster.test"]
		if !ok {
			return "", fmtp.Errorf("Cluster not found: %s", cluster)
		}
		nodePool, ok := s.RootModule().Resources["huaweicloud_cce_node_pool.test"]
		if !ok {
			return "", fmtp.Errorf("Node pool not found: %s", nodePool)
		}
		if cluster.Primary.ID == "" || nodePool.Primary.ID == "" {
			return "", fmtp.Errorf("resource not found: %s/%s", cluster.Primary.ID, nodePool.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", cluster.Primary.ID, nodePool.Primary.ID), nil
	}
}

func testAccCheckCCENodePoolExists(n string, cluster string, nodePool *nodepools.NodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}
		c, ok := s.RootModule().Resources[cluster]
		if !ok {
			return fmtp.Errorf("Cluster not found: %s", c)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}
		if c.Primary.ID == "" {
			return fmtp.Errorf("Cluster id is not set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		cceClient, err := config.CceV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
		}

		found, err := nodepools.Get(cceClient, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmtp.Errorf("Node Pool not found")
		}

		*nodePool = *found

		return nil
	}
}

func testAccCCENodePool_Base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_compute_keypair" "test" {
  name = "%[2]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, common.TestVpc(rName), rName)
}

func testAccCCENodePool_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  extend_params {
    docker_base_size = 20
    postinstall      = <<EOF
#! /bin/bash
date
EOF
  }
}
`, testAccCCENodePool_Base(rName), rName)
}

func testAccCCENodePool_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 2
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = true
  min_node_count           = 2
  max_node_count           = 9
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  extend_params {
    docker_base_size = 20
    postinstall      = <<EOF
#! /bin/bash
date
EOF
  }

  lifecycle {
    ignore_changes = [
      extend_params
    ]
  }
}
`, testAccCCENodePool_Base(rName), updateName)
}

func testAccCCENodePool_volume_extendParams(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size          = 40
    volumetype    = "SSD"
    extend_params = {
      test_key = "test_val"
    }
  }

  data_volumes {
    size          = 100
    volumetype    = "SSD"
    extend_params = {
      test_key1 = "test_val1"
    }
  }

  data_volumes {
    size          = 100
    volumetype    = "SSD"
    extend_params = {
      test_key2 = "test_val2"
    }
  }

  extend_params {
    docker_base_size = 20
    postinstall      = <<EOF
#! /bin/bash
date
EOF
  }

  lifecycle {
    ignore_changes = [
      extend_params
    ]
  }
}
`, testAccCCENodePool_Base(rName), rName)
}

func testAccCCENodePool_volume_encryption(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
    kms_key_id = huaweicloud_kms_key.test.id
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = huaweicloud_kms_key.test.id
  }
}
`, testAccCCENodePool_Base(rName), rName, rName)
}

func testAccCCENodePool_SecurityGroups(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_compute_keypair" "test" {
  name       = "%[2]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_networking_secgroup" "test" {
  count                 = 4
  name                 = "%[2]s-secgroup-${count.index}"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "rule1" {
  security_group_id = huaweicloud_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = huaweicloud_vpc_subnet.eni_test.cidr
}

resource "huaweicloud_networking_secgroup_rule" "rule2" {
  security_group_id = huaweicloud_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "10250"
  protocol          = "tcp"
  remote_ip_prefix  = huaweicloud_vpc_subnet.test.cidr
}

resource "huaweicloud_networking_secgroup_rule" "rule3" {
  security_group_id = huaweicloud_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "30000-32767"
  protocol          = "udp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "rule4" {
  security_group_id = huaweicloud_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "30000-32767"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "rule5" {
  security_group_id = huaweicloud_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "rule6" {
  security_group_id = huaweicloud_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_group_id   = huaweicloud_networking_secgroup.test[0].id
}

resource "huaweicloud_networking_secgroup_rule" "rule7" {
  security_group_id = huaweicloud_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "22"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "CentOS 7.6"
  flavor_id                = "c7.large.2"
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  security_groups = [
    huaweicloud_networking_secgroup.test[0].id,
    huaweicloud_networking_secgroup.test[1].id
  ]

  pod_security_groups = [
    huaweicloud_networking_secgroup.test[2].id,
    huaweicloud_networking_secgroup.test[3].id
  ]

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccCluster_turbo(rName, 1), rName)
}

func testAccCCENodePool_tagsLabelsTaints(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  tags = {
	test1 = "val1"
	test2 = "val2"
  }

  labels = {
	test1 = "val1"
	test2 = "val2"
  }

  taints {
	key    = "test_key"
	value  = "test_value"
	effect = "NoSchedule"
  }

}
`, testAccCCENodePool_Base(rName), rName)
}

func testAccCCENodePool_tagsLabelsTaints_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  tags = {
	test1        = "val1_update"
	test2_update = "val2_update"
  }

  labels = {
	test1        = "val1_update"
	test2_update = "val2_update"
  }

  taints {
	key    = "test_key"
	value  = "test_value_update"
	effect = "NoSchedule"
  }

  taints {
	key    = "test_key_update"
	value  = "test_value_update"
	effect = "NoSchedule"
  }
}
`, testAccCCENodePool_Base(rName), rName)
}

func testAccCCENodePool_prePaid(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccCCENodePool_Base(rName), rName)
}

func testAccCCENodePool_serverGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_servergroup" "test" {
  name     = "%[2]s"
  policies = ["anti-affinity"]
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  ecs_group_id             = huaweicloud_compute_servergroup.test.id

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccCCENodePool_Base(rName), rName)
}

func testAccCCENodePool_storage(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_compute_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = huaweicloud_kms_key.test.id
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = huaweicloud_kms_key.test.id
  }

  storage {
    selectors {
      name              = "cceUse"
      type              = "evs"
      match_label_size  = "100"
      match_label_count = "1"
    }

    selectors {
      name                           = "user"
      type                           = "evs"
      match_label_size               = "100"
      match_label_metadata_encrypted = "1"
      match_label_metadata_cmkid     = huaweicloud_kms_key.test.id
      match_label_count              = "1"
    }

    groups {
      name           = "vgpaas"
      selector_names = ["cceUse"]
      cce_managed    = true

      virtual_spaces {
        name        = "kubernetes"
        size        = "10%%"
        lvm_lv_type = "linear"
      }

      virtual_spaces {
        name        = "runtime"
        size        = "90%%"
      }
    }

    groups {
      name           = "vguser"
      selector_names = ["user"]

      virtual_spaces {
        name        = "user"
        size        = "100%%"
        lvm_lv_type = "linear"
        lvm_path    = "/workspace"
      }
    }
  }
}
`, testAccCCENodePool_Base(rName), rName)
}
