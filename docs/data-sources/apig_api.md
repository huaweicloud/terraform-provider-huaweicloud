---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api"
description: |-
  Use this data source to get the configuration details of the API within HuaweiCloud.
---

# huaweicloud_apig_api

Use this data source to get the configuration details of the API within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "api_id" {}

data "huaweicloud_apig_api" "test" {
  instance_id = var.instance_id
  api_id      = var.api_id     
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the API belong.

* `api_id` - (Required, String) Specifies the ID of the API.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - The name of the API.

* `type` - The type of the API.

* `request_method` - The request method of the API.

* `request_path` - The request address of the API.

* `request_protocol` - The request protocol of the API.

* `security_authentication` - The security authentication mode of the API request.

* `simple_authentication` - Whether the authentication of the application code is enabled.

* `authorizer_id` - The ID of the authorizer to which the API request used.

* `tags` - The list of tags configuration.

* `group_id` - The group ID corresponding to the API.

* `group_name` - The group name corresponding to the API.

* `group_version` - The version of group corresponding to the API.

* `request_params` - The configuration list of the front-end parameters.
  The [request_params](#api_request_params) structure is documented below.

* `backend_params` - The configuration list of the backend parameters.
  The [backend_params](#api_backend_params) structure is documented below.

* `body_description` - The description of the API request body.

* `cors` - Whether CORS is supported.

* `description` - The description of the API.

* `matching` - The matching mode of the API.
  + **Exact**
  + **Prefix**

* `response_id` - The response ID of the corresponding APIG group.

* `success_response` - The example response for a successful request.

* `failure_response` - The example response for a failure request.

* `mock` - The mock backend details.
  The [mock](#api_mock) structure is documented below.

* `mock_policy` - The policy backends of the mock.
  The [mock_policy](#api_mock_policy) structure is documented below.
  
* `func_graph` - The FunctionGraph backend details.
  The [func_graph](#api_func_graph) structure is documented below.

* `func_graph_policy` - The policy backends of the FunctionGraph.
  The [func_graph_policy](#api_func_graph_policy) structure is documented below.

* `web` - The web backend details.
  The [web](#api_web) structure is documented below.

* `web_policy` - The policy backends of the web.
  The [web_policy](#api_web_policy) structure is documented below.
  
* `env_id` - The name of the environment where the API is published.

* `env_name` - The name of the environment where the API is published.

* `publish_id` - The ID of publish corresponding to the API.

* `backend_type` - The backend type of the API.

* `published_at` - The published time of the API, in RFC3339 format.

* `registered_at` - The registered time of the API, in RFC3339 format.

* `updated_at` - The latest update time of the API, in RFC3339 format.

<a name="api_request_params"></a>
The `request_params` block supports:

* `id` - The ID of the request parameter.

* `name` - The name of the request parameter.

* `required` - Whether this parameter is required.

* `passthrough` - Whether to transparently transfer the parameter.

* `enumeration` - The enumerated value.

* `location` - Where this parameter is located.

* `description` - The parameter description.

* `type` - The parameter type.

* `maximum` - The maximum value or length (string parameter) for parameter.

* `minimum` - The minimum value or length (string parameter) for parameter.

* `example` - The parameter example.

* `default` - The default value of the parameter.

* `valid_enable` - Whether to enable the parameter validation.
  + **1**: enable
  + **2**: disable

<a name="api_backend_params"></a>
The `backend_params` block supports:

* `id` - The ID of the backend parameter.

* `request_id` - The ID of the corresponding request parameter.

* `type` - The name of parameter.

* `name` - The name of parameter.

* `location` - Where the parameter is located.

* `value` - The value of the parameter.

* `description` - The description of the constant or system parameter.

* `system_param_type` - The type of the system parameter.

<a name="api_mock"></a>
The `mock` block supports:

* `id` - The ID of the mock backend configuration.

* `status_code` - The custom status code of the mock response.

* `response` - The response of the mock backend configuration.

* `authorizer_id` - The ID of the backend custom authorization.

<a name="api_mock_policy"></a>
The `mock_policy` block supports:

* `id` - The ID of the mock backend policy.

* `name` - The backend policy name.

* `status_code` - The custom status code of the mock response.

* `response` - The response of the backend policy.

* `conditions` - The policy conditions.
  The [conditions](#policy_conditions) structure is documented below.

* `effective_mode` - The effective mode of the backend policy.

* `backend_params` - The configuration list of backend parameters.
  The [backend_params](#api_backend_params) structure is documented below.

* `authorizer_id` - The ID of the backend custom authorization.

<a name="api_func_graph"></a>
The `func_graph` block supports:

* `id` - The ID of the FunctionGraph backend configuration.

* `function_urn` - The URN of the FunctionGraph function.

* `version` - The version of the FunctionGraph function.

* `function_alias_urn` - The alias URN of the FunctionGraph function.  

* `network_type` - The network architecture (framework) type of the FunctionGraph function.
  **V1**: Non-VPC network framework.
  **V2**: VPC network framework.

* `request_protocol` - The request protocol of the FunctionGraph function.  

* `timeout` - The timeout for API requests to backend service.

* `version` - The version of the FunctionGraph function.

* `invocation_type` - The invocation type.

* `authorizer_id` - The ID of the backend custom authorization.

<a name="api_func_graph_policy"></a>
The `func_graph_policy` block supports:

* `id` - The ID of the FunctionGraph backend policy.

* `name` - The name of the backend policy.

* `function_urn` - The URN of the FunctionGraph function.

* `version` - The version of the FunctionGraph function.

* `function_alias_urn` - The alias URN of the FunctionGraph function.  

* `network_type` - The network architecture (framework) type of the FunctionGraph function.
  **V1**: Non-VPC network framework.
  **V2**: VPC network framework.

* `request_protocol` - The request protocol of the FunctionGraph function.  

* `conditions` - The policy conditions.
  The [conditions](#policy_conditions) structure is documented below.
  
* `invocation_type` - The invocation mode of the FunctionGraph function.

* `effective_mode` - The effective mode of the backend policy.

* `timeout` - The timeout for API requests to backend service.

* `backend_params` - The configaiton list of the backend parameters.
  The [backend_params](#api_backend_params) structure is documented below.

* `authorizer_id` - The ID of the backend custom authorization.

<a name="api_web"></a>
The `web` block supports:

* `id` - The ID of the backend configuration.

* `path` - The backend request path.

* `host_header` - The proxy host header.

* `vpc_channel_id` - The VPC channel ID.

* `backend_address` - The backend service address.

* `request_method` - The backend request method of the API.

* `request_protocol` - The web protocol type of the API request.

* `timeout` - The timeout for API requests to backend service.

* `retry_count` - The number of retry attempts to request the backend service.

* `ssl_enable` - Whether to enable two-way authentication.

* `authorizer_id` - The ID of the backend custom authorization.

<a name="api_web_policy"></a>
The `web_policy` block supports:

* `id` - The ID of the web policy.

* `name` - The name of the web policy.

* `path` - The backend request address.

* `request_method` - The backend request method of the API.

* `request_protocol` - The backend request protocol.

* `conditions` - The policy conditions.
  The [conditions](#policy_conditions) structure is documented below.

* `host_header` - The proxy host header.

* `vpc_channel_id` - The VPC channel ID.

* `backend_address` - The backend service address

* `effective_mode` - The effective mode of the backend policy.

* `timeout` - The timeout for API requests to backend service.

* `retry_count` - The number of retry attempts to request the backend service.

* `backend_params` - The configuration list of the backend parameters.
  The [backend_params](#api_backend_params) structure is documented below.

* `authorizer_id` - The ID of the backend custom authorization.

<a name="policy_conditions"></a>
The `conditions` block supports:

* `id` - The ID of the backend policy condition.

* `value` - The value corresponding to the parameter name.

* `param_name` - The request parameter name.

* `sys_name` - The gateway built-in parameter name.

* `cookie_name` - The cookie parameter name.

* `frontend_authorizer_name` - The frontend authentication parameter name.

* `source` - The type of the backend policy.

* `type` - The condition type of the backend policy.
  + **Equal**
  + **Enumerated**
  + **Matching**

* `request_id` - The ID of the corresponding request parameter.

* `request_location` - The location of the corresponding request parameter.
