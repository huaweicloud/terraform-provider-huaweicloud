package huaweicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cbr/v3/vaults"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCBRV3Vault_serverBasic(t *testing.T) {
	var vault vaults.Vault
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRV3VaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_serverBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3VaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", "server"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccCBRV3Vault_serverUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3VaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", "server"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "300"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
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

func TestAccCBRV3Vault_serverReplication(t *testing.T) {
	var vault vaults.Vault
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRV3VaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_serverReplication(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3VaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", "server"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
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

func TestAccCBRV3Vault_volumeBasic(t *testing.T) {
	var vault vaults.Vault
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRV3VaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_volumeBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3VaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", "disk"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "50"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccCBRV3Vault_volumeUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3VaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", "disk"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
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

func TestAccCBRV3Vault_turboBasic(t *testing.T) {
	var vault vaults.Vault
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRV3VaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_turboBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3VaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", "turbo"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccCBRV3Vault_turboUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3VaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", "turbo"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
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

func TestAccCBRV3Vault_turboReplication(t *testing.T) {
	var vault vaults.Vault
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRV3VaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_turboReplication(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3VaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", "turbo"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
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

func testAccCheckCBRV3VaultDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.CbrV3Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Huaweicloud CBR client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cbr_vault" {
			continue
		}

		_, err := vaults.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Vault still exists")
		}
	}

	return nil
}

func testAccCheckCBRV3VaultExists(n string, vault *vaults.Vault) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CbrV3Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating Huaweicloud CBR client: %s", err)
		}

		found, err := vaults.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] test found is: %#v", found)
		vault = found

		return nil
	}
}

func testAccCBRV3Vault_policy(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name        = "%s"
  type        = "backup"
  time_period = 20

  backup_cycle {
    days            = "MO,TU"
    execution_times = ["06:00", "18:00"]
  }
}
`, rName)
}

//Vaults of type 'server'
func testAccCBRV3Vault_serverBase(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/20"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_compute_keypair" "test" {
  name = "%s"
  lifecycle {
    ignore_changes = [
      public_key,
    ]
  }
}

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type  = "SAS"
  system_disk_size  = 50
  key_pair          = huaweicloud_compute_keypair.test.name

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, rName, rName, rName, rName)
}

func testAccCBRV3Vault_serverBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCBRV3Vault_serverBase(rName), rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCBRV3Vault_serverUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s-update"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 300
  enterprise_project_id = "%s"
  policy_id             = huaweicloud_cbr_policy.test.id

  resources {
    id = huaweicloud_compute_instance.test.id
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, testAccCBRV3Vault_serverBase(rName), testAccCBRV3Vault_policy(rName), rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCBRV3Vault_serverReplication(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "replication"
  size                  = 200
  enterprise_project_id = "%s"
}
`, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

//Vaults of type 'disk'
func testAccCBRV3Vault_volumeBase(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s_1"
  volume_type       = "SAS"
  size              = 20
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`, rName)
}

func testAccCBRV3Vault_volumeBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  consistent_level      = "crash_consistent"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 50
  enterprise_project_id = "%s"
}
`, testAccCBRV3Vault_volumeBase(rName), rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCBRV3Vault_volumeUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s-update"
  consistent_level      = "crash_consistent"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 100
  auto_expand           = true
  enterprise_project_id = "%s"
  policy_id             = huaweicloud_cbr_policy.test.id

  resources {
    id = huaweicloud_evs_volume.test.id
  }
}
`, testAccCBRV3Vault_volumeBase(rName), testAccCBRV3Vault_policy(rName), rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

//Vaults of type 'turbo'
func testAccCBRV3Vault_turboBase(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/20"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/22"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_sfs_turbo" "test" {
  name              = "%s"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`, rName, rName, rName, rName)
}

func testAccCBRV3Vault_turboBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  consistent_level      = "crash_consistent"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 800
  enterprise_project_id = "%s"
}
`, testAccCBRV3Vault_turboBase(rName), rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCBRV3Vault_turboUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s-update"
  consistent_level      = "crash_consistent"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 1000
  enterprise_project_id = "%s"
  policy_id             = huaweicloud_cbr_policy.test.id

  resources {
    id = huaweicloud_sfs_turbo.test.id
  }
}
`, testAccCBRV3Vault_turboBase(rName), testAccCBRV3Vault_policy(rName), rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCBRV3Vault_turboReplication(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  consistent_level      = "crash_consistent"
  type                  = "turbo"
  protection_type       = "replication"
  size                  = 1000
  enterprise_project_id = "%s"
}
`, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}
