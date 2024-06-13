package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceApi_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_apig_api.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApi_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "name", name),
					resource.TestCheckResourceAttr(dataSource, "type", "Public"),
					resource.TestCheckResourceAttr(dataSource, "request_method", "GET"),
					resource.TestCheckResourceAttr(dataSource, "request_path", "/user_info/{user_age}"),
					resource.TestCheckResourceAttr(dataSource, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(dataSource, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(dataSource, "simple_authentication", "false"),
					resource.TestCheckResourceAttrPair(dataSource, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttr(dataSource, "request_params.#", "2"),
					resource.TestCheckResourceAttr(dataSource, "request_params.0.%", "12"),
					resource.TestCheckResourceAttr(dataSource, "request_params.0.name", "X-TEST-ENUM"),
					resource.TestCheckResourceAttr(dataSource, "request_params.0.passthrough", "true"),
					resource.TestCheckResourceAttr(dataSource, "request_params.0.type", "STRING"),
					resource.TestCheckResourceAttr(dataSource, "request_params.0.minimum", "10"),
					resource.TestCheckResourceAttr(dataSource, "request_params.0.maximum", "20"),
					resource.TestCheckResourceAttr(dataSource, "request_params.1.name", "user_age"),
					resource.TestCheckResourceAttr(dataSource, "request_params.1.required", "true"),
					resource.TestCheckResourceAttr(dataSource, "request_params.1.minimum", "0"),
					resource.TestCheckResourceAttr(dataSource, "request_params.1.maximum", "200"),
					resource.TestCheckResourceAttr(dataSource, "request_params.1.location", "PATH"),
					resource.TestCheckResourceAttr(dataSource, "backend_params.#", "2"),
					resource.TestCheckResourceAttr(dataSource, "backend_params.0.%", "8"),
					resource.TestCheckResourceAttr(dataSource, "backend_params.0.name", "userAge"),
					resource.TestCheckResourceAttr(dataSource, "backend_params.1.name", "x-test-id"),
					resource.TestCheckResourceAttr(dataSource, "backend_params.1.system_param_type", "backend"),
					resource.TestCheckResourceAttr(dataSource, "description", "Created by script"),
					resource.TestCheckResourceAttr(dataSource, "matching", "Exact"),
					resource.TestCheckResourceAttr(dataSource, "success_response", "Success response"),
					resource.TestCheckResourceAttr(dataSource, "failure_response", "Failed response"),
					resource.TestCheckResourceAttr(dataSource, "web.0.%", "11"),
					resource.TestCheckResourceAttr(dataSource, "web.0.path", "/getUserAge/{userAge}"),
					resource.TestCheckResourceAttr(dataSource, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(dataSource, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(dataSource, "web.0.timeout", "30000"),
					resource.TestCheckResourceAttr(dataSource, "web.0.retry_count", "1"),
					resource.TestCheckResourceAttrPair(dataSource, "web.0.authorizer_id", "huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttrPair(dataSource, "web.0.vpc_channel_id", "huaweicloud_apig_vpc_channel.test", "id"),
					resource.TestCheckResourceAttr(dataSource, "web_policy.0.backend_params.#", "3"),
					resource.TestCheckResourceAttr(dataSource, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "web_policy.0.%", "14"),
					resource.TestCheckResourceAttr(dataSource, "web_policy.0.effective_mode", "ANY"),
					resource.TestCheckResourceAttr(dataSource, "web_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "web_policy.0.conditions.0.source", "param"),
					resource.TestCheckResourceAttr(dataSource, "web_policy.0.conditions.0.type", "Equal"),
					resource.TestCheckResourceAttr(dataSource, "web_policy.0.conditions.0.value", "28"),
					resource.TestCheckResourceAttrPair(dataSource, "web_policy.0.vpc_channel_id", "huaweicloud_apig_vpc_channel.test", "id"),
					resource.TestCheckResourceAttr(dataSource, "mock.#", "0"),
					resource.TestCheckResourceAttr(dataSource, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(dataSource, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.#", "0"),
					resource.TestMatchResourceAttr(dataSource, "registered_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceApi_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_api" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  api_id      = huaweicloud_apig_api.test.id
}
`, testAccApi_basic(testAccApi_base(name), name))
}

func TestAccDataSourceApi_functionGraph(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_apig_api.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApi_functionGraph(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "name", name),
					resource.TestCheckResourceAttr(dataSource, "type", "Public"),
					resource.TestCheckResourceAttr(dataSource, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(dataSource, "request_method", "POST"),
					resource.TestCheckResourceAttr(dataSource, "request_path", fmt.Sprintf("/test/function/%s", name)),
					resource.TestCheckResourceAttr(dataSource, "backend_type", "FUNCTION"),
					resource.TestCheckResourceAttr(dataSource, "func_graph.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "func_graph.0.invocation_type", "async"),
					resource.TestCheckResourceAttr(dataSource, "func_graph.0.timeout", "6000"),
					resource.TestCheckResourceAttrPair(dataSource, "func_graph.0.function_urn", "huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttrPair(dataSource, "func_graph.0.version", "huaweicloud_fgs_function.test", "version"),
					resource.TestCheckResourceAttrPair(dataSource, "func_graph.0.authorizer_id", "huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.invocation_type", "async"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.backend_params.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.backend_params.0.%", "8"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.backend_params.0.description", "created by terraform script"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.backend_params.0.location", "QUERY"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.backend_params.0.type", "CONSTANT"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.conditions.0.%", "10"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.conditions.0.source", "system"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.conditions.0.sys_name", "reqPath"),
					resource.TestCheckResourceAttr(dataSource, "func_graph_policy.0.conditions.0.type", "Equal"),
					resource.TestCheckResourceAttr(dataSource, "web.#", "0"),
					resource.TestCheckResourceAttr(dataSource, "web_policy.#", "0"),
					resource.TestCheckResourceAttr(dataSource, "mock.#", "0"),
					resource.TestCheckResourceAttr(dataSource, "mock_policy.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceApi_functionGraph(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_api" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  api_id      = huaweicloud_apig_api.test.id
}
`, testAccApi_functionGraph(name))
}

func TestAccDataSourceApi_mock(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_apig_api.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApi_mock(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "name", name),
					resource.TestCheckResourceAttr(dataSource, "type", "Private"),
					resource.TestCheckResourceAttr(dataSource, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(dataSource, "backend_type", "MOCK"),
					resource.TestCheckResourceAttr(dataSource, "matching", "Prefix"),
					resource.TestCheckResourceAttrPair(dataSource, "response_id", "huaweicloud_apig_response.test", "id"),
					resource.TestCheckResourceAttrPair(dataSource, "authorizer_id", "huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(dataSource, "body_description", "This is request body description"),
					resource.TestCheckResourceAttr(dataSource, "cors", "true"),
					resource.TestCheckResourceAttr(dataSource, "mock.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "mock.0.response", "Mock backend description"),
					resource.TestCheckResourceAttr(dataSource, "mock_policy.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "mock_policy.0.effective_mode", "ALL"),
					resource.TestCheckResourceAttr(dataSource, "mock_policy.0.conditions.0.sys_name", "reqMethod"),
					resource.TestCheckResourceAttrPair(dataSource, "publish_id", "huaweicloud_apig_api_publishment.test", "publish_id"),
					resource.TestMatchResourceAttr(dataSource, "published_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`))),
			},
		},
	})
}

func testAccDataSourceApi_mock(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"
  availability_zones    = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}

resource "huaweicloud_apig_group" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_apig_response" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  group_id    = huaweicloud_apig_group.test.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[2]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python3.6"
  code_type   = "inline"
}

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id      = huaweicloud_apig_instance.test.id
  name             = "%[2]s"
  function_urn     = huaweicloud_fgs_function.test.urn
  function_version = "latest"
  type             = "FRONTEND"
}

resource "huaweicloud_apig_api" "test" {
  instance_id             = huaweicloud_apig_instance.test.id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s"
  type                    = "Private"
  request_protocol        = "HTTP"
  request_method          = "POST"
  request_path            = "/test/mock"
  security_authentication = "AUTHORIZER"
  matching                = "Prefix"
  response_id             = huaweicloud_apig_response.test.id
  authorizer_id           = huaweicloud_apig_custom_authorizer.test.id
  body_description        = "This is request body description"
  cors                    = true

  mock {
    response = "Mock backend description"
  }

  mock_policy {
    name           = "%[2]s"
    response       = "Mock backend policy description"
    effective_mode = "ALL"

    conditions {
    source   = "system"
    type     = "Equal"
    value    = "GET"
    sys_name = "reqMethod"
    }
  }
}

resource "huaweicloud_apig_environment" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_apig_api_publishment" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_id      = huaweicloud_apig_api.test.id
}

data "huaweicloud_apig_api" "test" {
  depends_on = [
    huaweicloud_apig_api_publishment.test
  ]

  instance_id = huaweicloud_apig_instance.test.id
  api_id      = huaweicloud_apig_api.test.id
}
`, common.TestBaseNetwork(name), name)
}
