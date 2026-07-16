package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodepools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getNodePoolFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CceV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v3 client: %s", err)
	}
	return nodepools.Get(client, state.Primary.Attributes["cluster_id"], state.Primary.ID).Extract()
}

func TestAccNodePool_basic(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		updateName   = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_pool.test"

		baseConfig = testAccNodePool_base(name)

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_basic_step1(name, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.9"),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "1"),
				),
			},
			{
				Config: testAccNodePool_basic_step2(updateName, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.9"),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "scall_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", "9"),
					resource.TestCheckResourceAttr(resourceName, "scale_down_cooldown_time", "100"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNodePoolImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"ignore_initial_node_count", "extend_params", "enable_force_new", "partition",
				},
			},
		},
	})
}

func testAccNodePool_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  cluster_version        = "v1.34"
  enterprise_project_id  = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  enable_distribute_management = true
  container_network_type       = "eni"
  // eni_subnet_id is required if you want to enable the distribute cloud.
  eni_subnet_id                = join(",", [
    huaweicloud_vpc_subnet.test.ipv4_subnet_id, // Center subnet should be used.
  ])

  dynamic "masters" {
    for_each = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

    content {
      availability_zone = masters.value
    }
  }

  lifecycle {
    ignore_changes = [
      // After cluster created, the IEC subnet ID will append to the eni container subnet list.
      eni_subnet_id,
    ]
  }
}

data "huaweicloud_cce_flavor_specifications" "test" {
  cluster_type = "VirtualMachine"
}

locals {
  flavor_specifications = try(coalesce([
      for v in data.huaweicloud_cce_flavor_specifications.test.cluster_flavor_specs : try([
        for vv in v.available_master_flavors: vv.name if contains(vv.azs, data.huaweicloud_availability_zones.test.names[0])
      ][0], "") if v.name == huaweicloud_cce_cluster.test.flavor_id && !v.is_sold_out
    ])[0], null)
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}
`, common.TestVpc(name), name)
}

func testAccNodePool_basic_step1(name, baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  initial_node_count       = 1
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

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
      labels, extend_params
    ]
  }
}
`, baseConfig, name)
}

func testAccNodePool_basic_step2(name, baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id                = huaweicloud_cce_cluster.test.id
  name                      = "%[2]s"
  os                        = "EulerOS 2.9"
  flavor_id                 = local.flavor_specifications
  availability_zone         = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                  = huaweicloud_kps_keypair.test.name
  scall_enable              = true
  min_node_count            = 2
  max_node_count            = 9
  initial_node_count        = 2
  ignore_initial_node_count = false
  scale_down_cooldown_time  = 100
  priority                  = 1
  type                      = "vm"
  partition                 = "center"
  enterprise_project_id     = var.enterprise_project_id != "" ? var.enterprise_project_id : null

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
      labels, extend_params
    ]
  }
}
`, baseConfig, name)
}

func testAccNodePoolImportStateIdFunc(resName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterId, nodePoolId string
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of CCE node pool is not found in the tfstate", resName)
		}
		clusterId = rs.Primary.Attributes["cluster_id"]
		nodePoolId = rs.Primary.ID
		if clusterId == "" || nodePoolId == "" {
			return "", fmt.Errorf("the CCE node pool ID is not exist or related CCE cluster ID is missing")
		}
		return fmt.Sprintf("%s/%s", clusterId, nodePoolId), nil
	}
}

func TestAccNodePool_tagsLabelsTaints(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_tagsLabelsTaints_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.test1", "val1"),
					resource.TestCheckResourceAttr(resourceName, "tags.test2", "val2"),
					resource.TestCheckResourceAttr(resourceName, "labels.test1", "val1"),
					resource.TestCheckResourceAttr(resourceName, "labels.test2", "val2"),
					resource.TestCheckResourceAttr(resourceName, "labels.distribution.io/category", "Default"),
					resource.TestCheckResourceAttr(resourceName, "labels.distribution.io/partition", "center"),
					resource.TestCheckResourceAttr(resourceName, "labels.distribution.io/publicbordergroup", "center"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.value", "test_value"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(resourceName, "tag_policy_on_existing_nodes", "ignore"),
					resource.TestCheckResourceAttr(resourceName, "label_policy_on_existing_nodes", "ignore"),
					resource.TestCheckResourceAttr(resourceName, "taint_policy_on_existing_nodes", "ignore"),
				),
			},
			{
				Config: testAccNodePool_tagsLabelsTaints_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.test1", "val1_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.test2_update", "val2_update"),
					resource.TestCheckResourceAttr(resourceName, "labels.test1", "val1_update"),
					resource.TestCheckResourceAttr(resourceName, "labels.test2_update", "val2_update"),
					resource.TestCheckResourceAttr(resourceName, "labels.distribution.io/category", "Default"),
					resource.TestCheckResourceAttr(resourceName, "labels.distribution.io/partition", "center"),
					resource.TestCheckResourceAttr(resourceName, "labels.distribution.io/publicbordergroup", "center"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.value", "test_value_update"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(resourceName, "taints.1.key", "new_test_key"),
					resource.TestCheckResourceAttr(resourceName, "taints.1.value", "new_test_value"),
					resource.TestCheckResourceAttr(resourceName, "taints.1.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(resourceName, "tag_policy_on_existing_nodes", "refresh"),
					resource.TestCheckResourceAttr(resourceName, "label_policy_on_existing_nodes", "refresh"),
					resource.TestCheckResourceAttr(resourceName, "taint_policy_on_existing_nodes", "refresh"),
				),
			},
		},
	})
}

func testAccNodePool_tagsLabelsTaints_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 2
  initial_node_count       = 1
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  tag_policy_on_existing_nodes   = "ignore"
  label_policy_on_existing_nodes = "ignore"
  taint_policy_on_existing_nodes = "ignore"

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
    "distribution.io/category"          = "Default"
    "distribution.io/partition"         = "center"
    "distribution.io/publicbordergroup" = "center"
    "test1"                             = "val1"
    "test2"                             = "val2"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }
}
`, testAccNodePool_base(name), name)
}

func testAccNodePool_tagsLabelsTaints_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 2
  initial_node_count       = 1
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  tag_policy_on_existing_nodes   = "refresh"
  label_policy_on_existing_nodes = "refresh"
  taint_policy_on_existing_nodes = "refresh"

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
    "distribution.io/category"          = "Default"
    "distribution.io/partition"         = "center"
    "distribution.io/publicbordergroup" = "center"
    "test1"                             = "val1_update"
    "test2_update"                      = "val2_update"
  }

  taints {
    key    = "test_key"
    value  = "test_value_update"
    effect = "NoSchedule"
  }

  taints {
    key    = "new_test_key"
    value  = "new_test_value"
    effect = "NoSchedule"
  }
}
`, testAccNodePool_base(name), name)
}

func TestAccNodePool_volume_encryption(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_volume_encryption(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "root_volume.0.kms_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "data_volumes.0.kms_key_id"),
				),
			},
		},
	})
}

func testAccNodePool_volume_encryption(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%[2]s"
  pending_days = "7"
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 2
  initial_node_count       = 1
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

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

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, testAccNodePool_base(name), name)
}

func TestAccNodePool_prePaid(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_prePaid(name, "false"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "0"),
				),
			},
			{
				Config: testAccNodePool_prePaid(name, "true"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "0"),
				),
			},
		},
	})
}

func testAccNodePool_prePaid(rName, autoRenew string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  initial_node_count       = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[3]s"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, testAccNodePool_base(rName), rName, autoRenew)
}

func TestAccNodePool_SecurityGroups(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_SecurityGroups_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups.0",
						"huaweicloud_networking_secgroup.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups.1",
						"huaweicloud_networking_secgroup.test.1", "id"),
				),
			},
			{
				Config: testAccNodePool_SecurityGroups_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups.0",
						"huaweicloud_networking_secgroup.test.2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups.1",
						"huaweicloud_networking_secgroup.test.3", "id"),
				),
			},
		},
	})
}

func testAccNodePool_SecurityGroups_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test" {
  count                 = 4
  name                 = "%[2]s-secgroup-${count.index}"
  delete_default_rules = true
}
`, testAccNodePool_base(name), name)
}

func testAccNodePool_SecurityGroups_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  initial_node_count       = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

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

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, testAccNodePool_SecurityGroups_base(name), name)
}

func testAccNodePool_SecurityGroups_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  password                 = "test_123456"
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  initial_node_count       = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  security_groups = [
    huaweicloud_networking_secgroup.test[2].id,
    huaweicloud_networking_secgroup.test[3].id
  ]

  pod_security_groups = [
    huaweicloud_networking_secgroup.test[0].id,
    huaweicloud_networking_secgroup.test[1].id
  ]

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, testAccNodePool_SecurityGroups_base(name), name)
}

func TestAccNodePool_serverGroup(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_serverGroup(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "ecs_group_id",
						"huaweicloud_compute_servergroup.test", "id"),
				),
			},
		},
	})
}

func testAccNodePool_serverGroup(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_servergroup" "test" {
  name     = "%[2]s"
  policies = ["anti-affinity"]
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  initial_node_count       = 1
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  ecs_group_id             = huaweicloud_compute_servergroup.test.id
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, testAccNodePool_base(rName), rName)
}

func TestAccNodePool_storage(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_storage_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "storage.0.selectors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "storage.0.groups.#"),
				),
			},
			{
				Config: testAccNodePool_storage_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "storage.0.selectors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "storage.0.groups.#"),
				),
			},
		},
	})
}

func testAccNodePool_storage_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%[2]s"
  pending_days = "7"
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 2
  initial_node_count       = 1
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

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

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, testAccNodePool_base(rName), rName)
}

func testAccNodePool_storage_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%[2]s"
  pending_days = "7"
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 2
  initial_node_count       = 1
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  root_volume {
    size       = 50
    volumetype = "SSD"
  }
  
  data_volumes {
    size       = 200
    volumetype = "SSD"
    kms_key_id = huaweicloud_kms_key.test.id
  }

  data_volumes {
    size       = 200
    volumetype = "SSD"
    kms_key_id = huaweicloud_kms_key.test.id
  }

  storage {
    selectors {
      name              = "cceUse"
      type              = "evs"
      match_label_size  = "200"
      match_label_count = "1"
    }

    selectors {
      name                           = "user"
      type                           = "evs"
      match_label_size               = "200"
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
        size        = "20%%"
        lvm_lv_type = "linear"
      }

      virtual_spaces {
        name        = "runtime"
        size        = "80%%"
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

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, testAccNodePool_base(rName), rName)
}

func TestAccNodePool_extensionScaleGroups(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		updateName   = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_pool.test"

		baseConfig = testAccNodePool_base(name)

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_extensionScaleGroups_step1(name, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "scall_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", "10"),
					resource.TestCheckResourceAttr(resourceName, "scale_down_cooldown_time", "100"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "partition", "center"),
					resource.TestCheckResourceAttr(resourceName, "extension_scale_groups.#", "0"),
				),
			},
			{
				Config: testAccNodePool_extensionScaleGroups_step2(updateName, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "extension_scale_groups.#", "2"),
				),
			},
			{
				Config: testAccNodePool_extensionScaleGroups_step3(updateName, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "extension_scale_groups.#", "2"),
				),
			},
			{
				Config: testAccNodePool_extensionScaleGroups_step4(updateName, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "extension_scale_groups.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNodePoolImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"initial_node_count", "ignore_initial_node_count", "partition",
				},
			},
		},
	})
}

func testAccNodePool_extensionScaleGroups_step1(name, baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = true
  min_node_count           = 1
  max_node_count           = 10
  initial_node_count       = 2
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, baseConfig, name)
}

func testAccNodePool_extensionScaleGroups_step2(name, baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = true
  min_node_count           = 1
  max_node_count           = 10
  initial_node_count       = 2
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  extension_scale_groups {
    metadata {
      name = "group1"
    }

    spec {
      flavor = local.flavor_specifications
      az     = try(data.huaweicloud_availability_zones.test.names[0], null)

      autoscaling {
        extension_priority = 1
        enable             = true
        min_node_count     = 1
        max_node_count     = 3
      }
    }
  }

  extension_scale_groups {
    metadata {
      name = "group2"
    }

    spec {
      flavor = local.flavor_specifications
      az     = try(data.huaweicloud_availability_zones.test.names[0], null)

      autoscaling {
        extension_priority = 2
        enable             = true
        min_node_count     = 1
        max_node_count     = 3
      }
    }
  }

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, baseConfig, name)
}

func testAccNodePool_extensionScaleGroups_step3(name, baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = true
  min_node_count           = 1
  max_node_count           = 10
  initial_node_count       = 2
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  extension_scale_groups {
    metadata {
      name = "group1"
    }

    spec {
      flavor = local.flavor_specifications
      az     = try(data.huaweicloud_availability_zones.test.names[0], null)

      autoscaling {
        extension_priority = 1
        enable             = true
        min_node_count     = 1
        max_node_count     = 5
      }
    }
  }

  extension_scale_groups {
    metadata {
      name = "group2"
    }

    spec {
      flavor = local.flavor_specifications
      az     = try(data.huaweicloud_availability_zones.test.names[0], null)

      autoscaling {
        extension_priority = 3
        enable             = true
        min_node_count     = 1
        max_node_count     = 4
      }
    }
  }

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, baseConfig, name)
}

func testAccNodePool_extensionScaleGroups_step4(name, baseConfig string) string {
	return testAccNodePool_extensionScaleGroups_step1(name, baseConfig)
}

func TestAccNodePool_without_data_volumes(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_pool.test"

		baseConfig = testAccNodePool_base(name)

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_without_data_volumes(name, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func testAccNodePool_without_data_volumes(name, baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = local.flavor_specifications
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  initial_node_count       = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = "center"
  enterprise_project_id    = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  root_volume {
    size       = 40
    volumetype = "GPSSD"
  }

  storage {
    selectors {
      name = "cceUse"
      type = "system"
    }

    groups {
      name           = "vgpaas"
      selector_names = ["cceUse"]
      cce_managed    = true

      virtual_spaces {
        name        = "share"
        size        = "100%%"
        lvm_lv_type = "linear"
      }
    }
  }

  extension_scale_groups {
    metadata {
      name = "%[2]s-group1"
    }

    spec {
      flavor = local.flavor_specifications
      az     = try(data.huaweicloud_availability_zones.test.names[0], null)

      autoscaling {
        extension_priority = 1
        enable             = true
      }
    }
  }

  lifecycle {
    ignore_changes = [
      # When creating a node pool, the tags "distribution.io/category", "distribution.io/partition", and
      # "distribution.io/publicbordergroup" are automatically added.
      labels
    ]
  }
}
`, baseConfig, name)
}
