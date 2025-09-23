---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_debug"
description: |-
  Use this resource to debug API within HuaweiCloud.
---

# huaweicloud_apig_api_debug

Use this resource to debug API within HuaweiCloud.

-> This resource is only a one-time action resource for debugging API. Deleting this resource will not clear
   the corresponding debug record, but will only remove the resource information from the tfstate file.

## Example Usage

### Debug the API in DEVELOPER mode

```hcl
variable "instance_id" {}
variable "api_id" {}

resource "huaweicloud_apig_api_debug" "test" {
  instance_id = var.instance_id
  api_id      = var.api_id
  mode        = "DEVELOPER"
  scheme      = "HTTPS"
  method      = "GET"
  path        = "/terraform/mock"
}
```

### Debug the API with request body and headers

```hcl
variable "instance_id" {}
variable "api_id" {}

resource "huaweicloud_apig_api_debug" "test" {
  instance_id = var.instance_id
  api_id      = var.api_id
  mode        = "DEVELOPER"
  scheme      = "HTTPS"
  method      = "POST"
  path        = "/terraform/mock"
  body        = "{\"terraform\": \"test\"}"
  
  header = jsonencode({
    "Content-Type": ["application/json"],
  })
  
  query = jsonencode({
    "terraform": ["test"]
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the dedicated instance to which the API belongs is
located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the API belongs.

* `api_id` - (Required, String, NonUpdatable) Specifies the ID of the API to be debugged.

* `mode` - (Required, String, NonUpdatable) Specifies the debug mode.  
  The valid values are as follows:
  + **DEVELOPER**: Debug the API definition that has not been published.
  + **MARKET**: Debug the API purchased from the cloud store.
  + **CONSUMER**: Debug the API definition in the specified runtime environment.

* `scheme` - (Required, String, NonUpdatable) Specifies the request protocol.  
  The valid values are as follows:
  + **HTTP**: Make a request via the HTTP protocol.
  + **HTTPS**: Make a request via the HTTPS protocol.

* `method` - (Required, String, NonUpdatable) Specifies the request method of the API.

* `path` - (Required, String, NonUpdatable) Specifies the request path of the API.  
  The path must start with `\` and have a maximum length of `1,024` characters.

* `body` - (Optional, String, NonUpdatable) Specifies the request message body of the API.  
  The maximum length is `2,097,152` bytes.
  If it exceeds this limit, the excess part will be truncated.

* `header` - (Optional, String, NonUpdatable) Specifies the request header parameters of the API, in JSON format.  
  Each parameter value is a list of strings. Parameter names have the following constraints:
  + Composed of English letters, numbers, dots(.), and hyphens(-)
  + Must start with an English letter and have a maximum length of `32` bytes
  + Cannot start with `X-Apig-` or `X-Sdk-` (case-insensitive)
  + Cannot be `X-Stage` (case-insensitive)
  + When mode is **MARKET** or **CONSUMER**, cannot be `X-Auth-Token` or `Authorization` (case-insensitive)

* `query` - (Optional, String, NonUpdatable) Specifies the request query parameters of the API, in JSON format.  
  Each parameter value is a list of strings. Parameter names have the following constraints:
  + Composed of English letters, numbers, dots(.), underscores(_), and hyphens(-)
  + Must start with an English letter and have a maximum length of `32` bytes
  + Cannot start with `X-Apig-` or `X-Sdk-` (case-insensitive)
  + Cannot be `X-Stage` (case-insensitive)

* `stage` - (Optional, String, NonUpdatable) Specifies the runtime environment for debug request.  
  This parameter is only valid when mode is **CONSUMER**.
  If not provided, the default value is **RELEASE**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `request` - The debug request message content.

* `response` - The debug response message content.

* `latency` - The debug latency in milliseconds.
