package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before executing this use case, please create several pieces of data first.
func TestAccDatasourceRAMSharedPrincipals_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ram_shared_principals.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byPrincipal   = "data.huaweicloud_ram_shared_principals.principal_filter"
		dcByPrincipal = acceptance.InitDataSourceCheck(byPrincipal)

		byResourceUrn   = "data.huaweicloud_ram_shared_principals.resource_urn_filter"
		dcByResourceUrn = acceptance.InitDataSourceCheck(byResourceUrn)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMEnableFlag(t)
			acceptance.TestAccPreCheckRAMSharedPrincipalsQueryFields(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceRAMSharedPrincipals_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "shared_principals.0.resource_share_id"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_principals.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_principals.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_principals.0.updated_at"),

					dcByPrincipal.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byPrincipal, "shared_principals.#"),

					dcByResourceUrn.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byResourceUrn, "shared_principals.#"),

					resource.TestCheckOutput("resource_share_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceRAMSharedPrincipals_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ram_shared_principals" "test" {
  resource_owner = "self"
}

data "huaweicloud_ram_shared_principals" "principal_filter" {
  resource_owner = "self"
  principal      = "%[1]s"
}

data "huaweicloud_ram_shared_principals" "resource_urn_filter" {
  resource_owner = "self"
  resource_urn   = "%[2]s"
}

# Filter by resource_share_id
locals {
  resource_share_id = data.huaweicloud_ram_shared_principals.test.shared_principals.0.resource_share_id
}

data "huaweicloud_ram_shared_principals" "resource_share_id_filter" {
  resource_owner    = "self"
  resource_share_id = local.resource_share_id
}

locals {
  resource_share_id_filter_result = [
    for v in data.huaweicloud_ram_shared_principals.resource_share_id_filter.shared_principals[*].resource_share_id : v == local.resource_share_id
  ]
}

output "resource_share_id_filter_is_useful" {
  value = length(local.resource_share_id_filter_result) > 0 && alltrue(local.resource_share_id_filter_result)
}
`, acceptance.HW_RAM_SHARE_ACCOUNT_ID, acceptance.HW_RAM_SHARE_RESOURCE_URN)
}
