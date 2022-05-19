package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
)

func TestAccVaults_BasicServer(t *testing.T) {
	randName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVaults_basic(testAccVault_serverBasic(randName)),
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
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.tags.key", "value"),
				),
			},
		},
	})
}

func TestAccVaults_ReplicaServer(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVaults_basic(testAccVault_serverReplication(randName)),
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

func TestAccVaults_BasicVolume(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVaults_basic(testAccVault_volumeBasic(randName)),
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
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccVaults_BasicTurbo(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVaults_basic(testAccVault_turboBasic(randName)),
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
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccVaults_ReplicaTurbo(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVaults_basic(testAccVault_turboReplication(randName)),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "replication"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "1000"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccVaults_basic(config string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cbr_vaults" "test" {
  name = huaweicloud_cbr_vault.test.name
}
`, config)
}
