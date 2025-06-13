package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceExternalVaults_basic(t *testing.T) {
	var (
		dataSourceName              = "data.huaweicloud_cbr_external_vaults.test"
		dc                          = acceptance.InitDataSourceCheck(dataSourceName)
		dataSourceNameWithId        = "data.huaweicloud_cbr_external_vaults.test_with_id"
		dcWithId                    = acceptance.InitDataSourceCheck(dataSourceNameWithId)
		dataSourceNameWithFilers    = "data.huaweicloud_cbr_external_vaults.test_with_filters"
		dcWithFilers                = acceptance.InitDataSourceCheck(dataSourceNameWithFilers)
		dataSourceNameAllAttributes = "data.huaweicloud_cbr_external_vaults.test_all_attributes"
		dcAllAttributes             = acceptance.InitDataSourceCheck(dataSourceNameAllAttributes)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceExternalVaults_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.provider_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.auto_bind"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.auto_expand"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.smn_notify"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.threshold"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vaults.0.locked"),
				),
			},
			{
				Config: testAccDataSourceExternalVaults_withVaultID(),
				Check: resource.ComposeTestCheckFunc(
					dcWithId.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceNameWithId, "vaults.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceNameWithId, "vaults.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceNameWithId, "vaults.0.name"),
				),
			},
			{
				Config: testAccDataSourceExternalVaults_withFilters(),
				Check: resource.ComposeTestCheckFunc(
					dcWithFilers.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceNameWithFilers, "vaults.#"),
				),
			},
			{
				Config: testAccDataSourceExternalVaults_validateAllAttributes(),
				Check: resource.ComposeTestCheckFunc(
					dcAllAttributes.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.bind_rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.bind_rules.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.resources.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.resources.0.extra_info.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.allocated"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.used"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.cloud_type"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.consistent_level"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.protect_type"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.object_type"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceNameAllAttributes, "vaults.0.billing.0.is_multi_az"),
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
`, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME_1)
}

func testAccDataSourceExternalVaults_withVaultID() string {
	return fmt.Sprintf(`
data "huaweicloud_cbr_external_vaults" "test" {
  external_project_id = "%[1]s"
  region_id           = "%[2]s"
}

data "huaweicloud_cbr_external_vaults" "test_with_id" {
  external_project_id = "%[1]s"
  region_id           = "%[2]s"
  vault_id            = data.huaweicloud_cbr_external_vaults.test.vaults.0.id
}
`, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME_1)
}

func testAccDataSourceExternalVaults_withFilters() string {
	return fmt.Sprintf(`
data "huaweicloud_cbr_external_vaults" "test_with_filters" {
  external_project_id = "%[1]s"
  region_id           = "%[2]s"
  cloud_type          = "public"
  object_type         = "server"
}
`, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME_1)
}

func testAccDataSourceExternalVaults_validateAllAttributes() string {
	return fmt.Sprintf(`
data "huaweicloud_cbr_external_vaults" "test_all_attributes" {
  external_project_id = "%[1]s"
  region_id           = "%[2]s"
}
`, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME_1)
}
