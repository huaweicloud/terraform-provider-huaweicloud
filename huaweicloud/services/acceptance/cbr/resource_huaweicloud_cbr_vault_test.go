package cbr

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
)

func getVaultResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CbrV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR v3 client: %s", err)
	}
	return vaults.Get(c, state.Primary.ID).Extract()
}

func TestAccVault_BasicServer(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_serverBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
				),
			},
			{
				Config: testAccVault_serverUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "300"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "policy_id",
						"${huaweicloud_cbr_policy.test.id}"),
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

func TestAccVault_ReplicaServer(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_serverReplication(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
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

func TestAccVault_prePaidServer(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_serverPrePaid(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccVault_serverPrePaidUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "policy_id",
						"${huaweicloud_cbr_policy.test.id}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"period_unit",
					"period",
					"auto_renew",
					"auto_pay",
				},
			},
		},
	})
}

func TestAccVault_BasicVolume(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_volumeBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "50"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
				),
			},
			{
				Config: testAccVault_volumeUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
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

func TestAccVault_BasicTurbo(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_turboBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccVault_turboUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
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

func TestAccVault_ReplicaTurbo(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_turboReplication(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
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

func TestAccVault_AutoBind(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_autoBind(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "auto_bind", "true"),
					resource.TestCheckResourceAttr(resourceName, "bind_rules.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "bind_rules.key", "value"),
				),
			},
			{
				Config: testAccVault_autoBindUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "auto_bind", "false"),
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

func testAccVault_policy(rName string) string {
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

func testAccEvsVolumeConfiguration_basic() string {
	return `
variable "volume_configuration" {
  type = list(object({
    volume_type = string
    size        = number
    device_type = string
  }))
  default = [
    {volume_type = "SSD", size = 100, device_type = "VBD"},
    {volume_type = "SSD", size = 100, device_type = "SCSI"},
  ]
}`
}

func testAccEvsVolumeConfiguration_update() string {
	return `
variable "volume_configuration" {
  type = list(object({
    volume_type = string
    size        = number
    device_type = string
  }))
  default = [
    {volume_type = "GPSSD", size = 100, device_type = "VBD"},
    {volume_type = "SAS", size = 100, device_type = "SCSI"},
  ]
}`
}

func testAccVaultBasicConfiguration(config, rName string) string {
	return fmt.Sprintf(`
%s

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

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%s"
  delete_default_rules = true
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
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  key_pair          = huaweicloud_compute_keypair.test.name

  system_disk_type = "SSD"
  system_disk_size = 50

  security_group_ids = [
    huaweicloud_networking_secgroup.test.id
  ]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = var.volume_configuration[count.index].volume_type
  name              = "%s_${tostring(count.index)}"
  size              = var.volume_configuration[count.index].size
  device_type       = var.volume_configuration[count.index].device_type
}

resource "huaweicloud_compute_volume_attach" "test" {
  count = length(var.volume_configuration)

  instance_id = huaweicloud_compute_instance.test.id
  volume_id   = huaweicloud_evs_volume.test[count.index].id
}`, config, rName, rName, rName, rName, rName, rName)
}

func testAccVault_serverBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "%s"

  resources {
    server_id = huaweicloud_compute_instance.test.id
    excludes  = huaweicloud_compute_volume_attach.test[*].volume_id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccVaultBasicConfiguration(testAccEvsVolumeConfiguration_basic(), rName),
		rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccVault_serverUpdate(rName string) string {
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
    server_id = huaweicloud_compute_instance.test.id

    excludes = huaweicloud_compute_volume_attach.test[*].volume_id
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, testAccVaultBasicConfiguration(testAccEvsVolumeConfiguration_update(), rName), testAccVault_policy(rName),
		rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccVault_serverReplication(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "replication"
  size                  = 200
  enterprise_project_id = "%s"
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccVault_serverPrePaid(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "%s"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1

  resources {
    server_id = huaweicloud_compute_instance.test.id

    excludes = huaweicloud_compute_volume_attach.test[*].volume_id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccVaultBasicConfiguration(testAccEvsVolumeConfiguration_update(), rName),
		rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccVault_serverPrePaidUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s-update"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "%s"
  policy_id             = huaweicloud_cbr_policy.test.id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  resources {
    server_id = huaweicloud_compute_instance.test.id

    excludes = huaweicloud_compute_volume_attach.test[*].volume_id
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, testAccVaultBasicConfiguration(testAccEvsVolumeConfiguration_update(), rName), testAccVault_policy(rName),
		rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccVault_volumeBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 50
  enterprise_project_id = "%s"

  resources {
    includes = huaweicloud_compute_volume_attach.test[*].volume_id
  }
}
`, testAccVaultBasicConfiguration(testAccEvsVolumeConfiguration_basic(), rName),
		rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccVault_volumeUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s-update"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 100
  auto_expand           = true
  enterprise_project_id = "%s"
  policy_id             = huaweicloud_cbr_policy.test.id

  resources {
    includes = huaweicloud_compute_volume_attach.test[*].volume_id
  }
}
`, testAccVaultBasicConfiguration(testAccEvsVolumeConfiguration_basic(), rName),
		testAccVault_policy(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// Vaults of type 'turbo'
func testAccVault_turboBase(rName string) string {
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

resource "huaweicloud_sfs_turbo" "test1" {
  name              = "%s-1"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}`, rName, rName, rName, rName)
}

func testAccVault_turboBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 800
  enterprise_project_id = "%s"

  resources {
    includes = [
      huaweicloud_sfs_turbo.test1.id
    ]
  }
}
`, testAccVault_turboBase(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccVault_turboUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_sfs_turbo" "test2" {
  name              = "%s-2"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s-update"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 1000
  enterprise_project_id = "%s"
  policy_id             = huaweicloud_cbr_policy.test.id

  resources {
    includes = [
      huaweicloud_sfs_turbo.test1.id,
      huaweicloud_sfs_turbo.test2.id
    ]
  }
}
`, testAccVault_turboBase(rName), testAccVault_policy(rName), rName, rName,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccVault_turboReplication(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name                  = "%s"
  type                  = "turbo"
  protection_type       = "replication"
  size                  = 1000
  enterprise_project_id = "%s"
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccVault_autoBind(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name            = "%s"
  type            = "server"
  protection_type = "backup"
  size            = 800
  auto_bind       = true
  bind_rules      = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccVault_autoBindUpdate(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name            = "%s"
  type            = "server"
  protection_type = "backup"
  size            = 800
  auto_bind       = false
}
`, rName)
}
