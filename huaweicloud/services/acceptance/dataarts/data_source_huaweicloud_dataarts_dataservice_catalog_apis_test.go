package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Query exclusive API that they are not published.
func TestAccDataSourceDataServiceCatalogApis_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_dataarts_dataservice_catalog_apis.test"
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataServiceCatalogApis_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "apis.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(rName, "apis.0.id"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.name"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.description"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.type"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.manager"),
					resource.TestMatchResourceAttr(rName, "apis.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceDataServiceCatalogApis_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  name = "%[1]s"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[1]s"
  data_location = "DLI"

  columns {
    name        = "configuration"
    type        = "string"
    description = "The configuration for automatic creation, in JSON format"
  }
}

// Under root path.
resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  workspace_id = "%[2]s"
  dlm_type     = "EXCLUSIVE"

  name = "%[1]s"
}

resource "huaweicloud_dataarts_dataservice_api" "test" {
  count = 2

  workspace_id = "%[2]s"
  dlm_type     = "EXCLUSIVE"

  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
  name         = format("%[1]s_%%d", count.index)
  description  = "Created by terraform script"
  auth_type    = "NONE"
  manager      = "%[3]s"
  path         = format("/%[1]s/%%d", count.index)
  protocol     = "PROTOCOL_TYPE_HTTPS"
  request_type = "REQUEST_TYPE_GET"
  visibility   = "PROJECT"

  datasource_config {
    type          = "DLI"
    connection_id = "%[4]s"
    database      = huaweicloud_dli_database.test.name
    datatable     = huaweicloud_dli_table.test.name
    queue         = "%[5]s"

    response_params {
      name  = "configuration"
      type  = "REQUEST_PARAMETER_TYPE_STRING"
      field = "configuration"
    }
  }
}

data "huaweicloud_dataarts_dataservice_catalog_apis" "test" {
  depends_on = [
    huaweicloud_dataarts_dataservice_api.test
  ]

  workspace_id = "%[2]s"
  dlm_type     = "EXCLUSIVE"

  catalog_id = huaweicloud_dataarts_dataservice_catalog.test.id
}
`, name,
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}
