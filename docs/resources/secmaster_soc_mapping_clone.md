---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_mapping_clone"
description: |-
  Manages a resource to clone mapping within HuaweiCloud.
---

# huaweicloud_secmaster_soc_mapping_clone

Manages a resource to clone mapping within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "mapping_id" {}
variable "name" {}

resource "huaweicloud_secmaster_soc_mapping_clone" "test" {
  workspace_id = var.workspace_id
  mapping_id   = var.mapping_id
  name         = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Sepcifies the ID of the workspace to which the mapping belongs.

* `mapping_id` - (Required, String, NonUpdatable) Sepcifies the ID of the mapping.

* `name` - (Required, String, NonUpdatable) Sepcifies the name of the clone generated mapping.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
