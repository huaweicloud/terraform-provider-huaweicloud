---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_response"
description: ""
---

# huaweicloud_apig_response

Manages an APIG (API) custom response resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "response_name" {}
variable "response_headers" {
  type = list(object({
    key   = string
    value = string
  }))
}

resource "huaweicloud_apig_response" "test" {
  instance_id = var.instance_id
  group_id    = var.group_id
  name        = var.response_name

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  
    dynamic "headers" {
      for_each = var.response_headers

      content {
        key   = headers.value["key"]
        value = headers.value["value"]
      }
    }
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
  The [rule](#custom_response_rule) structure is documented below.

<a name="custom_response_rule"></a>
The `rule` block supports:

* `error_type` - (Required, String) Specifies the error type of the API response rule.
  The valid values and the related default status code are as follows:
  + **ACCESS_DENIED**: (**403**) Access denied.
  + **AUTH_FAILURE**: (**401**) Authentication failed.
  + **AUTH_HEADER_MISSING**: (**401**) The identity source is missing.
  + **AUTHORIZER_CONF_FAILURE**: (**500**) There has been a custom authorizer error.
  + **AUTHORIZER_FAILURE**: (**500**) Custom authentication failed.
  + **AUTHORIZER_IDENTITIES_FAILURE**: (**401**) The identity source of the custom authorizer is invalid.
  + **BACKEND_TIMEOUT**: (**504**) Communication with the backend service timed out.
  + **BACKEND_UNAVAILABLE**: (**502**) The backend service is unavailable.
  + **NOT_FOUND**: (**404**) No API is found.
  + **REQUEST_PARAMETERS_FAILURE**: (**400**) The request parameters are incorrect.
  + **THROTTLED**: (**429**) The request was rejected due to request throttling.
  + **UNAUTHORIZED**: (**401**) The app you are using has not been authorized to call the API.
  + **DEFAULT_4XX**: (**NONE**) Another 4XX error occurred.
  + **DEFAULT_5XX**: (**NONE**) Another 5XX error occurred.
  + **THIRD_AUTH_CONF_FAILURE**: (**500**) Third-party authorizer configuration error.
  + **THIRD_AUTH_FAILURE**: (**401**) Third-party authentication failed.
  + **THIRD_AUTH_IDENTITIES_FAILURE**: (**401**) Identity source of the third-party authorizer is invalid.

* `body` - (Required, String) Specifies the body template of the API response rule, e.g.
  `{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}`

* `status_code` - (Optional, Int) Specifies the HTTP status code of the API response rule.
  The valid value is range from `200` to `599`.

* `headers` - (Optional, List) Specifies the configuration of the custom response headers.  
  The [headers](#custom_response_rule_headers) structure is documented below.

<a name="custom_response_rule_headers"></a>
The `headers` block supports:

* `key` - (Required, String) Specifies the key name of the response header.
   The valid length is limited from `1` to `128`, only English letters, digits and hyphens (-) are allowed.

* `value` - (Required, String) Specifies the value for the specified response header key.
  The valid length is limited from `1` to `1,024`.

## Attribute Reference

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
