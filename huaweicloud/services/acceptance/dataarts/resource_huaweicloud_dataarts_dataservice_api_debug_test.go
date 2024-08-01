package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataServiceApiDebug_basic(t *testing.T) {
	rName := "huaweicloud_dataarts_dataservice_api_debug.test"

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataServiceApiDebug_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "url", "/demo"),
					resource.TestCheckResourceAttrSet(rName, "result"),
					resource.TestCheckResourceAttrSet(rName, "timeout"),
					resource.TestCheckResourceAttrSet(rName, "request_header"),
					resource.TestCheckResourceAttrSet(rName, "response_header"),
				),
			},
		},
	})
}

// Deploy an API that can be used for debugging, publishing, and authorization.
func testAccDataServiceApi_develop() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_dataservice_instances" "test" {
  workspace_id = "%[1]s"
}

# Under root path.
resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  dlm_type     = "EXCLUSIVE"
}

resource "huaweicloud_dli_database" "test" {
  name = "%[2]s"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[2]s"
  data_location = "DLI"

  columns {
    name = "test_request_field"
    type = "string"
  }
  columns {
    name = "test_response_field"
    type = "string"
  }
}

resource "huaweicloud_dataarts_dataservice_api" "test" {
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"

  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
  name         = "%[2]s"
  auth_type    = "NONE"
  manager      = "%[3]s"
  path         = "/demo"
  protocol     = "PROTOCOL_TYPE_HTTPS"
  request_type = "REQUEST_TYPE_POST"
  visibility   = "WORKSPACE"

  request_params {
    name      = "test_request_field"
    position  = "REQUEST_PARAMETER_POSITION_BODY"
    type      = "REQUEST_PARAMETER_TYPE_STRING"
    necessary = true
  }

  datasource_config {
    type          = "DLI"
    connection_id = "%[4]s"
    queue         = "%[5]s"
    database      = huaweicloud_dli_database.test.name
    datatable     = huaweicloud_dli_table.test.name

    backend_params {
      name     = "test_request_field"
      mapping  = "test_request_field"
      condition = "CONDITION_TYPE_EQ"
    }

    response_params {
      name  = "test_response_field"
      type  = "REQUEST_PARAMETER_TYPE_STRING"
      field = "test_response_field"
    }
  }
}

`, acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}

func testAccDataServiceApiDebug_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_debug" "test" {
  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")

  params = jsonencode({
    "page_num": "1",
    "page_size": "100",
    "test_request_field": "{\"foo\": \"bar\"}"
  })
}
`, testAccDataServiceApi_develop(), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
