---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_host_manual_detection"
description: |-
  Manages an HSS host manual detection operation resource within HuaweiCloud.
---

# huaweicloud_hss_host_manual_detection

Manages an HSS host manual detection operation resource within HuaweiCloud.

-> This resource is a one-time action resource using to operation HSS host manual detection. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "host_id" {}
variable "type" {}

resource "huaweicloud_hss_host_manual_detection" "test" {
  host_id = var.host_id
  type    = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `host_id` - (Required, String, NonUpdatable) Specifies the host ID.

* `type` - (Required, String, NonUpdatable) Specifies the type of detection.  
  The valid values are as follows:
  + **pwd**: Weak password detection.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
