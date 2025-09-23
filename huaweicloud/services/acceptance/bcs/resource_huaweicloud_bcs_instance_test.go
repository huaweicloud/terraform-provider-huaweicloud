package bcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bcs/v2/blockchains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBCSInstance_basic(t *testing.T) {
	var instance blockchains.BCSInstance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(3, "!@$%^-_=+[{}]:,./?"),
		acctest.RandIntRange(1, 3))
	resourceName := "huaweicloud_bcs_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBCSInstanceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBCSInstance_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBCSInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "4"),
					resource.TestCheckResourceAttr(resourceName, "consensus", "etcdraft"),
					resource.TestCheckResourceAttr(resourceName, "fabric_version", "2.2"),
					resource.TestCheckResourceAttr(resourceName, "blockchain_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "efs"),
					resource.TestCheckResourceAttr(resourceName, "org_disk_size", "1228"),
					resource.TestCheckResourceAttr(resourceName, "security_mechanism", "ECDSA"),
					resource.TestCheckResourceAttr(resourceName, "database_type", "goleveldb"),
					resource.TestCheckResourceAttr(resourceName, "orderer_node_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "peer_orgs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "peer_orgs.0.org_name", "organization01"),
					resource.TestCheckResourceAttr(resourceName, "peer_orgs.0.count", "1"),
					resource.TestCheckResourceAttr(resourceName, "channels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "channels.0.name", "channeldemo001"),
					resource.TestCheckResourceAttr(resourceName, "channels.0.org_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "channels.0.org_names.0", "organization01"),
					resource.TestCheckResourceAttr(resourceName, "sfs_turbo.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sfs_turbo.0.flavor", "sfs.turbo.40MBps"),
					resource.TestCheckResourceAttr(resourceName, "sfs_turbo.0.type", "efs-ha"),
					resource.TestCheckResourceAttr(resourceName, "sfs_turbo.0.share_type", "STANDARD"),
					resource.TestCheckResourceAttrPair(resourceName, "sfs_turbo.0.availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_storage",
					"eip_enable",
					"enterprise_project_id",
					"fabric_version",
					"orderer_node_num",
					"org_disk_size",
					"password",
					"volume_type",
					"sfs_turbo",
				},
			},
		},
	})
}

func TestAccBCSInstance_kafka(t *testing.T) {
	var instance blockchains.BCSInstance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(3, "!@$%^-_=+[{}]:,./?"),
		acctest.RandIntRange(1, 3))
	resourceName := "huaweicloud_bcs_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBCSInstanceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBCSInstance_kafka(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBCSInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "4"),
					resource.TestCheckResourceAttr(resourceName, "consensus", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "fabric_version", "1.4"),
					resource.TestCheckResourceAttr(resourceName, "blockchain_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "efs"),
					resource.TestCheckResourceAttr(resourceName, "org_disk_size", "500"),
					resource.TestCheckResourceAttr(resourceName, "security_mechanism", "ECDSA"),
					resource.TestCheckResourceAttr(resourceName, "database_type", "couchdb"),
					resource.TestCheckResourceAttr(resourceName, "orderer_node_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "couchdb.0.user_name", "Administrator"),
					resource.TestCheckResourceAttr(resourceName, "channels.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "peer_orgs.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sfs_turbo.0.share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "sfs_turbo.0.type", "efs-ha"),
					resource.TestCheckResourceAttr(resourceName, "sfs_turbo.0.flavor", "sfs.turbo.standard"),
					resource.TestCheckResourceAttrSet(resourceName, "sfs_turbo.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "kafka.0.flavor", "c3.mini"),
					resource.TestCheckResourceAttr(resourceName, "kafka.0.storage_size", "600"),
					resource.TestCheckResourceAttr(resourceName, "kafka.0.availability_zone.#", "1"),
				),
			},
		},
	})
}

func testAccCheckBCSInstanceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := cfg.BcsV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating BCS client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "huaweicloud_bcs_instance" {
				continue
			}

			id := rs.Primary.ID
			instance, err := blockchains.Get(client, id).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault400); ok {
					return nil
				}
				if _, ok := err.(golangsdk.ErrDefault401); ok {
					return nil
				}
				return err
			}
			if instance.Basic.ID != "" {
				return fmt.Errorf("%s (%s) still exists", rs.Type, id)
			}
		}
		return nil
	}
}

func testAccCheckBCSInstanceExists(name string, instance *blockchains.BCSInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := cfg.BcsV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating BCS client: %s", err)
		}

		found, err := blockchains.Get(client, id).Extract()
		if err != nil {
			return fmt.Errorf("error checking %s exist, err=%s", name, err)
		}
		if found.Basic.ID == "" {
			return fmt.Errorf("resource %s does not exist", name)
		}

		instance = found
		return nil
	}
}

func TestAccBCSInstance_updateWithEpsId(t *testing.T) {
	var instance blockchains.BCSInstance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(3, "!@$%^-_=+[{}]:,./?"),
		acctest.RandIntRange(1, 3))
	resourceName := "huaweicloud_bcs_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBCSInstanceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBCSInstance_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBCSInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testBCSInstance_basicTestWithEps(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBCSInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testBCSInstance_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "test" {
  name          = "%s"
  cidr          = "192.168.0.0/24"
  gateway_ip    = "192.168.0.1"
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.129.250"
  vpc_id        = huaweicloud_vpc.test.id
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s2.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  cluster_version        = "v1.19.16-r1"
  delete_sfs             = true
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%s"

  lifecycle {
    ignore_changes = [
      public_key,
    ]
  }
}

resource "huaweicloud_cce_node" "test" {
  name                  = "%s"
  cluster_id            = huaweicloud_cce_cluster.test.id
  flavor_id             = "s6.xlarge.2"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  key_pair              = huaweicloud_kps_keypair.test.name
  max_pods              = 30
  ecs_performance_type  = "normal"
  os                    = "CentOS 7.6"

  root_volume {
    size       = 40
    volumetype = "SAS"
  }
  data_volumes {
    size = 100
    volumetype = "SAS"
  }
}
`, rName, rName, rName, rName, rName)
}

func testBCSInstance_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_bcs_instance" "test" {
  depends_on = [huaweicloud_cce_node.test]

  name                  = "%s"
  cce_cluster_id        = huaweicloud_cce_cluster.test.id
  consensus             = "etcdraft"
  edition               = 4
  enterprise_project_id = "%s"
  fabric_version        = "2.2"
  password              = "%s"
  volume_type           = "efs"
  org_disk_size         = 1228
  security_mechanism    = "ECDSA"
  orderer_node_num      = 3
  delete_storage        = true

  peer_orgs {
    org_name = "organization01"
    count    = 1
  }
  channels {
    name      = "channeldemo001"
    org_names = [
      "organization01",
    ]
  }
  sfs_turbo {
    flavor            = "sfs.turbo.40MBps"
    type              = "efs-ha"
    share_type        = "STANDARD"
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
  }
}
`, testBCSInstance_base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, password)
}

func testBCSInstance_basicTestWithEps(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_bcs_instance" "test" {
  depends_on = [huaweicloud_cce_node.test]

  name                  = "%s"
  cce_cluster_id        = huaweicloud_cce_cluster.test.id
  consensus             = "etcdraft"
  edition               = 4
  enterprise_project_id = "%s"
  fabric_version        = "2.2"
  password              = "%s"
  volume_type           = "efs"
  org_disk_size         = 1228
  security_mechanism    = "ECDSA"
  orderer_node_num      = 3
  delete_storage        = true

  peer_orgs {
    org_name = "organization01"
    count    = 1
  }
  channels {
    name      = "channeldemo001"
    org_names = [
      "organization01",
    ]
  }
  sfs_turbo {
    flavor            = "sfs.turbo.40MBps"
    type              = "efs-ha"
    share_type        = "STANDARD"
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
  }
}
`, testBCSInstance_base(rName), rName, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST, password)
}

func testBCSInstance_kafka(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_bcs_instance" "test" {
  depends_on = [huaweicloud_cce_node.test]

  name                  = "%s"
  cce_cluster_id        = huaweicloud_cce_cluster.test.id
  consensus             = "kafka"
  edition               = 4
  enterprise_project_id = "%s"
  fabric_version        = "1.4"
  password              = "%s"
  volume_type           = "efs"
  org_disk_size         = 500
  database_type         = "couchdb"
  orderer_node_num      = 2
  delete_storage        = true
  delete_obs            = true

  couchdb {
    user_name = "Administrator"
    password  = "%s"
  }
  peer_orgs {
    org_name = "organization01"
    count    = 2
  }
  peer_orgs {
    org_name = "organization02"
    count    = 2
  }
  channels {
    name      = "channeldemo001"
    org_names = [
      "organization01",
      "organization02",
    ]
  }
  channels {
    name      = "channeldemo002"
    org_names = [
      "organization02",
    ]
  }
  sfs_turbo {
    share_type        = "STANDARD"
    type              = "efs-ha"
    flavor            = "sfs.turbo.standard"
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
  }
  kafka {
    flavor            = "c3.mini"
    storage_size      = 600
    availability_zone = [
  	data.huaweicloud_availability_zones.test.names[0],
    ]
  }
}
`, testBCSInstance_base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, password, password)
}
