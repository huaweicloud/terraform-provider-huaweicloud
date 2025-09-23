---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_api_debug"
description: |-
  Use this resource to debug API for DataArts Data Service within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_api_debug

Use this resource to debug API for DataArts Data Service within HuaweiCloud.

-> 1. Only exclusive API can be debugged.
   <br>2. This resource is only a one-time action resource for debugging the API. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "debug_api_id" {}
variable "exclusive_cluster_id" {}

# The API already has an request parameter and names 'test_request_field'.
resource "huaweicloud_dataarts_dataservice_api_debug" "test" {
  workspace_id = var.workspace_id
  api_id       = var.debug_api_id
  instance_id  = var.exclusive_cluster_id

  params = jsonencode({
    "page_num": "1",     # Default parameter
    "page_size": "100",  # Default parameter
    "test_request_field": "{\"foo\": \"bar\"}"
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of workspace where the API is located.

  Changing this parameter will create a new resource.

* `api_id` - (Required, String, ForceNew) Specifies the ID of the catalog where the API is located.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the exclusive cluster ID.
  Changing this parameter will create a new resource.

* `params` - (Optional, String, ForceNew) Specifies the request parameters in which to debug the API, in JSON format.
  Changing this parameter will create a new resource.

  -> There are two default and required parameter `page_num` and `page_size` that need to be defined.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also request ID for this debug operation.

* `url` - The result detail of this API debug.

* `result` - The ID of the catalog where the API is located.

* `timeout` - The timeout of this API debug.

* `request_header` - The request header of this API debug result, in JSON format.

* `response_header` - The response header of this API debug result, in JSON format.
