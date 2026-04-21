package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDataServiceCatalogApis_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_dataservice_catalog_apis.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDataServiceCatalogApis_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttrSet(all, "region"),
					resource.TestMatchResourceAttr(all, "apis.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "apis.0.id"),
					resource.TestCheckResourceAttrSet(all, "apis.0.name"),
					resource.TestCheckResourceAttrSet(all, "apis.0.manager"),
					resource.TestCheckResourceAttrSet(all, "apis.0.type"),
					resource.TestCheckResourceAttrSet(all, "apis.0.description"),
					resource.TestMatchResourceAttr(all, "apis.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataDataServiceCatalogApis_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_dataservice_catalog_apis" "all" {
  workspace_id = "%[1]s"
  catalog_id   = "0" # the ID of root catalog is '0'
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
