---
subcategory: "VPCEP"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_endpoint_upgrade"
description: |-
  Use this resource to upgrade VPCEP endpoint from basic version to professional version within HuaweiCloud.
---

# huaweicloud_vpcep_endpoint_upgrade

Use this resource to upgrade VPCEP endpoint from basic version to professional version within HuaweiCloud.

-> This resource is only a one-time action resource for upgrading VPCEP endpoint. Deleting this
  resource will not clear the corresponding request record, but will only remove the resource information from the
  tfstate file.

## Example Usage

```hcl
variable "endpoint_id" {}

resource "huaweicloud_vpcep_endpoint_upgrade" "test" {
  endpoint_id = var.endpoint_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the endpoint is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `endpoint_id` - (Required, String, NonUpdatable) Specifies the ID of the endpoint to be upgraded.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
