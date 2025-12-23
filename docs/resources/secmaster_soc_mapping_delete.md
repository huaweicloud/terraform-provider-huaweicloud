---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_mapping_delete"
description: |-
  Manages a resource to delete mapping within HuaweiCloud.
---

# huaweicloud_secmaster_soc_mapping_delete

Manages a resource to delete mapping within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "mapping_id" {}

resource "huaweicloud_secmaster_soc_mapping_delete" "test" {
  workspace_id = var.workspace_id
  mapping_id   = var.mapping_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `mapping_id` - (Required, String, NonUpdatable) Specifies the mapping ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
