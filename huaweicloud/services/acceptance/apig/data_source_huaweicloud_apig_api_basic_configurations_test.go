package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceApiBasicConfigurations_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_api_basic_configurations.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byId   = "data.huaweicloud_apig_api_basic_configurations.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_api_basic_configurations.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byExactName   = "data.huaweicloud_apig_api_basic_configurations.filter_by_exact_name"
		dcByExactName = acceptance.InitDataSourceCheck(byExactName)

		byNotFoundName   = "data.huaweicloud_apig_api_basic_configurations.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byGroupId   = "data.huaweicloud_apig_api_basic_configurations.filter_by_group_id"
		dcByGroupId = acceptance.InitDataSourceCheck(byGroupId)

		byType   = "data.huaweicloud_apig_api_basic_configurations.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byRequestMethod   = "data.huaweicloud_apig_api_basic_configurations.filter_by_request_method"
		dcByRequestMethod = acceptance.InitDataSourceCheck(byRequestMethod)

		byRequestPath   = "data.huaweicloud_apig_api_basic_configurations.filter_by_request_path"
		dcByRequestPath = acceptance.InitDataSourceCheck(byRequestPath)

		byRequestProtocol   = "data.huaweicloud_apig_api_basic_configurations.filter_by_request_protocol"
		dcByRequestProtocol = acceptance.InitDataSourceCheck(byRequestProtocol)

		bySecurityAuthentication   = "data.huaweicloud_apig_api_basic_configurations.filter_by_security_authentication"
		dcBySecurityAuthentication = acceptance.InitDataSourceCheck(bySecurityAuthentication)

		byVpcChannelName   = "data.huaweicloud_apig_api_basic_configurations.filter_by_vpc_channel_name"
		dcByVpcChannelName = acceptance.InitDataSourceCheck(byVpcChannelName)

		byEnvId   = "data.huaweicloud_apig_api_basic_configurations.filter_by_env_id"
		dcByEnvId = acceptance.InitDataSourceCheck(byEnvId)

		byEnvName   = "data.huaweicloud_apig_api_basic_configurations.filter_by_env_name"
		dcByEnvName = acceptance.InitDataSourceCheck(byEnvName)

		byBackEndType   = "data.huaweicloud_apig_api_basic_configurations.filter_by_backend_type"
		dcByBackEndType = acceptance.InitDataSourceCheck(byBackEndType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiBasicConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "configurations.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byId, "configurations.0.group_name"),
					resource.TestCheckResourceAttr(byId, "configurations.0.tags.#", "1"),
					resource.TestCheckResourceAttr(byId, "configurations.0.tags.0", "foo"),
					resource.TestCheckResourceAttrSet(byId, "configurations.0.group_version"),
					resource.TestCheckResourceAttrSet(byId, "configurations.0.publish_id"),
					resource.TestCheckResourceAttrSet(byId, "configurations.0.backend_type"),
					resource.TestCheckResourceAttr(byId, "configurations.0.simple_authentication", "false"),
					resource.TestCheckResourceAttr(byId, "configurations.0.cors", "true"),
					resource.TestCheckResourceAttr(byId, "configurations.0.matching", "Exact"),
					resource.TestCheckResourceAttr(byId, "configurations.0.description", "Created by script"),
					resource.TestMatchResourceAttr(byId, "configurations.0.registered_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byId, "configurations.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byId, "configurations.0.published_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByExactName.CheckResourceExists(),
					resource.TestCheckOutput("is_exact_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
					dcByGroupId.CheckResourceExists(),
					resource.TestCheckOutput("is_group_id_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByRequestMethod.CheckResourceExists(),
					resource.TestCheckOutput("is_request_method_filter_useful", "true"),
					dcByRequestPath.CheckResourceExists(),
					resource.TestCheckOutput("is_request_path_filter_useful", "true"),
					dcByRequestProtocol.CheckResourceExists(),
					resource.TestCheckOutput("is_request_protocol_filter_useful", "true"),
					dcBySecurityAuthentication.CheckResourceExists(),
					resource.TestCheckOutput("is_security_authentication_filter_useful", "true"),
					dcByVpcChannelName.CheckResourceExists(),
					resource.TestCheckOutput("is_vpc_channel_name_filter_useful", "true"),
					dcByEnvId.CheckResourceExists(),
					resource.TestCheckOutput("is_env_id_filter_useful", "true"),
					dcByEnvName.CheckResourceExists(),
					resource.TestCheckOutput("is_env_name_filter_useful", "true"),
					dcByBackEndType.CheckResourceExists(),
					resource.TestCheckOutput("is_backend_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceApiBasicConfigurations_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_custom_authorizer" "frontEnd" {
  instance_id      = local.instance_id
  name             = "%[2]s_front"
  function_urn     = huaweicloud_fgs_function.test[1].urn
  function_version = "latest"
  type             = "FRONTEND"
}

resource "huaweicloud_apig_api" "test" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/{%[2]s}"
  security_authentication = "AUTHORIZER"
  matching                = "Exact"
  authorizer_id           = huaweicloud_apig_custom_authorizer.frontEnd.id
  cors                    = true
  description             = "Created by script"
  tags                    = ["foo"]
  
  request_params {
    name     = "%[2]s"
    type     = "NUMBER"
    location = "PATH"
    required = true
    maximum  = 200
    minimum  = 0
  }

  web {
    path             = "/"
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
    retry_count      = 1
  }
}

resource "huaweicloud_apig_environment" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
}

resource "huaweicloud_apig_api_publishment" "test" {
  instance_id = local.instance_id
  env_id      = huaweicloud_apig_environment.test.id
  api_id      = huaweicloud_apig_api.test.id
}

`, testAccApi_base(name), name)
}

func testAccDataSourceApiBasicConfigurations_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_api_basic_configurations" "test" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id = local.instance_id
}

# Filter by ID
locals {
  api_id = huaweicloud_apig_api.test.id
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_id" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id = local.instance_id
  api_id      = local.api_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_id.configurations[*].id : v == local.api_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name (fuzzy search)
locals {
  api_name = huaweicloud_apig_api.test.name
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_name" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id = local.instance_id
  name        = local.api_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_name.configurations[*].name : strcontains(v, local.api_name)
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (exact search)
data "huaweicloud_apig_api_basic_configurations" "filter_by_exact_name" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id    = local.instance_id
  name           = local.api_name
  precise_search = "name,req_uri"
}

output "is_exact_name_filter_useful" {
  value = length(data.huaweicloud_apig_api_basic_configurations.filter_by_exact_name.configurations) == 1
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_not_found_name" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id    = local.instance_id
  name           = local.not_found_name
  precise_search = "name"
}

output "is_name_not_found_filter_useful" {
  value = length(data.huaweicloud_apig_api_basic_configurations.filter_by_not_found_name.configurations) == 0
}

# Filter by group ID
locals {
  group_id = huaweicloud_apig_api.test.group_id
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_group_id" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id = local.instance_id
  group_id    = local.group_id
}

locals {
  group_id_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_id.configurations[*].group_id : v == local.group_id
  ]
}

output "is_group_id_filter_useful" {
  value = length(local.group_id_filter_result) > 0 && alltrue(local.group_id_filter_result)
}

# Filter by type
locals {
  api_type = huaweicloud_apig_api.test.type
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_type" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id = local.instance_id
  type        = local.api_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_type.configurations[*].type : v == local.api_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by request method
locals {
  request_method = huaweicloud_apig_api.test.request_method
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_request_method" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id    = local.instance_id
  request_method = local.request_method
}

locals {
  request_method_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_request_method.configurations[*].request_method : v == local.request_method
  ]
}

output "is_request_method_filter_useful" {
  value = length(local.request_method_filter_result) > 0 && alltrue(local.request_method_filter_result)
}

# Filter by request path
locals {
  request_path = huaweicloud_apig_api.test.request_path
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_request_path" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id    = local.instance_id
  request_path   = local.request_path
  precise_search = "req_url"
}

locals {
  request_path_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_request_path.configurations[*].request_path : v == local.request_path
  ]
}

output "is_request_path_filter_useful" {
  value = length(local.request_path_filter_result) > 0 && alltrue(local.request_path_filter_result)
}

# Filter by request protocol
locals {
  request_protocol = huaweicloud_apig_api.test.request_protocol
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_request_protocol" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id      = local.instance_id
  request_protocol = local.request_protocol
}

locals {
  request_protocol_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_request_protocol.configurations[*].request_protocol : 
      v == local.request_protocol
  ]
}

output "is_request_protocol_filter_useful" {
  value = length(local.request_protocol_filter_result) > 0 && alltrue(local.request_protocol_filter_result)
}

# Filter by security authentication type
locals {
  security_authentication = huaweicloud_apig_api.test.security_authentication
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_security_authentication" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id             = local.instance_id
  security_authentication = local.security_authentication
}

locals {
  security_auth_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_security_authentication.configurations[*].security_authentication : 
      v == local.security_authentication
  ]
}

output "is_security_authentication_filter_useful" {
  value = length(local.security_auth_filter_result) > 0 && alltrue(local.security_auth_filter_result)
}

# Filter by vpc channel name
data "huaweicloud_apig_api_basic_configurations" "filter_by_vpc_channel_name" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id      = local.instance_id
  vpc_channel_name = huaweicloud_apig_channel.test.name
}

output "is_vpc_channel_name_filter_useful" {
  value = length(data.huaweicloud_apig_api_basic_configurations.filter_by_vpc_channel_name.configurations) > 0
}

# Filter by env ID
locals {
  env_id = huaweicloud_apig_environment.test.id
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_env_id" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id = local.instance_id
  env_id      = local.env_id
}

locals {
  env_id_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_env_id.configurations[*].env_id : v == local.env_id
  ]
}

output "is_env_id_filter_useful" {
  value = length(local.env_id_filter_result) > 0 && alltrue(local.env_id_filter_result)
}

# Filter by env name
locals {
  env_name = huaweicloud_apig_environment.test.name
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_env_name" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id = local.instance_id
  env_name    = local.env_name
}

locals {
  env_name_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_env_name.configurations[*].env_name : v == local.env_name
  ]
}

output "is_env_name_filter_useful" {
  value = length(local.env_name_filter_result) > 0 && alltrue(local.env_name_filter_result)
}

# Filter by backend type
# There is no "backend_type" field in the parent resource.
locals {
  backend_type = data.huaweicloud_apig_api_basic_configurations.filter_by_id.configurations[0].backend_type
}

data "huaweicloud_apig_api_basic_configurations" "filter_by_backend_type" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id  = local.instance_id
  backend_type = local.backend_type
}

locals {
  backend_type_filter_result = [
    for v in data.huaweicloud_apig_api_basic_configurations.filter_by_backend_type.configurations[*].backend_type : v == local.backend_type
  ]
}

output "is_backend_type_filter_useful" {
  value = length(local.backend_type_filter_result) > 0 && alltrue(local.backend_type_filter_result)
}
`, testAccDataSourceApiBasicConfigurations_base())
}
