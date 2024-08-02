package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataServiceApis_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_dataarts_dataservice_apis.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byId   = "data.huaweicloud_dataarts_dataservice_apis.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_dataarts_dataservice_apis.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNameNotFound   = "data.huaweicloud_dataarts_dataservice_apis.filter_by_name_not_found"
		dcByNameNotFound = acceptance.InitDataSourceCheck(byNameNotFound)

		byDesc   = "data.huaweicloud_dataarts_dataservice_apis.filter_by_desc"
		dcByDesc = acceptance.InitDataSourceCheck(byDesc)

		byCreator   = "data.huaweicloud_dataarts_dataservice_apis.filter_by_creator"
		dcByCreator = acceptance.InitDataSourceCheck(byCreator)

		byType   = "data.huaweicloud_dataarts_dataservice_apis.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byTableName   = "data.huaweicloud_dataarts_dataservice_apis.filter_by_datatable"
		dcByTableName = acceptance.InitDataSourceCheck(byTableName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataServiceApis_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "apis.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckResourceAttr(byId, "apis.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.id", "huaweicloud_dataarts_dataservice_api.test", "id"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.name", "huaweicloud_dataarts_dataservice_api.test", "name"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.type", "huaweicloud_dataarts_dataservice_api.test", "type"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.description", "huaweicloud_dataarts_dataservice_api.test", "description"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.protocol", "huaweicloud_dataarts_dataservice_api.test", "protocol"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.request_type", "huaweicloud_dataarts_dataservice_api.test", "request_type"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.manager", "huaweicloud_dataarts_dataservice_api.test", "manager"),
					resource.TestCheckResourceAttr(byId, "apis.0.datasource_config.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.datasource_config.0.type",
						"huaweicloud_dataarts_dataservice_api.test", "datasource_config.0.type"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.datasource_config.0.connection_id",
						"huaweicloud_dataarts_studio_data_connection.test", "id"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.datasource_config.0.database", "huaweicloud_dli_database.test", "name"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.datasource_config.0.datatable", "huaweicloud_dli_table.test", "name"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.datasource_config.0.queue", "huaweicloud_dli_queue.test", "name"),
					resource.TestCheckResourceAttr(byId, "apis.0.datasource_config.0.response_params.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.datasource_config.0.response_params.0.name",
						"huaweicloud_dataarts_dataservice_api.test", "datasource_config.0.response_params.0.name"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.datasource_config.0.response_params.0.type",
						"huaweicloud_dataarts_dataservice_api.test", "datasource_config.0.response_params.0.type"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.datasource_config.0.response_params.0.field",
						"huaweicloud_dataarts_dataservice_api.test", "datasource_config.0.response_params.0.field"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.datasource_config.0.response_params.0.description",
						"huaweicloud_dataarts_dataservice_api.test", "datasource_config.0.response_params.0.description"),
					resource.TestCheckResourceAttr(byId, "apis.0.request_params.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.request_params.0.name",
						"huaweicloud_dataarts_dataservice_api.test", "request_params.0.name"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.request_params.0.position",
						"huaweicloud_dataarts_dataservice_api.test", "request_params.0.position"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.request_params.0.type",
						"huaweicloud_dataarts_dataservice_api.test", "request_params.0.type"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.request_params.0.description",
						"huaweicloud_dataarts_dataservice_api.test", "request_params.0.description"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.request_params.0.necessary",
						"huaweicloud_dataarts_dataservice_api.test", "request_params.0.necessary"),
					resource.TestCheckResourceAttrPair(byId, "apis.0.request_params.0.example_value",
						"huaweicloud_dataarts_dataservice_api.test", "request_params.0.example_value"),
					resource.TestCheckResourceAttrSet(byId, "apis.0.create_user"),
					resource.TestMatchResourceAttr(rName, "apis.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
					dcByDesc.CheckResourceExists(),
					resource.TestCheckOutput("is_desc_filter_useful", "true"),
					dcByCreator.CheckResourceExists(),
					resource.TestCheckOutput("is_creator_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByTableName.CheckResourceExists(),
					resource.TestCheckOutput("is_datatable_filter_useful", "true"),
					waitForDeletionCooldownComplete(),
				),
			},
		},
	})
}

func testAccDataSourceDataServiceApis_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_dataarts_studio_data_connection" "test" {
  workspace_id = "%[1]s"
  type         = "DLI"
  env_type     = "0"
  name         = "%[2]s"
  config       = jsonencode({
    "cdm_property_enable": "false"
  })
}

resource "huaweicloud_dli_database" "test" {
  name = "%[2]s"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[2]s_resource_vpc"
  data_location = "DLI"

  columns {
    name        = "resource_id"
    type        = "string"
    description = "The resource ID, in UUID format"
  }
}

resource "huaweicloud_vpc" "test" {
  name = "%[2]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = "%[2]s"
  min_cu                = 64
  max_cu                = 64
  cidr                  = cidrsubnet(huaweicloud_vpc.test.cidr, 3, 1)
  enterprise_project_id = "0"
}

resource "huaweicloud_dli_queue" "test" {
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name
  resource_mode              = 1

  # basic configuration
  name     = "%[2]s"
  cu_count = 16
}

// Under root path.
resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
}

resource "huaweicloud_dataarts_dataservice_api" "test" {
  workspace_id = "%[1]s"
  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = huaweicloud_dataarts_dataservice_catalog.test.id
  name         = "%[2]s"
  description  = "Created by terraform script"
  auth_type    = "NONE"
  manager      = "%[3]s"
  path         = "/terraform/auto/resource_query/{resource_type}"
  protocol     = "PROTOCOL_TYPE_HTTP"
  request_type = "REQUEST_TYPE_GET"
  visibility   = "PROJECT"

  request_params {
    name          = "resource_type"
    position      = "REQUEST_PARAMETER_POSITION_PATH"
    type          = "REQUEST_PARAMETER_TYPE_STRING"
    description   = "The type of the terraform resource to be queried"
    necessary     = true
    example_value = "demo"
  }

  datasource_config {
    type          = "DLI"
    connection_id = huaweicloud_dataarts_studio_data_connection.test.id
    database      = huaweicloud_dli_database.test.name
    datatable     = huaweicloud_dli_table.test.name
    queue         = huaweicloud_dli_queue.test.name

    response_params {
      name        = "resource_id"
      type        = "REQUEST_PARAMETER_TYPE_STRING"
      field       = "resource_id"
      description = "The resource ID, in UUID format"
    }
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_REVIEWER_NAME)
}

func testAccDataSourceDataServiceApis_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_dataservice_apis" "test" {
  depends_on = [
    huaweicloud_dataarts_dataservice_api.test
  ]

  workspace_id = "%[2]s"
}

# Filter by ID
locals {
  api_id = data.huaweicloud_dataarts_dataservice_apis.test.apis[0].id
}

data "huaweicloud_dataarts_dataservice_apis" "filter_by_id" {
  workspace_id = "%[2]s"
  api_id       = local.api_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_apis.filter_by_id.apis[*].id : v == local.api_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  name = data.huaweicloud_dataarts_dataservice_apis.test.apis[0].name
}

data "huaweicloud_dataarts_dataservice_apis" "filter_by_name" {
  workspace_id = "%[2]s"
  name         = local.name # Fuzzy search 
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_apis.filter_by_name.apis[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "huaweicloud_dataarts_dataservice_apis" "filter_by_name_not_found" {
  workspace_id = "%[2]s"
  name         = local.not_found_name # This name is not exist 
}

locals {
  name_not_found_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_apis.filter_by_name_not_found.apis[*].name : strcontains(v, local.not_found_name)
  ]
}

output "is_name_not_found_filter_useful" {
  value = length(local.name_not_found_filter_result) == 0
}

# Filter by description
locals {
  description = data.huaweicloud_dataarts_dataservice_apis.test.apis[0].description
}

data "huaweicloud_dataarts_dataservice_apis" "filter_by_desc" {
  workspace_id = "%[2]s"
  description  = local.description # Fuzzy search 
}

locals {
  desc_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_apis.filter_by_desc.apis[*].description : strcontains(v, local.description)
  ]
}

output "is_desc_filter_useful" {
  value = length(local.desc_filter_result) > 0 && alltrue(local.desc_filter_result)
}

# Filter by create user
locals {
  create_user = data.huaweicloud_dataarts_dataservice_apis.test.apis[0].create_user
}

data "huaweicloud_dataarts_dataservice_apis" "filter_by_creator" {
  workspace_id = "%[2]s"
  create_user  = local.create_user
}

locals {
  creator_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_apis.filter_by_desc.apis[*].create_user : v == local.create_user
  ]
}

output "is_creator_filter_useful" {
  value = length(local.creator_filter_result) > 0 && alltrue(local.creator_filter_result)
}

# Filter by type
locals {
  type = data.huaweicloud_dataarts_dataservice_apis.test.apis[0].type
}

data "huaweicloud_dataarts_dataservice_apis" "filter_by_type" {
  workspace_id = "%[2]s"
  type         = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_apis.filter_by_type.apis[*].type : v == local.type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by datatable
locals {
  datatable = data.huaweicloud_dataarts_dataservice_apis.test.apis[0].datasource_config[0].datatable
}

data "huaweicloud_dataarts_dataservice_apis" "filter_by_datatable" {
  workspace_id = "%[2]s"
  datatable    = local.datatable
}

locals {
  datatable_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_apis.filter_by_type.apis : v.datasource_config[0].datatable == local.datatable
  ]
}

output "is_datatable_filter_useful" {
  value = length(local.datatable_filter_result) > 0 && alltrue(local.datatable_filter_result)
}
`, testAccDataSourceDataServiceApis_base(), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
