---
subcategory: "API Gateway (Dedicated APIG)"
---

# huaweicloud_apig_response

Manages an APIG (API) custom response resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "response_name" {}

resource "huaweicloud_apig_response" "test" {
  instance_id = var.instance_id
  group_id    = var.group_id
  name        = var.response_name

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the API custom response is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the API group and the
  API custom response belong.  
  Changing this will create a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the API group to which the API custom response
  belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the API custom response.  
  The valid length is limited from `1` to `64`, letters, digits, hyphens (-) and underscores (_) are allowed.

* `rule` - (Optional, List) Specifies the API custom response rules definition.  
  The [object](#custom_response_rule) structure is documented below.

<a name="custom_response_rule"></a>
The `rule` block supports:

* `error_type` - (Required, String) Specifies the error type of the API response rule.
  + **AUTH_FAILURE**: Authentication failed.
  + **AUTH_HEADER_MISSING**: The identity source is missing.
  + **AUTHORIZER_FAILURE**: Custom authentication failed.
  + **AUTHORIZER_CONF_FAILURE**: There has been a custom authorizer error.
  + **AUTHORIZER_IDENTITIES_FAILURE**: The identity source of the custom authorizer is invalid.
  + **BACKEND_UNAVAILABLE**: The backend service is unavailable.
  + **BACKEND_TIMEOUT**: Communication with the backend service timed out.
  + **THROTTLED**: The request was rejected due to request throttling.
  + **UNAUTHORIZED**: The app you are using has not been authorized to call the API.
  + **ACCESS_DENIED**: Access denied.
  + **NOT_FOUND**: No API is found.
  + **REQUEST_PARAMETERS_FAILURE**: The request parameters are incorrect.
  + **DEFAULT_4XX**: Another 4XX error occurred.
  + **DEFAULT_5XX**: Another 5XX error occurred.

* `body` - (Required, String) Specifies the body template of the API response rule, e.g.
  `{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}`

* `status_code` - (Optional, Int) Specifies the HTTP status code of the API response rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the API custom response.
* `created_at` - The creation time of the API custom response.
* `updated_at` - The latest update time of the API custom response.

## Import

API Responses can be imported using their `name` and IDs of the APIG dedicated instances and API groups to which the API
response belongs, separated by slashes, e.g.

```shell
$ terraform import huaweicloud_apig_response.test <instance_id>/<group_id>/<name>
```
