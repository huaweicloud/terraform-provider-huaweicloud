---
subcategory: "API Gateway (APIG)"
---

# huaweicloud\_api\_gateway\_api

Provides an API gateway API resource.

## Example Usage

```hcl
resource "huaweicloud_api_gateway_group" "tf_apigw_group" {
  name        = "tf_apigw_group"
  description = "your descpiption"
}

resource "huaweicloud_api_gateway_api" "tf_apigw_api" {
  group_id                 = huaweicloud_api_gateway_group.tf_apigw_group.id
  name                     = "tf_apigw_api"
  description              = "your descpiption"
  tags                     = ["tag1", "tag2"]
  visibility               = 2
  auth_type                = "IAM"
  backend_type             = "HTTP"
  request_protocol         = "HTTPS"
  request_method           = "GET"
  request_uri              = "/test/path1"
  example_success_response = "example response"

  http_backend {
    protocol   = "HTTPS"
    method     = "GET"
    uri        = "/web/openapi"
    url_domain = "myhuaweicloud.com"
    timeout    = 10000
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the API resource. If omitted, the provider-level region will work as default. Changing this creates a new API resource.

* `name` - (Required) Specifies the name of the API. An API name consists of 3–64 characters,
    starting with a letter. Only letters, digits, and underscores (_) are allowed.

* `group_id` - (Required) Specifies the ID of the API group.
    Changing this creates a new resource.

* `description` - (Optional) Specifies the description of the API.
    The description cannot exceed 255 characters.

* `visibility` - (Optional) Specifies whether the API is available to the public.
    The value can be 1 (public) and 2 (private). Defaults to 2.

* `auth_type` - (Required) Specifies the security authentication mode.
     The value can be 'App', 'IAM', and 'NONE'.

* `request_protocol` - (Optional) Specifies the request protocol. The value can be 'HTTP', 'HTTPS', and 'BOTH'
    which means the API can be accessed through both 'HTTP' and 'HTTPS'. Defaults to 'HTTPS'.

* `request_method` - (Required) Specifies the request method, including 'GET','POST','PUT' and etc..

* `request_uri` - (Required) Specifies the request path of the API. The value must comply with URI specifications.

* `backend_type` - (Required) Specifies the service backend type. The value can be:
    - 'HTTP': the web service backend
    - 'FUNCTION': the FunctionGraph service backend
    - 'MOCK': the Mock service backend
  
* `http_backend` - (Optional) Specifies the configuration when backend_type selected 'HTTP' (documented below).
* `function_backend` - (Optional) Specifies the configuration when backend_type selected 'FUNCTION' (documented below).
* `mock_backend` - (Optional) Specifies the configuration when backend_type selected 'MOCK' (documented below).

* `request_parameter` - (Optional) the request parameter list (documented below).
* `backend_parameter` - (Optional) the backend parameter list (documented below).

* `tags` - (Optional) the tags of API in format of string list.

* `version` - (Optional) Specifies the version of the API. A maximum of 16 characters are allowed.

* `cors` - (Optional) Specifies whether CORS is supported or not.

* `example_success_response` - (Required) Specifies the example response for a successful request.
    The length cannot exceed 20,480 characters.

* `example_failure_response` - (Optional) Specifies the example response for a failed request
    The length cannot exceed 20,480 characters.

The `http_backend` object supports the following:

* `protocol` - (Required) Specifies the backend request protocol. The value can be 'HTTP' and 'HTTPS'.
* `method` - (Optional) Specifies the backend request method, including 'GET','POST','PUT' and etc..
* `uri` - (Required) Specifies the backend request path. The value must comply with URI specifications.
* `vpc_channel` - (Optional) Specifies the VPC channel ID. This parameter and `url_domain` are alternative.
* `url_domain` - (Optional) Specifies the backend service address. An endpoint URL is in the format of
     "domain name (or IP address):port number", with up to 255 characters. This parameter and `vpc_channel` are alternative.
* `timeout` - (Optional) Timeout duration (in ms) for API Gateway to request for the backend service. Defaults to 50000. 

The `function_backend` object supports the following:

* `function_urn` - (Required) Specifies the function URN.
* `invocation_type` - (Required) Specifies the invocation mode, which can be 'async' or 'sync'.
* `version` - (Optional) Specifies the function version.
* `timeout` - (Optional) Timeout duration (in ms) for API Gateway to request for FunctionGraph. Defaults to 50000.

The `mock_backend` object supports the following:

* `result_content` (Optional) Specifies the return result.
* `version` (Optional) Specifies the version of the Mock backend.
* `description` (Optional) Specifies the description of the Mock backend. The description cannot exceed 255 characters.

The `request_parameter` object supports the following:

* `name` - (Required) Specifies the input parameter name. A parameter name consists of 1–32 characters, starting with a letter.
    Only letters, digits, periods (.), hyphens (-), and underscores (_) are allowed.
* `location` - (Required) Specifies the input parameter location, which can be 'PATH', 'QUERY' or 'HEADER'.
* `type` - (Required) Specifies the input parameter type, which can be 'STRING' or 'NUMBER'.
* `required` - (Optional) Specifies whether the parameter is mandatory or not.
* `default` - (Optional) Specifies the default value when the parameter is optional.
* `description` - (Optional) Specifies the description of the parameter. The description cannot exceed 255 characters.

The `backend_parameter` object supports the following:

* `name` - (Required) Specifies the parameter name. A parameter name consists of 1–32 characters, starting with a letter.
    Only letters, digits, periods (.), hyphens (-), and underscores (_) are allowed.
* `location` - (Required) Specifies the parameter location, which can be 'PATH', 'QUERY' or 'HEADER'.
* `type` - (Required) Specifies the parameter type, which can be 'REQUEST', 'CONSTANT', or 'SYSTEM'.
* `value` - (Required) Specifies the parameter value, which is a string of not more than 255 characters.
    The value varies depending on the parameter type:
    - 'REQUEST': parameter name in `request_parameter`
    - 'CONSTANT': real value of the parameter
    - 'SYSTEM': gateway parameter name
* `description` - (Optional) Specifies the description of the parameter. The description cannot exceed 255 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the API.
* `group_name` - The name of the API group to which the API belongs.

## Import

API can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_api_gateway_api.api "774438a28a574ac8a496325d1bf51807"
```
