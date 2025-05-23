---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_check"
description: |-
  Use this resource to check whether the API name or path already exists within HuaweiCloud.
---

# huaweicloud_apig_api_check

Use this resource to check whether the API name or path already exists within HuaweiCloud.

-> This resource is only a one-time resource for checking the API definition. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Verify whether the API name already exists in the same group

```hcl
variable "instance_id" {}
variable "check_name" {}
variable "group_id" {}

resource "huaweicloud_apig_api_check" "test" {
  instance_id = var.instance_id
  type        = "name"
  name        = var.check_name
  group_id    = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the API belongs.

* `type` - (Required, String) Specifies the type of the API to be checked.  
  The valid values are as follows:
  + **name**
  + **path**

* `api_id` - (Optional, String) Specifies the ID of the API to be excluded from the check.

* `name` - (Optional, String) Specifies the name of the API.  
  This parameter is required if `type` is set to **name**.

* `group_id` - (Optional, String) Specifies the ID of the group to which the API belongs.  
  This parameter is required when verifying whether the API definition under the specified group exists.

* `match_mode` - (Optional, String) Specifies the matching mode of the API.  
  This parameter is required if `type` is set to **path**.  
  The valid values are as follows:
  + **SWA**: Prefix match.
  + **NORMAL**: Exact match.

* `req_method` - (Optional, String) Specifies the request method of the API.  
  This parameter is required if `type` is set to **path**.  
  The valid values are as follows:
  + **GET**
  + **POST**
  + **PUT**
  + **DELETE**
  + **HEAD**
  + **PATCH**
  + **OPTIONS**
  + **ANY**

* `req_uri` - (Optional, String) Specifies the request path of the API.  
  This parameter is required if `type` is set to **path**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
