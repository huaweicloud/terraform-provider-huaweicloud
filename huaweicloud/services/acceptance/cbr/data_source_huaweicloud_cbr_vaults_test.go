package cbr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
)

func TestAccDataVaults_backupServer(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_cbr_vaults.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_cbr_vaults.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_cbr_vaults.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byConsistentLevel   = "data.huaweicloud_cbr_vaults.filter_by_consistent_level"
		dcByConsistentLevel = acceptance.InitDataSourceCheck(byConsistentLevel)

		byProtectionType   = "data.huaweicloud_cbr_vaults.filter_by_protection_type"
		dcByProtectionType = acceptance.InitDataSourceCheck(byProtectionType)

		bySize   = "data.huaweicloud_cbr_vaults.filter_by_size"
		dcBySize = acceptance.InitDataSourceCheck(bySize)

		byAutoExpand   = "data.huaweicloud_cbr_vaults.filter_by_auto_expand"
		dcByAutoExpand = acceptance.InitDataSourceCheck(byAutoExpand)

		byEpsId   = "data.huaweicloud_cbr_vaults.filter_by_eps_id"
		dcByEpsId = acceptance.InitDataSourceCheck(byEpsId)

		byPolicyId   = "data.huaweicloud_cbr_vaults.filter_by_policy_id"
		dcByPolicyId = acceptance.InitDataSourceCheck(byPolicyId)

		byStatus   = "data.huaweicloud_cbr_vaults.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVaults_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "vaults.#", regexp.MustCompile(`[1-9]\d*`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttr(byName, "name", name),
					resource.TestCheckResourceAttr(byName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(byName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(byName, "vaults.0.type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(byName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(byName, "vaults.0.size", "200"),
					resource.TestCheckResourceAttr(byName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(byName, "vaults.0.enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(byName, "vaults.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(byName, "vaults.0.tags.key", "value"),
					resource.TestCheckResourceAttr(byName, "vaults.0.auto_bind", "true"),
					resource.TestCheckResourceAttr(byName, "vaults.0.bind_rules.foo", "bar"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByConsistentLevel.CheckResourceExists(),
					resource.TestCheckOutput("is_consistent_level_filter_useful", "true"),
					dcByProtectionType.CheckResourceExists(),
					resource.TestCheckOutput("is_protection_type_filter_useful", "true"),
					dcBySize.CheckResourceExists(),
					resource.TestCheckOutput("is_size_filter_useful", "true"),
					dcByAutoExpand.CheckResourceExists(),
					resource.TestCheckOutput("is_auto_expand_filter_useful", "true"),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					dcByPolicyId.CheckResourceExists(),
					resource.TestCheckOutput("is_policy_id_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataVaults_basic_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_policy" "test" {
  name            = "%[2]s"
  type            = "backup"
  backup_quantity = 5

  backup_cycle {
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "0"
  backup_name_prefix    = "test-prefix-"
  is_multi_az           = true
  auto_bind             = true

  bind_rules = {
    foo = "bar"
  }

  resources {
    server_id = huaweicloud_compute_instance.test.id
    excludes  = slice(huaweicloud_compute_volume_attach.test[*].volume_id, 0, 2)
  }

  policy {
    id = huaweicloud_cbr_policy.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccVault_base(name), name)
}

func testAccDataVaults_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cbr_vaults" "test" {
  depends_on = [huaweicloud_cbr_vault.test]
}

# Filter by name
locals {
  vault_name = huaweicloud_cbr_vault.test.name
}

data "huaweicloud_cbr_vaults" "filter_by_name" {
  depends_on = [huaweicloud_cbr_vault.test]

  name = local.vault_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_cbr_vaults.filter_by_name.vaults[*].name : v == local.vault_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by type
locals {
  vault_type = huaweicloud_cbr_vault.test.type
}

data "huaweicloud_cbr_vaults" "filter_by_type" {
  depends_on = [huaweicloud_cbr_vault.test]

  type = local.vault_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_cbr_vaults.filter_by_type.vaults[*].type : v == local.vault_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by consistent_level
locals {
  consistent_level = huaweicloud_cbr_vault.test.consistent_level
}

data "huaweicloud_cbr_vaults" "filter_by_consistent_level" {
  depends_on = [huaweicloud_cbr_vault.test]

  consistent_level = local.consistent_level
}

locals {
  consistent_level_filter_result = [
    for v in data.huaweicloud_cbr_vaults.filter_by_consistent_level.vaults[*].consistent_level : v == local.consistent_level
  ]
}

output "is_consistent_level_filter_useful" {
  value = length(local.consistent_level_filter_result) > 0 && alltrue(local.consistent_level_filter_result)
}

# Filter by protection_type
locals {
  protection_type = huaweicloud_cbr_vault.test.protection_type
}

data "huaweicloud_cbr_vaults" "filter_by_protection_type" {
  depends_on = [huaweicloud_cbr_vault.test]

  protection_type = local.protection_type
}

locals {
  protection_type_filter_result = [
    for v in data.huaweicloud_cbr_vaults.filter_by_protection_type.vaults[*].protection_type : v == local.protection_type
  ]
}

output "is_protection_type_filter_useful" {
  value = length(local.protection_type_filter_result) > 0 && alltrue(local.protection_type_filter_result)
}

# Filter by size
locals {
  vault_size = huaweicloud_cbr_vault.test.size
}

data "huaweicloud_cbr_vaults" "filter_by_size" {
  depends_on = [huaweicloud_cbr_vault.test]

  size = local.vault_size
}

locals {
  size_filter_result = [
    for v in data.huaweicloud_cbr_vaults.filter_by_size.vaults[*].size : v == local.vault_size
  ]
}

output "is_size_filter_useful" {
  value = length(local.size_filter_result) > 0 && alltrue(local.size_filter_result)
}

# Filter by auto_expand
locals {
  auto_expand = huaweicloud_cbr_vault.test.auto_expand
}

data "huaweicloud_cbr_vaults" "filter_by_auto_expand" {
  depends_on = [huaweicloud_cbr_vault.test]

  auto_expand_enabled = local.auto_expand
}

locals {
  auto_expand_filter_result = [
    for v in data.huaweicloud_cbr_vaults.filter_by_auto_expand.vaults[*].auto_expand_enabled : v == local.auto_expand
  ]
}

output "is_auto_expand_filter_useful" {
  value = length(local.auto_expand_filter_result) > 0 && alltrue(local.auto_expand_filter_result)
}

# Filter by enterprise_project_id
locals {
  enterprise_project_id = huaweicloud_cbr_vault.test.enterprise_project_id
}

data "huaweicloud_cbr_vaults" "filter_by_eps_id" {
  depends_on = [huaweicloud_cbr_vault.test]

  enterprise_project_id = local.enterprise_project_id
}

locals {
  eps_id_filter_result = [
    for v in data.huaweicloud_cbr_vaults.filter_by_eps_id.vaults[*].enterprise_project_id : v == local.enterprise_project_id
  ]
}

output "is_eps_id_filter_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}

# Filter by policy_id
locals {
  policy_id = huaweicloud_cbr_policy.test.id
}

data "huaweicloud_cbr_vaults" "filter_by_policy_id" {
  depends_on = [huaweicloud_cbr_vault.test]

  policy_id = local.policy_id
}

locals {
  policy_id_filter_result = [
    for v in data.huaweicloud_cbr_vaults.filter_by_policy_id.vaults[*].policy_id : v == local.policy_id
  ]
}

output "is_policy_id_filter_useful" {
  value = length(local.policy_id_filter_result) > 0 && alltrue(local.policy_id_filter_result)
}

# Filter by status
locals {
  vault_status = huaweicloud_cbr_vault.test.status
}

data "huaweicloud_cbr_vaults" "filter_by_status" {
  depends_on = [huaweicloud_cbr_vault.test]

  status = local.vault_status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_cbr_vaults.filter_by_status.vaults[*].status : v == local.vault_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}
`, testAccDataVaults_basic_base(name))
}

func TestAccDataVaults_replicationServer(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_cbr_vaults.filter_by_type"
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
				Config: testAccDataVaults_filterByType(config),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataVaults_filterByType(config string) string {
	return fmt.Sprintf(`
%[1]s

# Filter by type
locals {
  vault_type = huaweicloud_cbr_vault.test.type
}

data "huaweicloud_cbr_vaults" "filter_by_type" {
  depends_on = [huaweicloud_cbr_vault.test]

  type = local.vault_type
}

locals {
  id_result_during_type_filter = [
    for v in data.huaweicloud_cbr_vaults.filter_by_type.vaults : v.id if v.type == local.vault_type
  ]
}

output "is_type_filter_useful" {
  value = contains(local.id_result_during_type_filter, huaweicloud_cbr_vault.test.id)
}
`, config)
}

func TestAccDataVaults_volume(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_cbr_vaults.filter_by_type"
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
				Config: testAccDataVaults_filterByType(config),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func TestAccDataVaults_backupTurbo(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_cbr_vaults.filter_by_type"
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
				Config: testAccDataVaults_filterByType(config),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func TestAccDataVaults_replicationTurbo(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_cbr_vaults.filter_by_type"
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
				Config: testAccDataVaults_filterByType(config),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}
