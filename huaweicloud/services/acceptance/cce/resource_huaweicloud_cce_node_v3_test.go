package cce

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCCENodeV3_basic(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "huaweicloud_cce_node.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENodeImportStateIdFunc(),
			},
			{
				Config: testAccCCENodeV3_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccCCENodeV3_auto_assign_eip(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_auto_assign_eip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestMatchResourceAttr(resourceName, "public_ip",
						regexp.MustCompile("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$")),
				),
			},
		},
	})
}

func TestAccCCENodeV3_existing_eip(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_existing_eip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestMatchResourceAttr(resourceName, "public_ip",
						regexp.MustCompile("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$")),
				),
			},
		},
	})
}

func TestAccCCENodeV3_volume_extendParams(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_volume_extendParams(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.extend_params.test_key", "test_val"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.0.extend_params.test_key", "test_val"),
				),
			},
		},
	})
}

func TestAccCCENodeV3_data_volume_encryption(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_data_volume_encryption(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "data_volumes.0.kms_key_id"),
				),
			},
		},
	})
}

func TestAccCCENodeV3_prePaid(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_prePaid(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccCCENodeV3_prePaid(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccCCENodeV3_password(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_password(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccCCENodeV3_storage(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_storage(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckCCENodeV3Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	cceClient, err := config.CceV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	var clusterId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "huaweicloud_cce_cluster" {
			clusterId = rs.Primary.ID
		}

		if rs.Type != "huaweicloud_cce_node" {
			continue
		}

		_, err := nodes.Get(cceClient, clusterId, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Node still exists")
		}
	}

	return nil
}

func testAccCheckCCENodeV3Exists(n string, cluster string, node *nodes.Nodes) resource.TestCheckFunc {
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

		found, err := nodes.Get(cceClient, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmtp.Errorf("Node not found")
		}

		*node = *found

		return nil
	}
}

func testAccCCENodeImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		cluster, ok := s.RootModule().Resources["huaweicloud_cce_cluster.test"]
		if !ok {
			return "", fmtp.Errorf("Cluster not found: %s", cluster)
		}
		node, ok := s.RootModule().Resources["huaweicloud_cce_node.test"]
		if !ok {
			return "", fmtp.Errorf("Node not found: %s", node)
		}
		if cluster.Primary.ID == "" || node.Primary.ID == "" {
			return "", fmtp.Errorf("resource not found: %s/%s", cluster.Primary.ID, node.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", cluster.Primary.ID, node.Primary.ID), nil
	}
}

func testAccCCENodeV3_Base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_compute_keypair" "test" {
  name = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, testAccCCEClusterV3_Base(rName), rName, rName)
}

func testAccCCENodeV3_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_compute_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCENodeV3_Base(rName), rName)
}

func testAccCCENodeV3_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_compute_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
  tags = {
    foo = "bar"
    key = "value_update"
  }
}
`, testAccCCENodeV3_Base(rName), updateName)
}

func testAccCCENodeV3_auto_assign_eip(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_compute_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  // Assign EIP
  iptype="5_bgp"
  bandwidth_charge_mode="traffic"
  sharetype= "PER"
  bandwidth_size= 100
}
`, testAccCCENodeV3_Base(rName), rName)
}

func testAccCCENodeV3_existing_eip(rName string) string {
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

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_compute_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  // Assign existing EIP
  eip_id = huaweicloud_vpc_eip.test.id
}
`, testAccCCENodeV3_Base(rName), rName)
}

func testAccCCENodeV3_volume_extendParams(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_compute_keypair.test.name

  root_volume {
    size         = 40
	volumetype   = "SSD"
	extend_params = {
	  test_key = "test_val"
	}
  }
  data_volumes {
    size         = 100
	volumetype   = "SSD"
	extend_params = {
	  test_key = "test_val"
	}
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCENodeV3_Base(rName), rName)
}

func testAccCCENodeV3_data_volume_encryption(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_compute_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = huaweicloud_kms_key.test.id
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCENodeV3_Base(rName), rName, rName)
}

func testAccCCENodeV3_prePaid(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%[2]s"
  flavor_id         = "s6.large.2"
  availability_zone = try(element(data.huaweicloud_availability_zones.test.names, 0), null)
  key_pair          = huaweicloud_compute_keypair.test.name

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "%[3]v"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccCCENodeV3_Base(rName), rName, isAutoRenew)
}

func testAccCCENodeV3_password(rName string) string {
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

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  password          = "Test@123"
  eip_id            = huaweicloud_vpc_eip.test.id

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  provisioner "remote-exec" {
    connection {
      user     = "root"
      password = "Test@123"
      host     = huaweicloud_vpc_eip.test.address
      port     = "22"
      timeout  = "60s"
    }

    inline = [
      "date"
    ]
  }
}
`, testAccCCENodeV3_Base(rName), rName)
}

func testAccCCENodeV3_storage(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_compute_keypair.test.name

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
      match_label_count = 1
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
`, testAccCCENodeV3_Base(rName), rName, rName)
}
