package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceExternalVaults_basic(t *testing.T) {
	var (
		dataSourceName   = "data.huaweicloud_cbr_external_vaults.test"
		dc               = acceptance.InitDataSourceCheck(dataSourceName)
		dataSourceWithId = "data.huaweicloud_cbr_external_vaults.test_with_id"
		dcWithId         = acceptance.InitDataSourceCheck(dataSourceWithId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCBRRegionName(t)
			acceptance.TestAccPreCheckCBRExternalProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceExternalVaults_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.bind_rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.bind_rules.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.allocated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.used"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.cloud_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.consistent_level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.protect_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.object_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.billing.0.is_multi_az"),

					dcWithId.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceWithId, "vaults.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceWithId, "vaults.0.id", "data.huaweicloud_cbr_external_vaults.test", "vaults.0.id"),

					resource.TestCheckOutput("is_vault_id_useful", "true"),
					resource.TestCheckOutput("is_cloud_type_useful", "true"),
					resource.TestCheckOutput("is_protect_type_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceExternalVaults_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cbr_external_vaults" "test" {
  external_project_id = "%[1]s"
  region_id           = "%[2]s"
}

# Filter using vault_id.
locals {
  vault_id = data.huaweicloud_cbr_external_vaults.test.vaults.0.id
}

data "huaweicloud_cbr_external_vaults" "test_with_id" {
  external_project_id = "%[1]s"
  region_id           = "%[2]s"
  vault_id            = local.vault_id
}

output "is_vault_id_useful" {
  value = length(data.huaweicloud_cbr_external_vaults.test_with_id.vaults.*.id) > 0 && alltrue(
    [for v in data.huaweicloud_cbr_external_vaults.test_with_id.vaults.*.id :v == local.vault_id]
  )
}

# Filter using cloud_type.
locals {
  cloud_type = data.huaweicloud_cbr_external_vaults.test.vaults.0.billing.0.cloud_type
}

data "huaweicloud_cbr_external_vaults" "test_with_cloud_type" {
  external_project_id = "%[1]s"
  region_id           = "%[2]s"
  cloud_type          = local.cloud_type
}

output "is_cloud_type_useful" {
  value = length(data.huaweicloud_cbr_external_vaults.test_with_cloud_type.vaults.*.billing.0.cloud_type) > 0 && alltrue(
    [for v in data.huaweicloud_cbr_external_vaults.test_with_cloud_type.vaults.*.billing.0.cloud_type :v == local.cloud_type]
  )
}

# Filter using protect_type.
locals {
  protect_type = data.huaweicloud_cbr_external_vaults.test.vaults.0.billing.0.protect_type
}

data "huaweicloud_cbr_external_vaults" "test_with_protect_type" {
  external_project_id = "%[1]s"
  region_id           = "%[2]s"
  protect_type        = local.protect_type
}

output "is_protect_type_useful" {
  value = length(data.huaweicloud_cbr_external_vaults.test_with_protect_type.vaults.*.billing.0.protect_type) > 0 && alltrue(
    [for v in data.huaweicloud_cbr_external_vaults.test_with_cloud_type.vaults.*.billing.0.protect_type :v == local.protect_type]
  )
}
`, acceptance.HW_CBR_EXTERNAL_PROJECT_ID, acceptance.HW_REGION_NAME_1)
}
