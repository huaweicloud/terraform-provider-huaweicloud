package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
)

func TestAccCbrVaultsV3_BasicServer(t *testing.T) {
	randName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCbrVaultsV3_serverBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "200"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.tags.key", "value"),
				),
			},
		},
	})
}

func TestAccCbrVaultsV3_ReplicaServer(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCbrVaultsV3_serverReplication(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
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

func TestAccCbrVaultsV3_BasicVolume(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCbrVaultsV3_volumeBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "50"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func TestAccCbrVaultsV3_BasicTurbo(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCbrVaultsV3_turboBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "800"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func TestAccCbrVaultsV3_ReplicaTurbo(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCbrVaultsV3_turboReplication(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "replication"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "1000"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func testAccCbrVaultsV3_serverBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cbr_vaults" "test" {
  name = huaweicloud_cbr_vault.test.name
}
`, testAccCBRV3Vault_serverBasic(rName))
}

func testAccCbrVaultsV3_serverReplication(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cbr_vaults" "test" {
  name = huaweicloud_cbr_vault.test.name
}
`, testAccCBRV3Vault_serverReplication(rName))
}

func testAccCbrVaultsV3_volumeBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cbr_vaults" "test" {
  name = huaweicloud_cbr_vault.test.name
}
`, testAccCBRV3Vault_volumeBasic(rName))
}

func testAccCbrVaultsV3_turboBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cbr_vaults" "test" {
  name = huaweicloud_cbr_vault.test.name
}
`, testAccCBRV3Vault_turboBasic(rName))
}

func testAccCbrVaultsV3_turboReplication(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cbr_vaults" "test" {
  name = huaweicloud_cbr_vault.test.name
}
`, testAccCBRV3Vault_turboReplication(rName))
}
