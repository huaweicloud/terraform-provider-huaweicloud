---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_rename_commands"
description: |-
  Use this data source to query the rename commands of a DCS instance within HuaweiCloud.
---

# huaweicloud_dcs_instance_rename_commands

Use this data source to query the rename commands of a DCS instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_instance_rename_commands" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the rename commands.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rename_commands` - The list of renamed commands for the instance.
