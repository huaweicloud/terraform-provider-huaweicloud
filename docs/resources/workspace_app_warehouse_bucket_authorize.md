---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_warehouse_bucket_authorize"
description: |-
  Use this resource to assign a bucket for the app repository within HuaweiCloud.
---

# huaweicloud_workspace_app_warehouse_bucket_authorize

Use this resource to assign a bucket for the app repository within HuaweiCloud.

-> This resource is only a one-time action resource for assigning an app repository bucket. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

## Example Usage

```hcl
variable "bucket_name" {}

resource "huaweicloud_workspace_app_warehouse_bucket_authorize" "test" {
  bucket_name = var.bucket_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the app repository bucket is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `bucket_name` - (Optional, String, NonUpdatable) Specifies the name of the bucket to be assigned.  
  + The valid length is limited `3` to `63` characters.
  + Only lowercase letters, digits, hyphens(-) and dots(.) are allowed.
  + Must start with a digit or lowercase letter.
    * IP addresses are not allowed.
    * Cannot start or end with "-" or ".".
    * Cannot have two consecutive "." (e.g.,"my..bucket").
    * Cannot have "." and "-" adjacent (e.g., "my-.bucket" and "my.-bucket").

  The app repository bucket will get a random name if the `bucket_name` is not given.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
