---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_version_unpublish"
description: |-
  Use this resource to unpublish API by version within HuaweiCloud.
---

# huaweicloud_apig_api_version_unpublish

Use this resource to unpublish API by version within HuaweiCloud.

-> This resource is only a one-time action resource for unpublish API by version. Deleting this resource
   will not restore the API version, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "version_id" {}

resource "huaweicloud_apig_api_version_unpublish" "test" {
  instance_id = var.instance_id
  version_id  = var.version_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the dedicated instance to which the API version
belongs is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the API version
belongs.

* `version_id` - (Required, String, NonUpdatable) Specifies the ID of the API version to be unpublish.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
