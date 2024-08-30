package dataarts

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getDataServiceApiStreamFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	workspaceId := state.Primary.Attributes["workspace_id"]
	dlmType := state.Primary.Attributes["dlm_type"]
	apiId := state.Primary.ID

	return dataarts.GetDataServiceApi(client, workspaceId, dlmType, apiId)
}

func TestAccDataServiceApi_basic(t *testing.T) {
	var (
		obj interface{}

		rName       = "huaweicloud_dataarts_dataservice_api.test"
		name        = acceptance.RandomAccResourceName()
		updateName  = acceptance.RandomAccResourceName()
		rc          = acceptance.InitResourceCheck(rName, &obj, getDataServiceApiStreamFunc)
		basicConfig = testAccDataServiceApi_base()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataServiceApi_basic_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "API_SPECIFIC_TYPE_CONFIGURATION"),
					resource.TestCheckResourceAttrPair(rName, "catalog_id",
						"huaweicloud_dataarts_dataservice_catalog.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "manager", acceptance.HW_DATAARTS_REVIEWER_NAME),
					resource.TestCheckResourceAttr(rName, "path", "/terraform/auto/resource_create/{resource_type}/{resource_name}"),
					resource.TestCheckResourceAttr(rName, "protocol", "PROTOCOL_TYPE_HTTPS"),
					resource.TestCheckResourceAttr(rName, "request_type", "REQUEST_TYPE_POST"),
					resource.TestCheckResourceAttr(rName, "request_params.#", "6"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.type", "DLI"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.connection_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttrPair(rName, "datasource_config.0.database",
						"huaweicloud_dli_database.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "datasource_config.0.datatable",
						"huaweicloud_dli_table.test", "name"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.queue", acceptance.HW_DATAARTS_DLI_QUEUE_NAME),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.0.name", "configuration"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.0.mapping", "configuration"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.0.condition", "CONDITION_TYPE_EQ"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.0.name", "resourceId"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.0.type", "REQUEST_PARAMETER_TYPE_STRING"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.0.field", "resource_id"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.0.description", "The resource ID, in UUID format"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.name", "bePlans"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.field", "plans"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.optional", "true"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.sort", "ASC"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.order", "1"),
					resource.TestCheckResourceAttrSet(rName, "create_user"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testAccDataServiceApi_basic_step2(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "catalog_id",
						"huaweicloud_dataarts_dataservice_catalog.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(rName, "manager", acceptance.HW_DATAARTS_REVIEWER_NAME),
					resource.TestCheckResourceAttr(rName, "path", "/terraform/auto/resource_query/{resource_type}/{resource_name}"),
					resource.TestCheckResourceAttr(rName, "protocol", "PROTOCOL_TYPE_HTTPS"),
					resource.TestCheckResourceAttr(rName, "request_type", "REQUEST_TYPE_GET"),
					resource.TestCheckResourceAttr(rName, "request_params.#", "2"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.type", "DLI"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.connection_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttrPair(rName, "datasource_config.0.database",
						"huaweicloud_dli_database.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "datasource_config.0.datatable",
						"huaweicloud_dli_table.test", "name"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.queue", acceptance.HW_DATAARTS_DLI_QUEUE_NAME),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.#", "0"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.#", "3"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.#", "0"),
					resource.TestCheckResourceAttrSet(rName, "create_user"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataServiceApiImportState(rName),
				ImportStateVerifyIgnore: []string{
					"auth_type",
					"catalog_id",
					"visibility",
				},
			},
		},
	})
}

func waitForDeletionCooldownComplete() resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		// After elastic resource pool is created, it cannot be deleted within one hour.
		// lintignore:R018
		time.Sleep(time.Hour)
		return nil
	}
}

func testDataServiceApiImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, dlmType, resourceId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		dlmType = rs.Primary.Attributes["dlm_type"]
		resourceId = rs.Primary.ID
		if workspaceId == "" || resourceId == "" {
			return "", fmt.Errorf("attribute 'workspace_id' or resource ID is missing")
		}
		if dlmType != "" {
			return fmt.Sprintf("%s/%s/%s", workspaceId, dlmType, resourceId), nil
		}
		return fmt.Sprintf("%s/%s", workspaceId, resourceId), nil
	}
}

func testAccDataServiceApi_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
// Under root path.
resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  count = 2

  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  name         = format("%[2]s_%%d", count.index)
}

resource "huaweicloud_dli_database" "test" {
  name = "%[2]s"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[2]s_resource_vpc"
  data_location = "DLI"

  columns {
    name        = "configuration"
    type        = "string"
    description = "The configuration for automatic creation, in JSON format"
  }
  columns {
    name        = "resource_id"
    type        = "string"
    description = "The resource ID, in UUID format"
  }
  columns {
    name        = "plans"
    type        = "string"
    description = "The plans to be executed"
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataServiceApi_basic_step1(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api" "test" {
  workspace_id = "%[2]s"
  dlm_type     = "EXCLUSIVE"

  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test[0].id
  name         = "%[3]s"
  description  = "Created by terraform script"
  auth_type    = "NONE"
  manager      = "%[4]s"
  path         = "/terraform/auto/resource_create/{resource_type}/{resource_name}"
  protocol     = "PROTOCOL_TYPE_HTTPS"
  request_type = "REQUEST_TYPE_POST"
  visibility   = "WORKSPACE"

  request_params {
    name          = "resource_type"
    position      = "REQUEST_PARAMETER_POSITION_PATH"
    type          = "REQUEST_PARAMETER_TYPE_STRING"
    description   = "The type of the terraform resource to be automatically created"
    necessary     = true
    example_value = "huaweicloud_vpc"
  }
  request_params {
    name          = "resource_name"
    position      = "REQUEST_PARAMETER_POSITION_PATH"
    type          = "REQUEST_PARAMETER_TYPE_STRING"
    description   = "The name of the terraform resource to be automatically created"
    necessary     = true
    example_value = "test"
  }
  request_params {
    name          = "count"
    position      = "REQUEST_PARAMETER_POSITION_QUERY"
    type          = "REQUEST_PARAMETER_TYPE_NUMBER"
    description   = "The name of the terraform resource to be automatically created"
    necessary     = false
    example_value = "3"
    default_value = "1"
  }
  request_params {
    name        = "configuration"
    position    = "REQUEST_PARAMETER_POSITION_BODY"
    type        = "REQUEST_PARAMETER_TYPE_STRING"
    description = "The configuration of the terraform resource, in JSON format"
    necessary   = true
  }
  request_params {
    name        = "resource_id"
    position    = "REQUEST_PARAMETER_POSITION_BODY"
    type        = "REQUEST_PARAMETER_TYPE_STRING"
    description = "The resource ID, in UUID format"
    necessary   = false
  }
  request_params {
    name          = "order"
    position      = "REQUEST_PARAMETER_POSITION_BODY"
    type          = "REQUEST_PARAMETER_TYPE_STRING"
    description   = "The filter parameter for resource configuration details"
    necessary     = false
    example_value = "asc"
    default_value = "desc"
  }

  datasource_config {
    type          = "DLI"
    connection_id = "%[5]s"
    database      = huaweicloud_dli_database.test.name
    datatable     = huaweicloud_dli_table.test.name
    queue         = "%[6]s"
    access_mode   = "SQL"

    backend_params {
      name      = "configuration"
      mapping   = "configuration"
      condition = "CONDITION_TYPE_EQ"
    }

    response_params {
      name        = "resourceId"
      type        = "REQUEST_PARAMETER_TYPE_STRING"
      field       = "resource_id"
      description = "The resource ID, in UUID format"
    }

    order_params {
      name     = "bePlans"
      field    = "plans"
      optional = true
      sort     = "ASC"
      order    = 1
    }
  }
}
`, basicConfig,
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}

func testAccDataServiceApi_basic_step2(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api" "test" {
  workspace_id = "%[2]s"
  dlm_type     = "EXCLUSIVE"

  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test[1].id
  name         = "%[3]s"
  description  = "Updated by terraform script"
  auth_type    = "IAM"
  manager      = "%[4]s"
  path         = "/terraform/auto/resource_query/{resource_type}/{resource_name}"
  protocol     = "PROTOCOL_TYPE_HTTPS"
  request_type = "REQUEST_TYPE_GET"
  visibility   = "PROJECT"

  request_params {
    name      = "resource_type"
    position  = "REQUEST_PARAMETER_POSITION_PATH"
    type      = "REQUEST_PARAMETER_TYPE_STRING"
    necessary = true
  }
  request_params {
    name          = "resource_name"
    position      = "REQUEST_PARAMETER_POSITION_PATH"
    type          = "REQUEST_PARAMETER_TYPE_STRING"
    description   = "The name of the terraform resource to be queried"
    necessary     = true
    example_value = "demo"
  }

  datasource_config {
    type          = "DLI"
    connection_id = "%[5]s"
    database      = huaweicloud_dli_database.test.name
    datatable     = huaweicloud_dli_table.test.name
    queue         = "%[6]s"

    response_params {
      name        = "configuration"
      type        = "REQUEST_PARAMETER_TYPE_STRING"
      field       = "configuration"
      description = "The resource configuration, in JSON format"
    }

    response_params {
      name        = "resource_id"
      type        = "REQUEST_PARAMETER_TYPE_STRING"
      field       = "resource_id"
      description = "The resource ID, in UUID format"
    }

    response_params {
      name        = "be_plans"
      type        = "REQUEST_PARAMETER_TYPE_STRING"
      field       = "plans"
      description = "The resource plan"
    }
  }
}
`, basicConfig,
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}

func TestAccDataServiceApi_simple(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dataarts_dataservice_api.test"
		name  = acceptance.RandomAccResourceName()
		rc    = acceptance.InitResourceCheck(rName, &obj, getDataServiceApiStreamFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create an API without request parameters.
				Config: testAccDataServiceApi_simple_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "API_SPECIFIC_TYPE_CONFIGURATION"),
					resource.TestCheckResourceAttrPair(rName, "catalog_id",
						"huaweicloud_dataarts_dataservice_catalog.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "manager", acceptance.HW_DATAARTS_REVIEWER_NAME),
					resource.TestCheckResourceAttr(rName, "path", fmt.Sprintf("/simple/%s", name)),
					resource.TestCheckResourceAttr(rName, "protocol", "PROTOCOL_TYPE_HTTPS"),
					resource.TestCheckResourceAttr(rName, "request_type", "REQUEST_TYPE_POST"),
					resource.TestCheckResourceAttr(rName, "request_params.#", "0"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.#", "0"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.0.name", "configuration"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.0.type", "REQUEST_PARAMETER_TYPE_STRING"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.0.field", "configuration"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.0.description", ""),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.#", "0"),
					resource.TestCheckResourceAttrSet(rName, "create_user"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				// Supports an request parameter and setting the corresponding backend parameter.
				Config: testAccDataServiceApi_simple_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "request_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "request_params.0.name", "configuration"),
					resource.TestCheckResourceAttr(rName, "request_params.0.position", "REQUEST_PARAMETER_POSITION_BODY"),
					resource.TestCheckResourceAttr(rName, "request_params.0.type", "REQUEST_PARAMETER_TYPE_STRING"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.0.name", "configuration"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.0.mapping", "configuration"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.0.condition", "CONDITION_TYPE_EQ"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.#", "0"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				// Remove the request parameter and update the field mapping to the order parameter.
				Config: testAccDataServiceApi_simple_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "request_params.#", "0"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.backend_params.#", "0"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.response_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.name", "order_configuration"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.field", "configuration"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.optional", "true"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.sort", "ASC"),
					resource.TestCheckResourceAttr(rName, "datasource_config.0.order_params.0.order", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataServiceApiImportState(rName),
				ImportStateVerifyIgnore: []string{
					"auth_type",
					"catalog_id",
					"visibility",
				},
			},
		},
	})
}

func testAccDataServiceApi_simple_base(name string) string {
	return fmt.Sprintf(`
// Under root path.
resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  name         = "%[2]s"
}

resource "huaweicloud_dli_database" "test" {
  name = "%[2]s"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[2]s"
  data_location = "DLI"

  columns {
    name        = "configuration"
    type        = "string"
    description = "The configuration for automatic creation, in JSON format"
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataServiceApi_simple_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api" "test" {
  workspace_id = "%[2]s"
  dlm_type     = "EXCLUSIVE"

  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
  name         = "%[3]s"
  auth_type    = "APP"
  manager      = "%[4]s"
  path         = "/simple/%[3]s"
  protocol     = "PROTOCOL_TYPE_HTTPS"
  request_type = "REQUEST_TYPE_POST"
  visibility   = "PROJECT"

  datasource_config {
    type          = "DLI"
    connection_id = "%[5]s"
    database      = huaweicloud_dli_database.test.name
    datatable     = huaweicloud_dli_table.test.name
    queue         = "%[6]s"

    response_params {
      name  = "configuration"
      type  = "REQUEST_PARAMETER_TYPE_STRING"
      field = "configuration"
    }
  }
}
`, testAccDataServiceApi_simple_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}

func testAccDataServiceApi_simple_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api" "test" {
  workspace_id = "%[2]s"
  dlm_type     = "EXCLUSIVE"

  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
  name         = "%[3]s"
  auth_type    = "APP"
  manager      = "%[4]s"
  path         = "/simple/%[3]s"
  protocol     = "PROTOCOL_TYPE_HTTPS"
  request_type = "REQUEST_TYPE_POST"
  visibility   = "PROJECT"

  request_params {
    name     = "configuration"
    position = "REQUEST_PARAMETER_POSITION_BODY"
    type     = "REQUEST_PARAMETER_TYPE_STRING"
  }

  datasource_config {
    type          = "DLI"
    connection_id = "%[5]s"
    database      = huaweicloud_dli_database.test.name
    datatable     = huaweicloud_dli_table.test.name
    queue         = "%[6]s"

    backend_params {
      name      = "configuration"
      mapping   = "configuration"
      condition = "CONDITION_TYPE_EQ"
    }
    response_params {
      name  = "configuration"
      type  = "REQUEST_PARAMETER_TYPE_STRING"
      field = "configuration"
    }
  }
}
`, testAccDataServiceApi_simple_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}

func testAccDataServiceApi_simple_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api" "test" {
  workspace_id = "%[2]s"
  dlm_type     = "EXCLUSIVE"

  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
  name         = "%[3]s"
  auth_type    = "APP"
  manager      = "%[4]s"
  path         = "/simple/%[3]s"
  protocol     = "PROTOCOL_TYPE_HTTPS"
  request_type = "REQUEST_TYPE_GET"
  visibility   = "PROJECT"

  datasource_config {
    type          = "DLI"
    connection_id = "%[5]s"
    database      = huaweicloud_dli_database.test.name
    datatable     = huaweicloud_dli_table.test.name
    queue         = "%[6]s"

    response_params {
      name  = "configuration"
      type  = "REQUEST_PARAMETER_TYPE_STRING"
      field = "configuration"
    }
    order_params {
      name     = "order_configuration"
      field    = "configuration"
      optional = true
      sort     = "ASC"
      order    = 1
    }
  }
}
`, testAccDataServiceApi_simple_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_REVIEWER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME)
}
