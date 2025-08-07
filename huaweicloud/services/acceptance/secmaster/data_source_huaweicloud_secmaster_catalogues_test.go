package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterCatalogues_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_catalogues.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterCatalogues_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.catalogue_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_card_area"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_display"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_landing_page"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_navigation"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.parent_catalogue"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.second_catalogue"),
				),
			},
		},
	})
}

func testDataSourceSecmasterCatalogues_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_catalogues" "test" {
  workspace_id = "%[1]s"
}

# Fields catalogue_type and catalogue_code are only used to perform the query param code.
# In fact, the data returned by the values of these two fields is empty. 
data "huaweicloud_secmaster_catalogues" "test-filter" {
  workspace_id   = "%[1]s"
  catalogue_type = "demo-type"
  catalogue_code = "demo-code"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
