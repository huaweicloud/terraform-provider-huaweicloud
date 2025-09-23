package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataServiceApiAuth_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
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
				// Just test whether the authorize request was executed successful.
				Config: testAccDataServiceApiAuth_basic_step1(name),
			},
			{
				// Just test whether the authorize status was cancelled successful.
				Config: testAccDataServiceApiAuth_basic_step2(name),
			},
		},
	})
}

// Debug and publish the API.
func testAccDataServiceApiAuth_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_dataservice_instances" "test" {
  workspace_id = "%[1]s"
}

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
  depends_on = [huaweicloud_dli_table.test]

  workspace_id = "%[1]s"

  dlm_type     = "EXCLUSIVE"
  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
  name         = "%[2]s"
  auth_type    = "APP"
  manager      = "%[3]s"
  path         = "/%[2]s/test"
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

resource "huaweicloud_dataarts_dataservice_api_debug" "test" {
  workspace_id = "%[1]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")

  params = jsonencode({
    "page_num": "1",
    "page_size": "100",
    "test_request_field": "{\"foo\": \"bar\"}"
  })

  max_retries = 6
}

resource "huaweicloud_dataarts_dataservice_api_publishment" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]

  workspace_id = "%[1]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
}

resource "huaweicloud_dataarts_dataservice_app" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]
  count      = 2

  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"

  name     = format("%[2]s_%%d", count.index)
  app_type = "APP"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}

// Authorize some APPs to access the API.
func testAccDataServiceApiAuth_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_auth" "test" {
  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  app_ids     = slice(huaweicloud_dataarts_dataservice_app.test[*].id, 0, 2)
}
`, testAccDataServiceApiAuth_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}

// Cancel the API access permissions for all APPs.
func testAccDataServiceApiAuth_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_auth_action" "test" {
  count = 2

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  app_id      = element(huaweicloud_dataarts_dataservice_app.test[*].id, count.index)
  type        = "APPLY_TYPE_APP_CANCEL_AUTHORIZE"
}
`, testAccDataServiceApiAuth_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}
