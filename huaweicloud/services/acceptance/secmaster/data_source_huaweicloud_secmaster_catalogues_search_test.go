package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterCataloguesSearch_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_catalogues_search.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterCataloguesSearch_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.catalogue_address"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.catalogue_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_card_area"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_display"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_landing_page"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_navigation"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.parent_catalogue"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.publisher_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.parent_alias_en"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.parent_alias_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.second_alias_en"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.second_alias_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.second_catalogue"),

					resource.TestCheckOutput("is_second_catalogue_filter_useful", "true"),
					resource.TestCheckOutput("is_publisher_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterCataloguesSearch_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_catalogues_search" "test" {
  workspace_id = "%[1]s"
}

# Filter by second_catalogue
locals {
  second_catalogue = data.huaweicloud_secmaster_catalogues_search.test.data[0].second_catalogue
}

data "huaweicloud_secmaster_catalogues_search" "filter_by_second_catalogue" {
  workspace_id     = "%[1]s"
  second_catalogue = local.second_catalogue
}

locals {
  list_by_second_catalogue = data.huaweicloud_secmaster_catalogues_search.filter_by_second_catalogue.data
}

output "is_second_catalogue_filter_useful" {
  value = length(local.list_by_second_catalogue) > 0 && alltrue(
    [for v in local.list_by_second_catalogue[*].second_catalogue : v == local.second_catalogue]
  )
}

# Filter by publisher_name
locals {
  publisher_name = data.huaweicloud_secmaster_catalogues_search.test.data[0].publisher_name
}

data "huaweicloud_secmaster_catalogues_search" "filter_by_publisher_name" {
  workspace_id   = "%[1]s"
  publisher_name = local.publisher_name
}

locals {
  list_by_publisher_name = data.huaweicloud_secmaster_catalogues_search.filter_by_publisher_name.data
}

output "is_publisher_name_filter_useful" {
  value = length(local.list_by_publisher_name) > 0 && alltrue(
    [for v in local.list_by_publisher_name[*].publisher_name : v == local.publisher_name]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
