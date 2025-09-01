---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_offline"
description: |-
  Use this resource to take an API offline by version within HuaweiCloud.
---

# huaweicloud_apig_api_offline

Use this resource to take an API offline by version within HuaweiCloud.

-> This resource is only a one-time action resource for taking API offline by version. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "version_id" {}

resource "huaweicloud_apig_api_offline" "test" {
  instance_id = var.instance_id
  version_id  = var.version_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the APIG dedicated instance is located.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the API belongs.

* `version_id` - (Required, String) Specifies the version ID of the API to be taken offline.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
