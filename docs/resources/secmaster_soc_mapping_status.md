---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_mapping_status"
description: |-
  Manages a SecMaster mapping status resource within HuaweiCloud.
---

# huaweicloud_secmaster_soc_mapping_status

Manages a SecMaster mapping status resource within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not change the current mapping status,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "mapping_id" {}
variable "status" {}

resource "huaweicloud_secmaster_soc_mapping_status" "test" {
  workspace_id = var.workspace_id
  mapping_id   = var.mapping_id
  status       = var.status
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `mapping_id` - (Required, String, NonUpdatable) Specifies the mapping ID.

* `status` - (Required, String) Specifies the status of the mapping.
  The valid values are as follows:
  + **enabled**: Enabled.
  + **disabled**: Disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
