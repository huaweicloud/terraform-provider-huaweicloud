---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
pageTitle: "HuaweiCloud: huaweicloud_dsc_operate_obs_audit"
description: |-
  Manages a resource to operate DSC OBS audit within HuaweiCloud.
---

# huaweicloud_dsc_operate_obs_audit

Manages a resource to operate DSC OBS audit within HuaweiCloud.

-> This resource is a one-time action resource used to operate DSC OBS audit.
Deleting this resource will not restore the audit status or undo the operation, but will only
remove the resource information from the tf state file.

## Example Usage

```hcl
variable "bucket_id" {
  type = string
}

resource "huaweicloud_dsc_operate_obs_audit" "test" {
  bucket_id        = var.bucket_id
  operation_status = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `bucket_id` - (Required, String, NonUpdatable) Specifies the OBS bucket asset ID.

* `operation_status` - (Required, Bool, NonUpdatable) Specifies whether to enable or disable audit.
  **true** means enabling audit, **false** means disabling audit.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message describing the operation result or error information.

* `status` - The returned status, such as '200' or '400'.
