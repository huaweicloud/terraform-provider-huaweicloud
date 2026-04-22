package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDataServiceAppAuthorizedApis_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_dataservice_app_authorized_apis.test"
		dc  = acceptance.InitDataSourceCheck(all)
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
				Config: testAccDataDataServiceAppAuthorizedApis_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "apis.#", "2"),
					resource.TestCheckResourceAttrSet(all, "apis.0.id"),
					resource.TestCheckResourceAttrSet(all, "apis.0.name"),
					resource.TestCheckResourceAttrSet(all, "apis.0.approval_time"),
					resource.TestCheckResourceAttrSet(all, "apis.0.manager"),
					resource.TestCheckResourceAttrSet(all, "apis.0.relationship_type"),
				),
			},
			{
				// Require authorization revocation before resource deletion.
				Config: testAccDataDataServiceAppAuthorizedApis_basic_step2(name),
			},
		},
	})
}

func testAccDataDataServiceAppAuthorizedApis_base(name string) string {
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
  count = 2

  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
  name         = format("%[2]s_%%d", count.index)
  auth_type    = "APP"
  manager      = "%[3]s"
  path         = format("/%[2]s/test_%%d", count.index)
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
      name      = "test_request_field"
      mapping   = "test_request_field"
      condition = "CONDITION_TYPE_EQ"
    }
    response_params {
      name  = "test_response_field"
      type  = "REQUEST_PARAMETER_TYPE_STRING"
      field = "test_response_field"
    }
  }

  depends_on = [huaweicloud_dli_table.test]
}

resource "huaweicloud_dataarts_dataservice_api_debug" "test" {
  count = 2

  workspace_id = "%[1]s"
  api_id       = huaweicloud_dataarts_dataservice_api.test[count.index].id
  instance_id  = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  max_retries  = 6

  params = jsonencode({
    "page_num": "1",
    "page_size": "100",
    "test_request_field": "{\"foo\": \"bar\"}"
  })
}

resource "huaweicloud_dataarts_dataservice_api_publishment" "test" {
  count = 2

  workspace_id = "%[1]s"
  api_id       = huaweicloud_dataarts_dataservice_api.test[count.index].id
  instance_id  = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")

  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]
}

resource "huaweicloud_dataarts_dataservice_app" "test" {
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  name         = "%[2]s"
  app_type     = "APP"

  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}

func testAccDataDataServiceAppAuthorizedApis_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_auth" "test" {
  count = 2

  workspace_id = "%[2]s"
  api_id       = huaweicloud_dataarts_dataservice_api.test[count.index].id
  instance_id  = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  app_ids      = [huaweicloud_dataarts_dataservice_app.test.id]

  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]
}

data "huaweicloud_dataarts_dataservice_app_authorized_apis" "test" {
  workspace_id = "%[2]s"
  app_id       = huaweicloud_dataarts_dataservice_app.test.id

  depends_on = [huaweicloud_dataarts_dataservice_api_auth.test]
}
`, testAccDataDataServiceAppAuthorizedApis_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDataDataServiceAppAuthorizedApis_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_auth_action" "test" {
  count = 2

  workspace_id = "%[2]s"
  api_id       = huaweicloud_dataarts_dataservice_api.test[count.index].id
  instance_id  = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  app_id       = huaweicloud_dataarts_dataservice_app.test.id
  type         = "APPLY_TYPE_APP_CANCEL_AUTHORIZE"
}
`, testAccDataDataServiceAppAuthorizedApis_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
