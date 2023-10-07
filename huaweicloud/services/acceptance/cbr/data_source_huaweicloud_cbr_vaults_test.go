package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
)

func TestAccDataVaults_backupServer(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_cbr_vaults.test"
		config         = testAccVault_backupServer_step1(testAccVault_base(name), name)
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVaults(config),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "200"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.tags.key", "value"),
				),
			},
		},
	})
}

func testAccDataVaults(config string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cbr_vaults" "test" {
  name = huaweicloud_cbr_vault.test.name
}
`, config)
}

func TestAccDataVaults_replicationServer(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_cbr_vaults.test"
		config         = testAccVault_replicationServer_step1(name)
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVaults(config),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "replication"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "200"),
				),
			},
		},
	})
}

func TestAccDataVaults_volume(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_cbr_vaults.test"
		config         = testAccVault_volume_step1(testAccVault_base(name), name)
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVaults(config),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "50"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id", "0"),
				),
			},
		},
	})
}

func TestAccDataVaults_backupTurbo(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_cbr_vaults.test"
		config         = testAccVault_backupTurbo_step1(name)
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVaults(config),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "800"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id", "0"),
				),
			},
		},
	})
}

func TestAccDataVaults_replicationTurbo(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_cbr_vaults.test"
		config         = testAccVault_replicationTurbo_step1(name)
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVaults(config),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "replication"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "1000"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id", "0"),
				),
			},
		},
	})
}
