package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getPublishmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return apig.GetVersionHistories(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["env_id"],
		state.Primary.Attributes["api_id"])
}

func TestAccApiPublishment_basic(t *testing.T) {
	var (
		histories []apis.ApiVersionInfo

		rName = acceptance.RandomAccResourceName()

		webBackend       = "huaweicloud_apig_api_publishment.web"
		rcWithWebBackend = acceptance.InitResourceCheck(webBackend, &histories, getPublishmentResourceFunc)
		fgsBackend       = "huaweicloud_apig_api_publishment.func_graph"
		rcWithFgsBackend = acceptance.InitResourceCheck(fgsBackend, &histories, getPublishmentResourceFunc)
		mockBackend      = "huaweicloud_apig_api_publishment.mock"
		rcWithMock       = acceptance.InitResourceCheck(mockBackend, &histories, getPublishmentResourceFunc)
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
		CheckDestroy:      rcWithWebBackend.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApiPublishment_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					// Publish the API with Web backend
					rcWithWebBackend.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(webBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(webBackend, "env_id", "huaweicloud_apig_environment.test", "id"),
					resource.TestCheckResourceAttrPair(webBackend, "env_name", "huaweicloud_apig_environment.test", "name"),
					resource.TestCheckResourceAttrPair(webBackend, "api_id", "huaweicloud_apig_api.web", "id"),
					resource.TestMatchResourceAttr(webBackend, "published_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(webBackend, "publish_id"),
					// Publish the API with FunctionGraph backend
					rcWithFgsBackend.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(fgsBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(fgsBackend, "env_id", "huaweicloud_apig_environment.test", "id"),
					resource.TestCheckResourceAttrPair(fgsBackend, "env_name", "huaweicloud_apig_environment.test", "name"),
					resource.TestCheckResourceAttrPair(fgsBackend, "api_id", "huaweicloud_apig_api.func_graph", "id"),
					resource.TestMatchResourceAttr(fgsBackend, "published_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(fgsBackend, "publish_id"),
					// Publish the API with Mock configuration
					rcWithMock.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(mockBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(mockBackend, "env_id", "huaweicloud_apig_environment.test", "id"),
					resource.TestCheckResourceAttrPair(mockBackend, "env_name", "huaweicloud_apig_environment.test", "name"),
					resource.TestCheckResourceAttrPair(mockBackend, "api_id", "huaweicloud_apig_api.mock", "id"),
					resource.TestMatchResourceAttr(mockBackend, "published_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(mockBackend, "publish_id"),
				),
			},
			{
				ResourceName:      webBackend,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      fgsBackend,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      mockBackend,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApiPublishment_basic(name string) string {
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
