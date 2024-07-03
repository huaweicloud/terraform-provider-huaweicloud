package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceApi_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		webBackend       = "data.huaweicloud_apig_api.web"
		dcWithWebBackend = acceptance.InitDataSourceCheck(webBackend)
		fgsBackend       = "data.huaweicloud_apig_api.func_graph"
		dcWithFgsBackend = acceptance.InitDataSourceCheck(fgsBackend)
		mockBackend      = "data.huaweicloud_apig_api.mock"
		dcWithMock       = acceptance.InitDataSourceCheck(mockBackend)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running acceptance test for each kind of APIs, please make sure the agency already assign the FGS service.
			acceptance.TestAccPreCheckFgsAgency(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApi_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Query the API with Web backend
					dcWithWebBackend.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(webBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(webBackend, "api_id", "huaweicloud_apig_api.web", "id"),
					resource.TestCheckResourceAttrPair(webBackend, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttr(webBackend, "name", name+"_web"),
					resource.TestCheckResourceAttr(webBackend, "type", "Public"),
					resource.TestCheckResourceAttr(webBackend, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(webBackend, "request_method", "GET"),
					resource.TestCheckResourceAttr(webBackend, "request_path", "/web/test"),
					resource.TestCheckResourceAttr(webBackend, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(webBackend, "simple_authentication", "true"),
					resource.TestCheckResourceAttr(webBackend, "matching", "Exact"),
					resource.TestCheckResourceAttr(webBackend, "success_response", "Success response"),
					resource.TestCheckResourceAttr(webBackend, "failure_response", "Failed response"),
					resource.TestCheckResourceAttr(webBackend, "description", "Created by script"),
					resource.TestCheckResourceAttr(webBackend, "tags.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "tags.0", "foo"),
					resource.TestCheckResourceAttr(webBackend, "request_params.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.name", "X-Service-Num"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.type", "STRING"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.location", "HEADER"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.maximum", "20"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.minimum", "10"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.example", "TERRAFORM01"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.passthrough", "true"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.valid_enable", "1"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.0.type", "REQUEST"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.0.name", "SerivceNum"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.0.location", "HEADER"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.0.value", "X-Service-Num"),
					resource.TestCheckResourceAttr(webBackend, "web.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "web.0.path", "/web/test/backend"),
					resource.TestCheckResourceAttrPair(webBackend, "web.0.vpc_channel_id", "huaweicloud_apig_channel.test", "id"),
					resource.TestCheckResourceAttr(webBackend, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(webBackend, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(webBackend, "web.0.timeout", "30000"),
					resource.TestCheckResourceAttr(webBackend, "web.0.retry_count", "1"),
					resource.TestCheckResourceAttrPair(webBackend, "web.0.authorizer_id", "huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.name", name+"_web_policy"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.request_method", "GET"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.effective_mode", "ANY"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.path", "/web/test/backend"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.timeout", "30000"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.retry_count", "1"),
					resource.TestCheckResourceAttrPair(webBackend, "web_policy.0.vpc_channel_id", "huaweicloud_apig_channel.test", "id"),
					resource.TestCheckResourceAttrPair(webBackend, "web_policy.0.authorizer_id", "huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.type", "SYSTEM"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.name", "SerivceNum"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.location", "HEADER"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.value", "X-Service-Num"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.system_param_type", "backend"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.0.source", "param"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.0.param_name", "X-Service-Num"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.0.type", "Equal"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.0.value", "0001"),
					resource.TestCheckResourceAttr(webBackend, "mock.#", "0"),
					resource.TestCheckResourceAttr(webBackend, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(webBackend, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(webBackend, "func_graph_policy.#", "0"),
					resource.TestMatchResourceAttr(webBackend, "registered_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(webBackend, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Query the API with FunctionGraph backend
					dcWithFgsBackend.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(fgsBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(fgsBackend, "api_id", "huaweicloud_apig_api.func_graph", "id"),
					resource.TestCheckResourceAttrPair(fgsBackend, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttr(fgsBackend, "name", name+"_fgs"),
					resource.TestCheckResourceAttr(fgsBackend, "type", "Public"),
					resource.TestCheckResourceAttr(fgsBackend, "request_protocol", "GRPCS"),
					resource.TestCheckResourceAttr(fgsBackend, "request_method", "POST"),
					resource.TestCheckResourceAttr(fgsBackend, "request_path", "/fgs/test"),
					resource.TestCheckResourceAttr(fgsBackend, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(fgsBackend, "simple_authentication", "true"),
					resource.TestCheckResourceAttr(fgsBackend, "matching", "Exact"),
					resource.TestCheckResourceAttr(fgsBackend, "request_params.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "backend_params.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.#", "1"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph.0.function_urn", "huaweicloud_fgs_function.test.1", "urn"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph.0.version", "huaweicloud_fgs_function.test.1", "versions.0.name"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.0.network_type", "V2"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.0.request_protocol", "GRPCS"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.0.timeout", "5000"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.0.invocation_type", "sync"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph.0.authorizer_id", "huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.#", "1"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.name", name+"_fgs_policy"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph_policy.0.function_urn", "huaweicloud_fgs_function.test.1", "urn"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph_policy.0.version",
						"huaweicloud_fgs_function.test.1", "versions.0.name"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.network_type", "V2"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.request_protocol", "GRPCS"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.timeout", "5000"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.invocation_type", "sync"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.effective_mode", "ANY"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph_policy.0.authorizer_id",
						"huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.source", "cookie"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.cookie_name", "regex_test"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.type", "Matching"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.value", fmt.Sprintf("^%s:\\w+$", name)),
					resource.TestCheckResourceAttr(fgsBackend, "web.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "mock.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "web_policy.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "mock_policy.#", "0"),
					resource.TestMatchResourceAttr(fgsBackend, "registered_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(fgsBackend, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Query the API with mock configuration
					dcWithMock.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(mockBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(mockBackend, "api_id", "huaweicloud_apig_api.mock", "id"),
					resource.TestCheckResourceAttrPair(mockBackend, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttr(mockBackend, "name", name+"_mock"),
					resource.TestCheckResourceAttr(mockBackend, "type", "Public"),
					resource.TestCheckResourceAttr(mockBackend, "request_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(mockBackend, "request_method", "POST"),
					resource.TestCheckResourceAttr(mockBackend, "request_path", "/mock/test"),
					resource.TestCheckResourceAttr(mockBackend, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(mockBackend, "simple_authentication", "true"),
					resource.TestCheckResourceAttr(mockBackend, "matching", "Exact"),
					resource.TestCheckResourceAttr(mockBackend, "request_params.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "backend_params.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "mock.#", "1"),
					resource.TestCheckResourceAttr(mockBackend, "mock.0.status_code", "201"),
					resource.TestCheckResourceAttr(mockBackend, "mock.0.response", "{'message':'hello world'}"),
					resource.TestCheckResourceAttrPair(mockBackend, "mock.0.authorizer_id", "huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.#", "1"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.name", name+"_mock_policy"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.status_code", "201"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.response", "{'message':'hello world'}"),
					resource.TestCheckResourceAttrPair(mockBackend, "mock_policy.0.authorizer_id", "huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.effective_mode", "ANY"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.0.source", "system"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.0.type", "Equal"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.0.value", "user"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.0.sys_name", "reqPath"),
					resource.TestCheckResourceAttr(mockBackend, "web.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "web_policy.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "func_graph_policy.#", "0"),
					resource.TestMatchResourceAttr(mockBackend, "registered_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(mockBackend, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceApi_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api" "web" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s_web"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/web/test"
  security_authentication = "APP"
  simple_authentication   = true
  matching                = "Exact"
  success_response        = "Success response"
  failure_response        = "Failed response"
  description             = "Created by script"
  tags                    = ["foo"]

  request_params {
    name         = "X-Service-Num"
    type         = "STRING"
    location     = "HEADER"
    maximum      = 20
    minimum      = 10
    example      = "TERRAFORM01"
    passthrough  = true
    valid_enable = 1 # enable
  }

  backend_params {
    type     = "REQUEST"
    name     = "SerivceNum"
    location = "HEADER"
    value    = "X-Service-Num"
  }

  web {
    path             = "/web/test/backend"
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
    retry_count      = 1
    authorizer_id    = huaweicloud_apig_custom_authorizer.test.id
  }

  web_policy {
    name             = "%[2]s_web_policy"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    path             = "/web/test/backend"
    timeout          = 30000
    retry_count      = 1
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    authorizer_id    = huaweicloud_apig_custom_authorizer.test.id

    backend_params {
      type              = "SYSTEM"
      name              = "SerivceNum"
      location          = "HEADER"
      value             = "X-Service-Num"
      system_param_type = "backend"
    }

    conditions {
      source     = "param"
      param_name = "X-Service-Num"
      type       = "Equal"
      value      = "0001"
    }
  }
}

resource "huaweicloud_apig_api" "func_graph" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s_fgs"
  type                    = "Public"
  request_protocol        = "GRPCS"
  request_method          = "POST"
  request_path            = "/fgs/test"
  security_authentication = "APP"
  simple_authentication   = true
  matching                = "Exact"
  description             = "Created by script"
  tags                    = ["foo"]

  func_graph {
    function_urn     = huaweicloud_fgs_function.test[1].urn
    version          = tolist(huaweicloud_fgs_function.test[1].versions)[0].name
    network_type     = "V2"
    request_protocol = "GRPCS"
    timeout          = 5000
    invocation_type  = "sync"
    authorizer_id    = huaweicloud_apig_custom_authorizer.test.id
  }

  func_graph_policy {
    name             = "%[2]s_fgs_policy"
    function_urn     = huaweicloud_fgs_function.test[1].urn
    version          = tolist(huaweicloud_fgs_function.test[1].versions)[0].name
    network_type     = "V2"
    request_protocol = "GRPCS"
    timeout          = 5000
    invocation_type  = "sync"
    effective_mode   = "ANY"
    authorizer_id    = huaweicloud_apig_custom_authorizer.test.id

    conditions {
      source      = "cookie"
      cookie_name = "regex_test"
      type        = "Matching"
      value       = "^%[2]s:\\w+$"
    }
  }
}

resource "huaweicloud_apig_api" "mock" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s_mock"
  type                    = "Public"
  request_protocol        = "HTTPS"
  request_method          = "POST"
  request_path            = "/mock/test"
  security_authentication = "APP"
  simple_authentication   = true
  matching                = "Exact"
  success_response        = "Success response"
  failure_response        = "Failed response"
  description             = "Created by script"

  mock {
    status_code   = 201
    response      = "{'message':'hello world'}"
    authorizer_id = huaweicloud_apig_custom_authorizer.test.id
  }

  mock_policy {
    name           = "%[2]s_mock_policy"
    status_code    = 201
    response       = "{'message':'hello world'}"
    authorizer_id  = huaweicloud_apig_custom_authorizer.test.id
    effective_mode = "ANY"

    conditions {
      source   = "system"
      type     = "Equal"
      value    = "user"
      sys_name = "reqPath"
    }
  }
}

resource "huaweicloud_apig_environment" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
}

# Publish the API (with Web backend) and query it
resource "huaweicloud_apig_api_publishment" "web" {
  instance_id = local.instance_id
  env_id      = huaweicloud_apig_environment.test.id
  api_id      = huaweicloud_apig_api.web.id
}

data "huaweicloud_apig_api" "web" {
  depends_on = [
    huaweicloud_apig_api_publishment.web
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.web.id
}

# Publish the API (with FunctionGraph backend) and query it
resource "huaweicloud_apig_api_publishment" "func_graph" {
  instance_id = local.instance_id
  env_id      = huaweicloud_apig_environment.test.id
  api_id      = huaweicloud_apig_api.func_graph.id
}

data "huaweicloud_apig_api" "func_graph" {
  depends_on = [
    huaweicloud_apig_api_publishment.func_graph
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.func_graph.id
}

# Publish the API (with Mock) and query it
resource "huaweicloud_apig_api_publishment" "mock" {
  instance_id = local.instance_id
  env_id      = huaweicloud_apig_environment.test.id
  api_id      = huaweicloud_apig_api.mock.id
}

data "huaweicloud_apig_api" "mock" {
  depends_on = [
    huaweicloud_apig_api_publishment.mock
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.mock.id
}
`, testAccApi_base(name), name)
}
