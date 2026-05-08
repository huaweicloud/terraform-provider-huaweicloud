---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_sessions_query"
description: |-
  Manages a DCS sessions query resource within HuaweiCloud.
---

# huaweicloud_dcs_sessions_query

Manages a DCS sessions query resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_dcs_sessions_query" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  clean_cache = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to log in to WebCli.
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `node_id` - (Required, String, NonUpdatable) Specifies the ID of a DDS instance node.

* `clean_cache` - (Optional, Bool, NonUpdatable) Whether to re-query and save the session list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
