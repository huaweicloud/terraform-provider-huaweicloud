package cbr

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
%[1]s

// base compute resources
%[2]s

resource "huaweicloud_compute_instance" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[3]s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]

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
  name              = "%[3]s_${tostring(count.index)}"
  size              = var.volume_configuration[count.index].size
  device_type       = var.volume_configuration[count.index].device_type
}

resource "huaweicloud_compute_volume_attach" "test" {
  count = length(var.volume_configuration)

  instance_id = huaweicloud_compute_instance.test.id
  volume_id   = huaweicloud_evs_volume.test[count.index].id
}`, config, common.TestBaseComputeResources(rName), rName)
}

func testAccVault_serverBasic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "%[3]s"

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
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s-update"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 300
  enterprise_project_id = "%[3]s"

  resources {
    server_id = huaweicloud_compute_instance.test.id

    excludes = huaweicloud_compute_volume_attach.test[*].volume_id
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, testAccVaultBasicConfiguration(testAccEvsVolumeConfiguration_update(), rName),
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
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s-update"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "%[3]s"

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
`, testAccVaultBasicConfiguration(testAccEvsVolumeConfiguration_update(), rName),
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
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s-update"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 100
  auto_expand           = true
  enterprise_project_id = "%[3]s"

  resources {
    includes = huaweicloud_compute_volume_attach.test[*].volume_id
  }
}
`, testAccVaultBasicConfiguration(testAccEvsVolumeConfiguration_basic(), rName),
		rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// Vaults of type 'turbo'
func testAccVault_turboBase(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test1" {
  name              = "%s-1"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}`, common.TestBaseNetwork(rName), rName)
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
%[1]s

resource "huaweicloud_sfs_turbo" "test2" {
  name              = "%[2]s-2"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s-update"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 1000
  enterprise_project_id = "%[3]s"

  resources {
    includes = [
      huaweicloud_sfs_turbo.test1.id,
      huaweicloud_sfs_turbo.test2.id
    ]
  }
}
`, testAccVault_turboBase(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
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

func TestAccVault_bindPolicies(t *testing.T) {
	var (
		vault        vaults.Vault
		randName     = acceptance.RandomAccResourceName()
		mainRcName   = "huaweicloud_cbr_vault.test"
		legacyRcName = "huaweicloud_cbr_vault.legacy"
	)

	mainRc := acceptance.InitResourceCheck(
		mainRcName,
		&vault,
		getVaultResourceFunc,
	)

	legacyRc := acceptance.InitResourceCheck(
		legacyRcName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckReplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      mainRc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_bindPolicies_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					mainRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(mainRcName, "name", randName),
					resource.TestCheckResourceAttr(mainRcName, "policy.#", "2"),
					legacyRc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(legacyRcName, "policy_id",
						"huaweicloud_cbr_policy.backup.0", "id"),
				),
			},
			{
				Config: testAccVault_bindPolicies_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					mainRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(mainRcName, "name", randName),
					resource.TestCheckResourceAttr(mainRcName, "policy.#", "2"),
					legacyRc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(legacyRcName, "policy_id",
						"huaweicloud_cbr_policy.backup.1", "id"),
				),
			},
			{
				ResourceName:      mainRcName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      legacyRcName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"policy_id",
				},
			},
		},
	})
}

func testAccVault_bindPolicies_base(name string) string {
	return fmt.Sprintf(`
variable "backup_configuration" {
  type = list(object({
    interval        = number
    days            = string
    execution_times = list(string)
  }))
  default = [{
    interval        = 5
    days            = null
    execution_times = ["06:00", "18:00"]
  },
  {
    interval        = null
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }]
}

resource "huaweicloud_cbr_policy" "backup" {
  count = 2

  name            = "%[1]s"
  type            = "backup"
  backup_quantity = 5

  backup_cycle {
    interval        = var.backup_configuration[count.index].interval
    days            = var.backup_configuration[count.index].days
    execution_times = var.backup_configuration[count.index].execution_times
  }
}

resource "huaweicloud_cbr_policy" "replication" {
  count = 2

  name                   = "%[1]s"
  type                   = "replication"
  destination_region     = "%[2]s"
  destination_project_id = "%[3]s"
  time_period            = 20

  backup_cycle {
    interval        = var.backup_configuration[count.index].interval
    days            = var.backup_configuration[count.index].days
    execution_times = var.backup_configuration[count.index].execution_times
  }
}

resource "huaweicloud_cbr_vault" "destination" {
  region          = "%[2]s"
  name            = "%[1]s"
  type            = "server"
  protection_type = "replication"
  size            = 200
}
`, name, acceptance.HW_DEST_REGION, acceptance.HW_DEST_PROJECT_ID)
}

func testAccVault_bindPolicies_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "legacy" {
  name            = "%[2]s_legacy"
  type            = "server"
  protection_type = "backup"
  size            = 200
  policy_id       = huaweicloud_cbr_policy.backup[0].id
}

resource "huaweicloud_cbr_vault" "test" {
  name            = "%[2]s"
  type            = "server"
  protection_type = "backup"
  size            = 200

  policy {
    id = huaweicloud_cbr_policy.backup[0].id
  }
  policy {
    id                   = huaweicloud_cbr_policy.replication[0].id
    destination_vault_id = huaweicloud_cbr_vault.destination.id
  }
}
`, testAccVault_bindPolicies_base(name), name)
}

func testAccVault_bindPolicies_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "legacy" {
  name            = "%[2]s_legacy"
  type            = "server"
  protection_type = "backup"
  size            = 200
  policy_id       = huaweicloud_cbr_policy.backup[1].id
}

resource "huaweicloud_cbr_vault" "test" {
  name            = "%[2]s"
  type            = "server"
  protection_type = "backup"
  size            = 200

  policy {
    id = huaweicloud_cbr_policy.backup[1].id
  }
  policy {
    id                   = huaweicloud_cbr_policy.replication[1].id
    destination_vault_id = huaweicloud_cbr_vault.destination.id
  }
}
`, testAccVault_bindPolicies_base(name), name)
}
