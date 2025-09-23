package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getApiFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return apis.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccApi_basic(t *testing.T) {
	var (
		api apis.APIResp

		webBackend       = "huaweicloud_apig_api.web"
		rcWithWebBackend = acceptance.InitResourceCheck(webBackend, &api, getApiFunc)
		fgsBackend       = "huaweicloud_apig_api.func_graph"
		rcWithFgsBackend = acceptance.InitResourceCheck(fgsBackend, &api, getApiFunc)
		mockBackend      = "huaweicloud_apig_api.mock"
		rcWithMock       = acceptance.InitResourceCheck(mockBackend, &api, getApiFunc)

		name        = acceptance.RandomAccResourceName()
		updateName  = acceptance.RandomAccResourceName()
		basicConfig = testAccApi_base(name)
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
				Config: testAccApi_basic_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					// Web backend
					rcWithWebBackend.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(webBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(webBackend, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttr(webBackend, "name", name+"_web"),
					resource.TestCheckResourceAttr(webBackend, "type", "Public"),
					resource.TestCheckResourceAttr(webBackend, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(webBackend, "request_method", "GET"),
					resource.TestCheckResourceAttr(webBackend, "request_path", "/web/test"),
					resource.TestCheckResourceAttr(webBackend, "security_authentication", "AUTHORIZER"),
					resource.TestCheckResourceAttrPair(webBackend, "authorizer_id", "huaweicloud_apig_custom_authorizer.front", "id"),
					resource.TestCheckResourceAttr(webBackend, "simple_authentication", "false"),
					resource.TestCheckResourceAttr(webBackend, "content_type", ""),
					resource.TestCheckResourceAttr(webBackend, "is_send_fg_body_base64", "true"),
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
					resource.TestCheckResourceAttrPair(webBackend, "web.0.authorizer_id", "huaweicloud_apig_custom_authorizer.backend", "id"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.name", name+"_web_policy"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.request_method", "GET"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.effective_mode", "ANY"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.path", "/web/test/backend"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.timeout", "30000"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.retry_count", "1"),
					resource.TestCheckResourceAttrPair(webBackend, "web_policy.0.vpc_channel_id", "huaweicloud_apig_channel.test", "id"),
					resource.TestCheckResourceAttrPair(webBackend, "web_policy.0.authorizer_id", "huaweicloud_apig_custom_authorizer.backend", "id"),
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
					// FunctionGraph backend
					rcWithFgsBackend.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(fgsBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
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
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph.0.authorizer_id", "huaweicloud_apig_custom_authorizer.backend", "id"),
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
						"huaweicloud_apig_custom_authorizer.backend", "id"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.source", "cookie"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.cookie_name", "regex_test"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.type", "Matching"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.value", fmt.Sprintf("^%s:\\w+$", name)),
					resource.TestCheckResourceAttr(fgsBackend, "web.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "mock.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "web_policy.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "mock_policy.#", "0"),
					// Mock configuration
					rcWithMock.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(mockBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
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
					resource.TestCheckResourceAttrPair(mockBackend, "mock.0.authorizer_id", "huaweicloud_apig_custom_authorizer.backend", "id"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.#", "1"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.name", name+"_mock_policy"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.status_code", "201"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.response", "{'message':'hello world'}"),
					resource.TestCheckResourceAttrPair(mockBackend, "mock_policy.0.authorizer_id",
						"huaweicloud_apig_custom_authorizer.backend", "id"),
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
				),
			},
			{
				Config: testAccApi_basic_step2(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					// Web backend
					rcWithWebBackend.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(webBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(webBackend, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttr(webBackend, "name", updateName+"_web"),
					resource.TestCheckResourceAttr(webBackend, "type", "Private"),
					resource.TestCheckResourceAttr(webBackend, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(webBackend, "request_method", "GET"),
					resource.TestCheckResourceAttr(webBackend, "request_path", "/web/new/test"),
					resource.TestCheckResourceAttr(webBackend, "security_authentication", "AUTHORIZER"),
					resource.TestCheckResourceAttrPair(webBackend, "authorizer_id", "huaweicloud_apig_custom_authorizer.front", "id"),
					resource.TestCheckResourceAttr(webBackend, "simple_authentication", "false"),
					resource.TestCheckResourceAttr(webBackend, "content_type", "application/json"),
					resource.TestCheckResourceAttr(webBackend, "is_send_fg_body_base64", "false"),
					resource.TestCheckResourceAttr(webBackend, "matching", "Exact"),
					resource.TestCheckResourceAttr(webBackend, "success_response", "Updated success response"),
					resource.TestCheckResourceAttr(webBackend, "failure_response", "Updated failed response"),
					resource.TestCheckResourceAttr(webBackend, "description", ""),
					resource.TestCheckResourceAttr(webBackend, "tags.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "tags.0", "key"),
					resource.TestCheckResourceAttr(webBackend, "request_params.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.name", "X-Service-Name"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.type", "STRING"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.location", "HEADER"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.maximum", "30"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.minimum", "5"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.example", "XF"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.passthrough", "false"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.enumeration", "AF,RF"),
					resource.TestCheckResourceAttr(webBackend, "request_params.0.valid_enable", "2"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.0.type", "REQUEST"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.0.name", "ServiceName"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.0.location", "HEADER"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.0.value", "X-Service-Name"),
					resource.TestCheckResourceAttr(webBackend, "backend_params.0.system_param_type", "backend"),
					resource.TestCheckResourceAttr(webBackend, "web.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "web.0.path", "/web/new/test/backend"),
					resource.TestCheckResourceAttrPair(webBackend, "web.0.vpc_channel_id", "huaweicloud_apig_channel.test", "id"),
					resource.TestCheckResourceAttr(webBackend, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(webBackend, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(webBackend, "web.0.timeout", "40000"),
					resource.TestCheckResourceAttr(webBackend, "web.0.retry_count", "2"),
					resource.TestCheckResourceAttrPair(webBackend, "web.0.authorizer_id", "huaweicloud_apig_custom_authorizer.backend", "id"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.name", updateName+"_web_policy"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.request_method", "GET"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.effective_mode", "ALL"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.path", "/web/new/test/backend"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.timeout", "40000"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.retry_count", "2"),
					resource.TestCheckResourceAttrPair(webBackend, "web_policy.0.vpc_channel_id", "huaweicloud_apig_channel.test", "id"),
					resource.TestCheckResourceAttrPair(webBackend, "web_policy.0.authorizer_id", "huaweicloud_apig_custom_authorizer.backend", "id"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.type", "SYSTEM"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.name", "ServiceName"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.location", "HEADER"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.value", "X-Service-Name"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.backend_params.0.system_param_type", "backend"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.0.source", "param"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.0.param_name", "X-Service-Name"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.0.type", "Equal"),
					resource.TestCheckResourceAttr(webBackend, "web_policy.0.conditions.0.value", "TF"),
					resource.TestCheckResourceAttr(webBackend, "mock.#", "0"),
					resource.TestCheckResourceAttr(webBackend, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(webBackend, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(webBackend, "func_graph_policy.#", "0"),
					// FunctionGraph backend
					rcWithFgsBackend.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(fgsBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(fgsBackend, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttr(fgsBackend, "name", updateName+"_fgs"),
					resource.TestCheckResourceAttr(fgsBackend, "type", "Private"),
					resource.TestCheckResourceAttr(fgsBackend, "request_protocol", "GRPCS"),
					resource.TestCheckResourceAttr(fgsBackend, "request_method", "POST"),
					resource.TestCheckResourceAttr(fgsBackend, "request_path", "/fgs/new/test"),
					resource.TestCheckResourceAttr(fgsBackend, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(fgsBackend, "simple_authentication", "false"),
					resource.TestCheckResourceAttr(fgsBackend, "matching", "Exact"),
					resource.TestCheckResourceAttr(fgsBackend, "description", ""),
					resource.TestCheckResourceAttr(fgsBackend, "request_params.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "backend_params.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.#", "1"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph.0.function_urn", "huaweicloud_fgs_function.test.1", "urn"),
					resource.TestCheckResourceAttrSet(fgsBackend, "func_graph.0.function_alias_urn"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.0.network_type", "V2"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.0.request_protocol", "GRPCS"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.0.timeout", "6000"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph.0.invocation_type", "sync"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph.0.authorizer_id", "huaweicloud_apig_custom_authorizer.backend", "id"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.#", "1"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.name", updateName+"_fgs_policy"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph_policy.0.function_urn", "huaweicloud_fgs_function.test.1", "urn"),
					resource.TestCheckResourceAttrSet(fgsBackend, "func_graph_policy.0.function_alias_urn"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.network_type", "V2"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.request_protocol", "GRPCS"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.timeout", "6000"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.invocation_type", "sync"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.effective_mode", "ALL"),
					resource.TestCheckResourceAttrPair(fgsBackend, "func_graph_policy.0.authorizer_id",
						"huaweicloud_apig_custom_authorizer.backend", "id"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.source", "cookie"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.cookie_name", "regex_test"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.type", "Matching"),
					resource.TestCheckResourceAttr(fgsBackend, "func_graph_policy.0.conditions.0.value", fmt.Sprintf("^cookie-%s:\\w+$", updateName)),
					resource.TestCheckResourceAttr(fgsBackend, "web.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "mock.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "web_policy.#", "0"),
					resource.TestCheckResourceAttr(fgsBackend, "mock_policy.#", "0"),
					// Mock configuration
					rcWithMock.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(mockBackend, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(mockBackend, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttr(mockBackend, "name", updateName+"_mock"),
					resource.TestCheckResourceAttr(mockBackend, "type", "Private"),
					resource.TestCheckResourceAttr(mockBackend, "request_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(mockBackend, "request_method", "POST"),
					resource.TestCheckResourceAttr(mockBackend, "request_path", "/mock/new/test"),
					resource.TestCheckResourceAttr(mockBackend, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(mockBackend, "simple_authentication", "false"),
					resource.TestCheckResourceAttr(mockBackend, "matching", "Exact"),
					resource.TestCheckResourceAttr(webBackend, "success_response", "Updated success response"),
					resource.TestCheckResourceAttr(webBackend, "failure_response", "Updated failed response"),
					resource.TestCheckResourceAttr(webBackend, "description", ""),
					resource.TestCheckResourceAttr(mockBackend, "request_params.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "backend_params.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "mock.#", "1"),
					resource.TestCheckResourceAttr(mockBackend, "mock.0.status_code", "202"),
					resource.TestCheckResourceAttr(mockBackend, "mock.0.response", "{'message':'hello world!'}"),
					resource.TestCheckResourceAttrPair(mockBackend, "mock.0.authorizer_id", "huaweicloud_apig_custom_authorizer.backend", "id"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.#", "1"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.name", updateName+"_mock_policy"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.status_code", "202"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.response", "{'message':'hello world!'}"),
					resource.TestCheckResourceAttrPair(mockBackend, "mock_policy.0.authorizer_id",
						"huaweicloud_apig_custom_authorizer.backend", "id"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.effective_mode", "ALL"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.0.source", "system"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.0.type", "Equal"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.0.value", "GET"),
					resource.TestCheckResourceAttr(mockBackend, "mock_policy.0.conditions.0.sys_name", "reqMethod"),
					resource.TestCheckResourceAttr(mockBackend, "web.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "web_policy.#", "0"),
					resource.TestCheckResourceAttr(mockBackend, "func_graph_policy.#", "0"),
				),
			},
			{
				ResourceName:      webBackend,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApiResourceImportStateFunc(webBackend),
			},
			{
				ResourceName:      fgsBackend,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApiResourceImportStateFunc(fgsBackend),
			},
			{
				ResourceName:      mockBackend,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApiResourceImportStateFunc(mockBackend),
			},
		},
	})
}

func testAccApiResourceImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{instance_id}/{name}', but '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccApi_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  count = 2

  functiongraph_version = "v2"
  agency                = "%[3]s"
  name                  = format("%[2]s_http_%%d", count.index)
  app                   = "default"
  handler               = "bootstrap"
  memory_size           = 128
  timeout               = 3
  runtime               = "http"
  code_type             = "inline"
  func_code             = "dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganN="

  versions {
    name = "%[2]s"

    aliases {
      name = "custom_alias"
    }
  }
}

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[4]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_group" "test" {
  name        = "%[2]s"
  instance_id = local.instance_id
}

// When the backend uses the load balance channel, the number of retry count cannot exceed the number of available
// backend servers in the load balance channel.
resource "huaweicloud_compute_instance" "test" {
  count = 3

  name               = format("%[2]s_%%d", count.index)
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = "%[5]s"
  }
}

resource "huaweicloud_apig_channel" "test" {
  instance_id        = local.instance_id
  name               = "%[2]s"
  port               = 80
  balance_strategy   = 1
  member_type        = "ecs"
  type               = 2

  health_check {
    protocol           = "TCP"
    threshold_normal   = 1 # minimum value
    threshold_abnormal = 1 # minimum value
    interval           = 1 # minimum value
    timeout            = 1 # minimum value
  }

  dynamic "member" {
    for_each = huaweicloud_compute_instance.test[*]

    content {
      id   = member.value.id
      name = member.value.name
    }
  }
}

resource "huaweicloud_apig_custom_authorizer" "front" {
  instance_id        = local.instance_id
  name               = "%[2]s_front"
  function_urn       = huaweicloud_fgs_function.test[0].urn
  function_alias_uri = format("%%s:!%%s", huaweicloud_fgs_function.test[0].urn,
    tolist(huaweicloud_fgs_function.test[0].versions)[0].aliases[0].name)
  network_type       = "V2"
  type               = "FRONTEND"
  is_body_send       = true
  user_data          = "Demo"
  cache_age          = 15

  identity {
    name     = "X-Service-Num"
    location = "HEADER"
  }
}

resource "huaweicloud_apig_custom_authorizer" "backend" {
  instance_id        = local.instance_id
  name               = "%[2]s_backend"
  function_urn       = huaweicloud_fgs_function.test[0].urn
  function_alias_uri = format("%%s:!%%s", huaweicloud_fgs_function.test[0].urn,
    tolist(huaweicloud_fgs_function.test[0].versions)[0].aliases[0].name)
  network_type       = "V2"
  type               = "BACKEND"
  is_body_send       = true
  user_data          = "Demo"
  cache_age          = 15
}
`, common.TestBaseComputeResources(name), name, acceptance.HW_FGS_AGENCY_NAME,
		acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}

func testAccApi_basic_step1(baseConfig, name string) string {
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
  security_authentication = "AUTHORIZER"
  authorizer_id           = huaweicloud_apig_custom_authorizer.front.id
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
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id
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
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id

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
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id
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
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id

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
    authorizer_id = huaweicloud_apig_custom_authorizer.backend.id
  }

  mock_policy {
    name           = "%[2]s_mock_policy"
    status_code    = 201
    response       = "{'message':'hello world'}"
    authorizer_id  = huaweicloud_apig_custom_authorizer.backend.id
    effective_mode = "ANY"

    conditions {
      source   = "system"
      type     = "Equal"
      value    = "user"
      sys_name = "reqPath"
    }
  }
}
`, baseConfig, name)
}

func testAccApi_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api" "web" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s_web"
  type                    = "Private"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/web/new/test"
  security_authentication = "AUTHORIZER"
  authorizer_id           = huaweicloud_apig_custom_authorizer.front.id
  matching                = "Exact"
  success_response        = "Updated success response"
  failure_response        = "Updated failed response"
  tags                    = ["key"]
  content_type            = "application/json"
  is_send_fg_body_base64  = false

  request_params {
    name         = "X-Service-Name"
    type         = "STRING"
    location     = "HEADER"
    maximum      = 30
    minimum      = 5
    example      = "XF"
    passthrough  = false
    enumeration  = "AF,RF"
    valid_enable = 2 # disable
  }

  backend_params {
    type              = "REQUEST"
    name              = "ServiceName"
    location          = "HEADER"
    value             = "X-Service-Name"
    system_param_type = "backend"
  }

  web {
    path             = "/web/new/test/backend"
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 40000
    retry_count      = 2
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id
  }

  web_policy {
    name             = "%[2]s_web_policy"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ALL"
    path             = "/web/new/test/backend"
    timeout          = 40000
    retry_count      = 2
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id

    backend_params {
      type              = "SYSTEM"
      name              = "ServiceName"
      location          = "HEADER"
      value             = "X-Service-Name"
      system_param_type = "backend"
    }

    conditions {
      source     = "param"
      param_name = "X-Service-Name"
      type       = "Equal"
      value      = "TF"
    }
  }
}

resource "huaweicloud_apig_api" "func_graph" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s_fgs"
  type                    = "Private"
  request_protocol        = "GRPCS"
  request_method          = "POST"
  request_path            = "/fgs/new/test"
  security_authentication = "APP"
  simple_authentication   = false
  matching                = "Exact"

  func_graph {
    function_urn       = huaweicloud_fgs_function.test[1].urn
    function_alias_urn = format("%%s:!%%s", huaweicloud_fgs_function.test[1].urn,
	  tolist(huaweicloud_fgs_function.test[1].versions)[0].aliases[0].name)
    network_type       = "V2"
    request_protocol   = "GRPCS"
    timeout            = 6000
    invocation_type    = "sync"
    authorizer_id      = huaweicloud_apig_custom_authorizer.backend.id
  }

  func_graph_policy {
    name               = "%[2]s_fgs_policy"
    function_urn       = huaweicloud_fgs_function.test[1].urn
    function_alias_urn = format("%%s:!%%s", huaweicloud_fgs_function.test[1].urn,
	  tolist(huaweicloud_fgs_function.test[1].versions)[0].aliases[0].name)
    network_type       = "V2"
    request_protocol   = "GRPCS"
    timeout            = 6000
    invocation_type    = "sync"
    effective_mode     = "ALL"
    authorizer_id      = huaweicloud_apig_custom_authorizer.backend.id

    conditions {
      source      = "cookie"
      cookie_name = "regex_test"
      type        = "Matching"
      value       = "^cookie-%[2]s:\\w+$"
    }
  }
}

resource "huaweicloud_apig_api" "mock" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s_mock"
  type                    = "Private"
  request_protocol        = "HTTPS"
  request_method          = "POST"
  request_path            = "/mock/new/test"
  security_authentication = "APP"
  simple_authentication   = false
  matching                = "Exact"
  success_response        = "Updated success response"
  failure_response        = "Updated failed response"

  mock {
    status_code   = 202
    response      = "{'message':'hello world!'}"
    authorizer_id = huaweicloud_apig_custom_authorizer.backend.id
  }

  mock_policy {
    name           = "%[2]s_mock_policy"
    status_code    = 202
    response       = "{'message':'hello world!'}"
    authorizer_id  = huaweicloud_apig_custom_authorizer.backend.id
    effective_mode = "ALL"

    conditions {
      source   = "system"
      type     = "Equal"
      value    = "GET"
      sys_name = "reqMethod"
    }
  }
}
`, baseConfig, name)
}

func TestAccApi_orchestration(t *testing.T) {
	var (
		api apis.APIResp

		resourceName = "huaweicloud_apig_api.test"
		rc           = acceptance.InitResourceCheck(resourceName, &api, getApiFunc)

		name        = acceptance.RandomAccResourceName()
		basicConfig = testAccApi_orchestration_base(name)
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
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApi_orchestration_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					// Web backend
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "Private"),
					resource.TestCheckResourceAttr(resourceName, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "request_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "request_path", "/orchestration/test"),
					resource.TestCheckResourceAttr(resourceName, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(resourceName, "request_params.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "request_params.0.name", "X-Service-Name"),
					resource.TestCheckResourceAttr(resourceName, "request_params.0.type", "STRING"),
					resource.TestCheckResourceAttr(resourceName, "request_params.0.location", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "request_params.0.maximum", "30"),
					resource.TestCheckResourceAttr(resourceName, "request_params.0.minimum", "5"),
					resource.TestCheckResourceAttr(resourceName, "request_params.0.orchestrations.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "request_params.0.orchestrations.0",
						"huaweicloud_apig_orchestration_rule.type_none_value", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "request_params.0.orchestrations.1",
						"huaweicloud_apig_orchestration_rule.type_list", "id"),
					resource.TestCheckResourceAttr(resourceName, "backend_params.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "backend_params.0.type", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "backend_params.0.name", "ServiceName"),
					resource.TestCheckResourceAttr(resourceName, "backend_params.0.location", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "backend_params.0.value", "X-Service-Name"),
					resource.TestCheckResourceAttr(resourceName, "backend_params.0.system_param_type", "backend"),
					resource.TestCheckResourceAttr(resourceName, "web.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web.0.path", "/orchestration/test/backend"),
					resource.TestCheckResourceAttrPair(resourceName, "web.0.vpc_channel_id", "huaweicloud_apig_channel.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "web.0.timeout", "40000"),
					resource.TestCheckResourceAttr(resourceName, "web.0.retry_count", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "web.0.authorizer_id", "huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.name", name+"_orchestration_policy"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.request_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.effective_mode", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.path", "/orchestration/test/backend/list"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.timeout", "40000"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.retry_count", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "web_policy.0.vpc_channel_id", "huaweicloud_apig_channel.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.backend_params.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.backend_params.0.type", "SYSTEM"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.backend_params.0.name", "ServiceName"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.backend_params.0.location", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.backend_params.0.value", "X-Service-Name"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.backend_params.0.system_param_type", "backend"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.0.source", "orchestration"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.0.mapped_param_name", "ServiceName"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.0.mapped_param_location", "header"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.0.value", "ValueAA"),
					resource.TestCheckResourceAttr(resourceName, "mock.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "func_graph_policy.#", "0"),
				),
			},
			{
				Config: testAccApi_orchestration_step2(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					// Web backend
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "request_params.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "request_params.0.orchestrations.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "request_params.0.orchestrations.0",
						"huaweicloud_apig_orchestration_rule.type_none_value", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "request_params.0.orchestrations.1",
						"huaweicloud_apig_orchestration_rule.type_hash", "id"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.path", "/orchestration/test/backend/hash"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.0.source", "orchestration"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.0.mapped_param_name", "ServiceName"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.0.mapped_param_location", "header"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.0.conditions.0.value", "HashValue"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApiResourceImportStateFunc(resourceName),
			},
		},
	})
}

func testAccApi_orchestration_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_orchestration_rule" "type_list" {
  instance_id = "%[2]s"
  name        = "%[3]s_type_list"
  strategy    = "list"

  mapped_param = jsonencode({
    "mapped_param_name": "ServiceName",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
  
  map = [
    jsonencode({
      "mapped_param_value": "ValueAA",
      "map_param_list": ["ValueA"]
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_hash" {
  instance_id = "%[2]s"
  name        = "%[3]s_type_hash"
  strategy    = "hash"

  mapped_param = jsonencode({
    "mapped_param_name": "ServiceName",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
}

resource "huaweicloud_apig_orchestration_rule" "type_none_value" {
  instance_id = "%[2]s"
  name        = "%[3]s_type_none_value"
  strategy    = "none_value"
  
  mapped_param = jsonencode({
    "mapped_param_name": "ServiceName",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
  
  map = [
    jsonencode({
      "mapped_param_value": "NoneValue"
    })
  ]
}
`, testAccApi_base(name),
		acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name)
}

func testAccApi_orchestration_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api" "test" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s"
  type                    = "Private"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/orchestration/test"
  security_authentication = "APP"

  request_params {
    name     = "X-Service-Name"
    type     = "STRING"
    location = "HEADER"
    maximum  = 30
    minimum  = 5

    orchestrations = [
      # None value type has the highest priority.
      huaweicloud_apig_orchestration_rule.type_none_value.id,
      huaweicloud_apig_orchestration_rule.type_list.id,
    ]
  }

  backend_params {
    type              = "REQUEST"
    name              = "ServiceName"
    location          = "HEADER"
    value             = "X-Service-Name"
    system_param_type = "backend"
  }

  web {
    path             = "/orchestration/test/backend"
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 40000
    retry_count      = 2
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id
  }

  web_policy {
    name             = "%[2]s_orchestration_policy"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ALL"
    path             = "/orchestration/test/backend/list"
    timeout          = 40000
    retry_count      = 2
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id

    backend_params {
      type              = "SYSTEM"
      name              = "ServiceName"
      location          = "HEADER"
      value             = "X-Service-Name"
      system_param_type = "backend"
    }

    conditions {
      source                = "orchestration"
      mapped_param_name     = "ServiceName"
      mapped_param_location = "header"
      value                 = "ValueAA"
    }
  }
}
`, baseConfig, name)
}

func testAccApi_orchestration_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api" "test" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s"
  type                    = "Private"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/orchestration/test"
  security_authentication = "APP"

  request_params {
    name     = "X-Service-Name"
    type     = "STRING"
    location = "HEADER"
    maximum  = 30
    minimum  = 5

    orchestrations = [
      # None value type has the highest priority.
      huaweicloud_apig_orchestration_rule.type_none_value.id,
      huaweicloud_apig_orchestration_rule.type_hash.id,
    ]
  }

  backend_params {
    type              = "REQUEST"
    name              = "ServiceName"
    location          = "HEADER"
    value             = "X-Service-Name"
    system_param_type = "backend"
  }

  web {
    path             = "/orchestration/test/backend"
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 40000
    retry_count      = 2
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id
  }

  web_policy {
    name             = "%[2]s_orchestration_policy"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ALL"
    path             = "/orchestration/test/backend/hash"
    timeout          = 40000
    retry_count      = 2
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    authorizer_id    = huaweicloud_apig_custom_authorizer.backend.id

    backend_params {
      type              = "SYSTEM"
      name              = "ServiceName"
      location          = "HEADER"
      value             = "X-Service-Name"
      system_param_type = "backend"
    }

    conditions {
      source                = "orchestration"
      mapped_param_name     = "ServiceName"
      mapped_param_location = "header"
      value                 = "HashValue"
    }
  }
}
`, baseConfig, name)
}
