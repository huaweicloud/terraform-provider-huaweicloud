package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getDataServiceApiPublishInfoFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	workspaceId := state.Primary.Attributes["workspace_id"]
	instanceId := state.Primary.Attributes["instance_id"]
	apiId := state.Primary.Attributes["api_id"]

	return dataarts.QueryApiPublishInfoByInstanceId(client, workspaceId, apiId, instanceId)
}

func TestAccDataServiceApiPublishment_basic(t *testing.T) {
	var (
		obj interface{}

		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_dataarts_dataservice_api_publishment.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDataServiceApiPublishInfoFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
			acceptance.TestAccPreCheckDataArtsDataServiceApigInstanceId(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Just test whether the publish request was successful.
				Config: testAccDataServiceApiPublish_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "publish_id"),
				),
			},
			{
				// Unpublish the API and test whether CheckDeleted logic is useful
				Config:             testAccDataServiceApiPublish_basic_step2(name),
				ExpectNonEmptyPlan: true,
			},
			{
				// Republish the API and publish it in the APIG service at the same time.
				// And test whether the publish request was successful and the API been published successfully on the
				// APIG side.
				Config: testAccDataServiceApiPublish_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "publish_id"),
					resource.TestCheckOutput("total_api_count_in_apig_side", "1"),
				),
			},
		},
	})
}

// Debug and publish the API.
func testAccDataServiceApiPublish_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"

  # Unpublish API on the Data Service will not unpublish the APIG API and will caused the delete operation of APIG group
  # triggers error return.
  # Force destroy allows deleting all APIs under the management of this group (even published APIs).
  force_destroy = true
}

data "huaweicloud_dataarts_dataservice_instances" "test" {
  workspace_id = "%[3]s"
}

resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  workspace_id = "%[3]s"
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

  workspace_id = "%[3]s"
  dlm_type     = "EXCLUSIVE"

  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
  name         = "%[2]s"
  auth_type    = "NONE"
  manager      = "%[4]s"
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
    connection_id = "%[5]s"
    queue         = "%[6]s"
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
  workspace_id = "%[3]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")

  params = jsonencode({
    "page_num": "1",
    "page_size": "100",
    "test_request_field": "test"
  })
  max_retries = 5
}
`, acceptance.HW_DATAARTS_INSTANCE_ID_IN_APIG,
		name,
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}

func testAccDataServiceApiPublish_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_publishment" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
}
`, testAccDataServiceApiPublish_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDataServiceApiPublish_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_publishment" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
}

resource "huaweicloud_dataarts_dataservice_api_action" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  type        = "UNPUBLISH"
}
`, testAccDataServiceApiPublish_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDataServiceApiPublish_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_publishment" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]

  workspace_id = "%[2]s"

  api_id           = huaweicloud_dataarts_dataservice_api.test.id
  instance_id      = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  apig_type        = "APIGW"
  apig_instance_id = "%[3]s"
  apig_group_id    = huaweicloud_apig_group.test.id
}

data "huaweicloud_apig_api_basic_configurations" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]

  instance_id = "%[3]s"
  group_id    = huaweicloud_apig_group.test.id
}

output "total_api_count_in_apig_side" {
  value = length(data.huaweicloud_apig_api_basic_configurations.test.configurations)
}
`, testAccDataServiceApiPublish_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_INSTANCE_ID_IN_APIG)
}

func TestAccDataServiceApiPublishment_action(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
			acceptance.TestAccPreCheckDataArtsDataServiceApigInstanceId(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Test whether the publish request (on Data Service side) was successful.
				Config: testAccDataServiceApiPublish_action_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("total_api_count_in_apig_side", "0"),
				),
			},
			{
				// Test whether repeated publish triggers error return.
				Config:      testAccDataServiceApiPublish_action_step2(name),
				ExpectError: regexp.MustCompile(`Api path already exist`),
			},
			{
				// Test whether the unpublish request was successful.
				Config: testAccDataServiceApiPublish_action_step3(name),
			},
			{
				// Test whether repeated unpublish triggers error return.
				Config:      testAccDataServiceApiPublish_action_step4(name),
				ExpectError: regexp.MustCompile(`This action is illegal`),
			},
			{
				// Test whether the publish request (both published on Data Service side and APIG side) was successful.
				Config: testAccDataServiceApiPublish_action_step5(name),
				Check: resource.ComposeTestCheckFunc(
					// Check publish result on the APIG side.
					resource.TestCheckOutput("total_api_count_in_apig_side", "1"),
				),
			},
			{
				// Test whether the unpublish request was successful.
				Config: testAccDataServiceApiPublish_action_step6(name),
			},
		},
	})
}

// Publish the API (just on Data Service side).
func testAccDataServiceApiPublish_action_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_publish" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
}

data "huaweicloud_apig_api_basic_configurations" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publish.test]

  instance_id = "%[3]s"
  group_id    = huaweicloud_apig_group.test.id
}

output "total_api_count_in_apig_side" {
  value = length(data.huaweicloud_apig_api_basic_configurations.test.configurations)
}
`, testAccDataServiceApiPublish_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_INSTANCE_ID_IN_APIG)
}

// Publish the API again and trigger the error return.
func testAccDataServiceApiPublish_action_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_publish" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
}

resource "huaweicloud_dataarts_dataservice_api_publish" "publish_again" {
  depends_on = [
    huaweicloud_dataarts_dataservice_api_debug.test,
    huaweicloud_dataarts_dataservice_api_publish.test,
  ]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
}
`, testAccDataServiceApiPublish_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}

// Unpublish the API.
func testAccDataServiceApiPublish_action_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_action" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  type        = "UNPUBLISH"
}
`, testAccDataServiceApiPublish_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}

// Unpublish the API again and trigger the error return.
func testAccDataServiceApiPublish_action_step4(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_action" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  type        = "UNPUBLISH"
}

resource "huaweicloud_dataarts_dataservice_api_action" "unpublish_again" {
  depends_on = [
    huaweicloud_dataarts_dataservice_api_debug.test,
    huaweicloud_dataarts_dataservice_api_action.test,
  ]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  type        = "UNPUBLISH"
}
`, testAccDataServiceApiPublish_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}

// Publish the API (Both on Data Service and APIG sides).
func testAccDataServiceApiPublish_action_step5(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_publish" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_debug.test]

  workspace_id = "%[2]s"

  api_id           = huaweicloud_dataarts_dataservice_api.test.id
  instance_id      = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  apig_type        = "APIGW"
  apig_instance_id = "%[3]s"
  apig_group_id    = huaweicloud_apig_group.test.id
}

data "huaweicloud_apig_api_basic_configurations" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publish.test]

  instance_id = "%[3]s"
  group_id    = huaweicloud_apig_group.test.id
}

output "total_api_count_in_apig_side" {
  value = length(data.huaweicloud_apig_api_basic_configurations.test.configurations)
}
`, testAccDataServiceApiPublish_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_INSTANCE_ID_IN_APIG)
}

// Unpublish the API (Just from Data Service side).
func testAccDataServiceApiPublish_action_step6(name string) string {
	return testAccDataServiceApiPublish_action_step3(name)
}
