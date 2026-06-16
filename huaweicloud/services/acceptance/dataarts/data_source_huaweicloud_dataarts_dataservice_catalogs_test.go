package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDataServiceCatalogs_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_dataservice_catalogs.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byCatalogId   = "data.huaweicloud_dataarts_dataservice_catalogs.test"
		dcByCatalogId = acceptance.InitDataSourceCheck(byCatalogId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataDataServiceCatalogs_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts DataService catalogs"),
			},
			{
				Config: testAccDataSourceDataServiceCatalogs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "catalogs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					dcByCatalogId.CheckResourceExists(),
					resource.TestMatchResourceAttr(byCatalogId, "catalogs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byCatalogId, "catalogs.0.id"),
					resource.TestCheckResourceAttrSet(byCatalogId, "catalogs.0.parent_id"),
					resource.TestCheckResourceAttrSet(byCatalogId, "catalogs.0.name"),
					resource.TestCheckResourceAttrSet(byCatalogId, "catalogs.0.description"),
					resource.TestCheckResourceAttrSet(byCatalogId, "catalogs.0.api_catalog_type"),
					resource.TestCheckResourceAttrSet(byCatalogId, "catalogs.0.created_at"),
					resource.TestCheckResourceAttrSet(byCatalogId, "catalogs.0.create_user"),
					resource.TestCheckResourceAttrSet(byCatalogId, "catalogs.0.updated_at"),
					resource.TestCheckResourceAttrSet(byCatalogId, "catalogs.0.update_user"),
				),
			},
		},
	})
}

func testAccDataDataServiceCatalogs_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_dataservice_catalogs" "test" {
  workspace_id = "%[1]s"
  catalog_id   = "%[1]s"
}
`, randUUID.String())
}

func testAccDataSourceDataServiceCatalogs_basic_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  parent_id    = "0"
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  name         = "%[2]s"
  description  = "Created by terraform script"
}

resource "huaweicloud_dataarts_dataservice_catalog" "child" {
  parent_id    = huaweicloud_dataarts_dataservice_catalog.test.id
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  name         = "%[2]s_child"
  description  = "Created by terraform script!"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataSourceDataServiceCatalogs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_dataservice_catalogs" "all" {
  depends_on = [
    huaweicloud_dataarts_dataservice_catalog.child,
  ]

  workspace_id = "%[2]s"
}

data "huaweicloud_dataarts_dataservice_catalogs" "test" {
  depends_on = [
    data.huaweicloud_dataarts_dataservice_catalogs.all,
  ]

  workspace_id = "%[2]s"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
}
`, testAccDataSourceDataServiceCatalogs_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
